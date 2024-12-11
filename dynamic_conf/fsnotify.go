package dynamic_conf

import (
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/zeromicro/go-zero/core/configcenter/subscriber"
	"github.com/zeromicro/go-zero/core/logx"
)

var _ subscriber.Subscriber = (*FsNotify)(nil)

type FsNotify struct {
	Path string `json:"path"`
	*fsnotify.Watcher
}

func NewFsNotify(path string) (*FsNotify, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &FsNotify{
		Path:    path,
		Watcher: watcher,
	}, nil
}

func (lf *FsNotify) AddListener(listener func()) error {
	go func() {
		for {
			select {
			case event, ok := <-lf.Watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					listener()
				}
			case err, ok := <-lf.Watcher.Errors:
				if !ok {
					return
				}
				logx.Errorf("error: %v", err)
			}
		}
	}()

	if err := lf.Watcher.Add(lf.Path); err != nil {
		return err
	}
	return nil
}

func (lf *FsNotify) Value() (string, error) {
	file, err := os.ReadFile(lf.Path)
	return string(file), err
}
