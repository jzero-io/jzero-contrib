package zinclog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/hpcloud/tail"
	"github.com/zeromicro/go-zero/core/logx"
)

// ZincLogstash implements Logstash interface for ZincSearch
type ZincLogstash struct {
	baseURL     string
	username    string
	password    string
	serviceName string
	logsPath    []string

	wg         sync.WaitGroup
	stopChan   chan struct{}
	notifyChan map[string]chan bool
	mu         sync.Mutex
}

// NewZincLogstash creates a new ZincLogstash instance
func NewZincLogstash(baseURL, username, password string, serviceName string, logsPath []string) *ZincLogstash {
	return &ZincLogstash{
		baseURL:     baseURL,
		username:    username,
		password:    password,
		serviceName: serviceName,
		logsPath:    logsPath,
		stopChan:    make(chan struct{}),
		notifyChan:  make(map[string]chan bool),
		mu:          sync.Mutex{},
	}
}

// SendLog sends log data to ZincSearch
func (z *ZincLogstash) SendLog(index, data string) error {
	url := fmt.Sprintf("%s/api/%s/_doc", z.baseURL, index)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(z.username, z.password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to send log to ZincSearch, status code: %d, body: %s",
			resp.StatusCode, string(respBody))
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

	notifyChan := make(chan bool)

	z.mu.Lock()
	defer z.mu.Unlock()
	z.notifyChan[path] = notifyChan

	defer func() {
		if tailer != nil {
			tailer.Stop()
		}
	}()

	// Initial setup
	fileSize, err = getFileSize(path)
	if err != nil {
		logx.Errorf("File %s doesn't exist: %s", path, err.Error())
		return
	}

	// Start tailing the file
	tailer, err = tail.TailFile(path, tail.Config{
		Follow:    true,
		ReOpen:    false,
		MustExist: true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: io.SeekEnd},
	})
	if err != nil {
		logx.Errorf("Failed to tail file %s: %v", path, err)
		return
	}

	// Start the tail goroutine
	go z.tailFile(path, tailer, &curIndex, notifyChan)

	// Monitor for file rotation
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-z.stopChan:
			notifyChan <- true
			return
		case <-ticker.C:
			newSize, err := getFileSize(path)
			if err != nil {
				continue
			}

			if newSize >= fileSize {
				fileSize = newSize
				continue
			}

			// File rotation detected
			fileSize = newSize

			if tailer != nil {
				tailer.Stop()
				tailer = nil
			}
			notifyChan <- true

			// Reopen the file
			for {
				select {
				case <-z.stopChan:
					return
				default:
					tailer, err = tail.TailFile(path, tail.Config{
						Follow:    true,
						ReOpen:    false,
						MustExist: true,
						Location:  &tail.SeekInfo{Offset: 0, Whence: io.SeekStart},
					})
					if err != nil {
						logx.Errorf("Failed to reopen file %s: %v", path, err)
						time.Sleep(100 * time.Millisecond)
						continue
					}

					notifyChan = make(chan bool)
					z.notifyChan[path] = notifyChan
					go z.tailFile(path, tailer, &curIndex, notifyChan)
					break
				}
				break
			}
		}
	}
}

func (z *ZincLogstash) tailFile(path string, tailer *tail.Tail, curIndex *string, notifyChan chan bool) {
	for {
		select {
		case line := <-tailer.Lines:
			if line == nil {
				continue
			}

			// Determine index based on current date
			year, month, day := time.Now().Date()
			indexNew := fmt.Sprintf("%s_%s_%d_%02d_%02d", z.serviceName,
				strings.ReplaceAll(path, "/", "_"), year, month, day)

			// Create index if needed
			if *curIndex != indexNew {
				err := z.createIndex(indexNew)
				if err != nil {
					logx.Errorf("Failed to create index %s: %v", indexNew, err)
					continue
				}
				*curIndex = indexNew
			}

			// Send log if not empty
			if line.Text != "" {
				if err := z.SendLog(*curIndex, line.Text); err != nil {
					logx.Errorf("Failed to send log from %s: %v", path, err)
				}
			}

		case <-notifyChan:
			logx.Infof("Stopping tail for file: %s", path)
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
		return err
	}

	url := fmt.Sprintf("%s/api/index", z.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(z.username, z.password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK || strings.Contains(string(b), "already exists") {
		return nil
	}

	return fmt.Errorf("failed to create index, status: %d, body: %s", resp.StatusCode, string(b))
}

func getFileSize(path string) (int64, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}
