package logtoconsole

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/nxadm/tail"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
)

// Must
// Deprecated: As of go-zero v1.7.0, use [logx.AddWriter(logx.NewWriter(os.Stdout))]
func Must(l logx.LogConf) {
	mustLogToConsole(l)
}

func mustLogToConsole(l logx.LogConf) {
	if l.Mode == "console" {
		return
	}

	var logPaths []string

	err := filepath.Walk(l.Path, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			return nil
		}
		if filepath.Ext(f.Name()) == ".log" {
			logPaths = append(logPaths, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	go func() {
		err := mr.MapReduceVoid(func(source chan<- string) {
			for _, v := range logPaths {
				source <- v
			}
		}, func(item string, writer mr.Writer[*tail.Tail], cancel func(error)) {
			t, err := tail.TailFile(
				filepath.Join(item), tail.Config{
					Follow: true,
					ReOpen: true,
					Poll:   true,
					Location: &tail.SeekInfo{
						Offset: 0,
						Whence: io.SeekEnd,
					},
					Logger: tail.DiscardingLogger,
				})
			if err != nil {
				cancel(err)
			}

			// Print the text of each received line
			for line := range t.Lines {
				fmt.Println(line.Text)
			}
		}, func(pipe <-chan *tail.Tail, cancel func(error)) {}, mr.WithWorkers(len(logPaths)))
		if err != nil {
			return
		}
	}()
}
