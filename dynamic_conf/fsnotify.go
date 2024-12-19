package dynamic_conf

import (
	"os"
	"path/filepath"

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
				if (event.Has(fsnotify.Write) || event.Has(fsnotify.Rename)) &&
					filepath.ToSlash(filepath.Clean(event.Name)) == filepath.Clean(filepath.ToSlash(lf.Path)) {
					logx.Infof("listen %s %s event", event.Name, event.Op)
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

	// see: https://github.com/fsnotify/fsnotify/issues/363
	if err := lf.Watcher.Add(filepath.Dir(lf.Path)); err != nil {
		return err
	}
	return nil
}

func (lf *FsNotify) Value() (string, error) {
	file, err := os.ReadFile(lf.Path)
	return string(file), err
}
