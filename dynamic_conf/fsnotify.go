package dynamic_conf

import (
	"os"
	"path/filepath"

	"github.com/a8m/envsubst"
	"github.com/eddieowens/opts"
	"github.com/fsnotify/fsnotify"
	"github.com/jaronnie/genius"
	"github.com/zeromicro/go-zero/core/configcenter/subscriber"
	"github.com/zeromicro/go-zero/core/logx"
)

var _ subscriber.Subscriber = (*FsNotify)(nil)

type FsNotify struct {
	path string

	// options
	options FsNotifyOpts

	*fsnotify.Watcher
}

type FsNotifyOpts struct {
	UseEnv bool

	UseKey string
}

func (opts FsNotifyOpts) DefaultOptions() FsNotifyOpts {
	return FsNotifyOpts{}
}

func WithUseEnv(useEnv bool) opts.Opt[FsNotifyOpts] {
	return func(o *FsNotifyOpts) {
		o.UseEnv = useEnv
	}
}

func WithUseKey(useKey string) opts.Opt[FsNotifyOpts] {
	return func(o *FsNotifyOpts) {
		o.UseKey = useKey
	}
}

func NewFsNotify(path string, op ...opts.Opt[FsNotifyOpts]) (*FsNotify, error) {
	o := opts.DefaultApply(op...)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &FsNotify{
		path:    path,
		Watcher: watcher,
		options: o,
	}, nil
}

func (f *FsNotify) AddListener(listener func()) error {
	go func() {
		for {
			select {
			case event, ok := <-f.Watcher.Events:
				if !ok {
					return
				}
				if (event.Has(fsnotify.Write) || event.Has(fsnotify.Rename)) &&
					filepath.ToSlash(filepath.Clean(event.Name)) == filepath.Clean(filepath.ToSlash(f.path)) {
					logx.Infof("listen %s %s event", event.Name, event.Op)
					listener()
				}
			case err, ok := <-f.Watcher.Errors:
				if !ok {
					return
				}
				logx.Errorf("error: %v", err)
			}
		}
	}()

	// see: https://github.com/fsnotify/fsnotify/issues/363
	if err := f.Watcher.Add(filepath.Dir(f.path)); err != nil {
		return err
	}
	return nil
}

func (f *FsNotify) Value() (string, error) {
	file, err := os.ReadFile(f.path)
	if err != nil {
		return "", err
	}

	if f.options.UseEnv {
		file, err = envsubst.Bytes(file)
		if err != nil {
			return "", err
		}
	}

	if f.options.UseKey != "" {
		g, err := genius.NewFromType(file, filepath.Ext(f.path))
		if err != nil {
			return "", err
		}
		sub := g.Sub(f.options.UseKey)
		file, err = sub.EncodeToType(filepath.Ext(f.path))
		if err != nil {
			return "", err
		}
	}

	return string(file), nil
}
