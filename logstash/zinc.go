package logstash

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hpcloud/tail"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// ZincLogstash implements Logstash interface for ZincSearch
type ZincLogstash struct {
	baseURL     string
	username    string
	password    string
	serviceName string
}

// NewZincLogstash creates a new ZincLogstash instance
func NewZincLogstash(baseURL, username, password string, serviceName string) *ZincLogstash {
	return &ZincLogstash{
		baseURL:     baseURL,
		username:    username,
		password:    password,
		serviceName: serviceName,
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

// FileWatcher manages watching multiple log files
type FileWatcher struct {
	paths      []string
	logstash   Logstash
	wg         sync.WaitGroup
	stopChan   chan struct{}
	notifyChan map[string]chan bool
	mu         sync.Mutex
}

// NewFileWatcher creates a new FileWatcher instance
func NewFileWatcher(paths []string, logstash Logstash) *FileWatcher {
	return &FileWatcher{
		paths:      paths,
		logstash:   logstash,
		stopChan:   make(chan struct{}),
		notifyChan: make(map[string]chan bool),
		mu:         sync.Mutex{},
	}
}

// Start begins watching all configured log files
func (fw *FileWatcher) Start() {
	for _, path := range fw.paths {
		fw.wg.Add(1)
		go fw.watchFile(path)
	}
}

// Stop gracefully stops all file watchers
func (fw *FileWatcher) Stop() {
	close(fw.stopChan)
	fw.wg.Wait()
}

func (fw *FileWatcher) watchFile(path string) {
	defer fw.wg.Done()

	var fileSize int64
	var err error
	var curIndex string
	var tailer *tail.Tail

	notifyChan := make(chan bool)

	fw.mu.Lock()
	defer fw.mu.Unlock()
	fw.notifyChan[path] = notifyChan

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
	go fw.tailFile(path, tailer, &curIndex, notifyChan)

	// Monitor for file rotation
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-fw.stopChan:
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
				case <-fw.stopChan:
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
					fw.notifyChan[path] = notifyChan
					go fw.tailFile(path, tailer, &curIndex, notifyChan)
					break
				}
				break
			}
		}
	}
}

func (fw *FileWatcher) tailFile(path string, tailer *tail.Tail, curIndex *string, notifyChan chan bool) {
	for {
		select {
		case line := <-tailer.Lines:
			if line == nil {
				continue
			}

			// Determine index based on current date
			year, month, day := time.Now().Date()
			indexNew := fmt.Sprintf("%s_%s_%d_%02d_%02d", fw.logstash.(*ZincLogstash).serviceName,
				strings.ReplaceAll(path, "/", "_"), year, month, day)

			// Create index if needed
			if *curIndex != indexNew {
				err := fw.createIndex(indexNew)
				if err != nil {
					logx.Errorf("Failed to create index %s: %v", indexNew, err)
					continue
				}
				*curIndex = indexNew
			}

			// Send log if not empty
			if line.Text != "" {
				if err := fw.logstash.SendLog(*curIndex, line.Text); err != nil {
					logx.Errorf("Failed to send log from %s: %v", path, err)
				}
			}

		case <-notifyChan:
			logx.Infof("Stopping tail for file: %s", path)
			return
		}
	}
}

func (fw *FileWatcher) createIndex(index string) error {
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

	url := fmt.Sprintf("%s/api/index", fw.logstash.(*ZincLogstash).baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(fw.logstash.(*ZincLogstash).username, fw.logstash.(*ZincLogstash).password)

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
