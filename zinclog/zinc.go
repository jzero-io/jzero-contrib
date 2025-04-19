package zinclog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/hpcloud/tail"
	"github.com/zeromicro/go-zero/core/logx"
)

// ZincLogstash implements log collection and shipping to ZincSearch
type ZincLogstash struct {
	baseURL     string
	username    string
	password    string
	serviceName string
	logsPath    []string

	wg          sync.WaitGroup
	stopChan    chan struct{}
	notifyChans map[string]chan struct{} // Changed to chan struct{} for signaling
	mu          sync.Mutex
}

// NewZincLogstash creates a new ZincLogstash instance
func NewZincLogstash(baseURL, username, password, serviceName string, logsPath []string) *ZincLogstash {
	return &ZincLogstash{
		baseURL:     strings.TrimRight(baseURL, "/"),
		username:    username,
		password:    password,
		serviceName: serviceName,
		logsPath:    logsPath,
		stopChan:    make(chan struct{}),
		notifyChans: make(map[string]chan struct{}),
	}
}

// SendLog sends log data to ZincSearch
func (z *ZincLogstash) SendLog(index, data string) error {
	if strings.TrimSpace(data) == "" {
		return nil
	}

	url := fmt.Sprintf("%s/api/%s/_doc", z.baseURL, index)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(z.username, z.password)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// Start begins watching all configured log files
func (z *ZincLogstash) Start() {
	for _, path := range z.logsPath {
		z.wg.Add(1)
		go z.watchFile(path)
	}
}

// Stop gracefully stops all file watchers
func (z *ZincLogstash) Stop() {
	close(z.stopChan)
	z.wg.Wait()
}

func (z *ZincLogstash) watchFile(path string) {
	defer z.wg.Done()

	var fileSize int64
	var err error
	var curIndex string
	var tailer *tail.Tail

	// Create notification channel
	notifyChan := make(chan struct{})

	// Register the channel
	z.mu.Lock()
	z.notifyChans[path] = notifyChan
	z.mu.Unlock()

	// Cleanup on exit
	defer func() {
		z.mu.Lock()
		delete(z.notifyChans, path)
		z.mu.Unlock()

		if tailer != nil {
			_ = tailer.Stop()
		}
	}()

	// Initial file check
	fileSize, err = getFileSize(path)
	if err != nil {
		logx.Errorf("File %s doesn't exist: %v", path, err)
		return
	}

	// Start tailing the file
	tailer, err = tail.TailFile(path, tail.Config{
		Follow:    true,
		ReOpen:    false,
		MustExist: true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: io.SeekEnd},
		Logger:    tail.DiscardingLogger, // Discard tail's internal logs
	})
	if err != nil {
		logx.Errorf("Failed to tail file %s: %v", path, err)
		return
	}

	// Start the tail processing goroutine
	go z.tailFile(path, tailer, &curIndex, notifyChan)

	// Monitor for file rotation
	ticker := time.NewTicker(100 * time.Millisecond) // Reduced frequency
	defer ticker.Stop()

	for {
		select {
		case <-z.stopChan:
			close(notifyChan)
			return
		case <-ticker.C:
			newSize, err := getFileSize(path)
			if err != nil {
				if os.IsNotExist(err) {
					logx.Infof("File %s no longer exists", path)
					close(notifyChan)
					return
				}
				continue
			}

			// File rotation detected when size decreases
			if newSize < fileSize {
				fileSize = newSize

				if tailer != nil {
					_ = tailer.Stop()
					tailer = nil
				}
				close(notifyChan)

				// Reopen the file with new tailer
				for i := 0; i < 5; i++ { // Retry up to 5 times
					select {
					case <-z.stopChan:
						return
					default:
						tailer, err = tail.TailFile(path, tail.Config{
							Follow:    true,
							ReOpen:    false,
							MustExist: true,
							Location:  &tail.SeekInfo{Offset: 0, Whence: io.SeekStart},
							Logger:    tail.DiscardingLogger,
						})
						if err == nil {
							notifyChan = make(chan struct{})
							z.mu.Lock()
							z.notifyChans[path] = notifyChan
							z.mu.Unlock()
							go z.tailFile(path, tailer, &curIndex, notifyChan)
							break
						}
						logx.Errorf("Failed to reopen file %s (attempt %d): %v", path, i+1, err)
						time.Sleep(time.Second)
					}
				}
				if err != nil {
					logx.Errorf("Giving up on file %s after retries", path)
					return
				}
			} else {
				fileSize = newSize
			}
		}
	}
}

func (z *ZincLogstash) tailFile(path string, tailer *tail.Tail, curIndex *string, notifyChan <-chan struct{}) {
	logx.Infof("Started tailing file: %s", path)

	for {
		select {
		case line, ok := <-tailer.Lines:
			if !ok {
				return
			}
			if line == nil || line.Text == "" {
				continue
			}

			// Generate index name
			year, month, day := time.Now().Date()
			indexNew := fmt.Sprintf("%s_%s_%d%02d%02d",
				z.serviceName,
				sanitizeForIndex(path),
				year, month, day)

			// Create index if needed
			if *curIndex != indexNew {
				if err := z.createIndex(indexNew); err != nil {
					logx.Errorf("Failed to create index %s: %v", indexNew, err)
					continue
				}
				*curIndex = indexNew
			}

			// Send log
			if err := z.SendLog(*curIndex, line.Text); err != nil {
				logx.Errorf("Failed to send log from %s: %v", path, err)
			}

		case <-notifyChan:
			logx.Infof("Stopped tailing file: %s", path)
			return
		}
	}
}

func (z *ZincLogstash) createIndex(index string) error {
	config := struct {
		Name        string `json:"name"`
		StorageType string `json:"storage_type"`
		ShardNum    int    `json:"shard_num"`
	}{
		Name:        index,
		StorageType: "disk",
		ShardNum:    1,
	}

	body, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	url := fmt.Sprintf("%s/api/index", z.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(z.username, z.password)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	respBody, _ := io.ReadAll(resp.Body)
	if strings.Contains(string(respBody), "already exists") {
		return nil
	}

	return fmt.Errorf("failed to create index (status %d): %s", resp.StatusCode, string(respBody))
}

// sanitizeForIndex cleans up a path string for use in an index name
func sanitizeForIndex(path string) string {
	base := filepath.Base(path)
	return strings.ReplaceAll(
		strings.TrimSuffix(base, filepath.Ext(base)),
		".", "_")
}

func getFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}
