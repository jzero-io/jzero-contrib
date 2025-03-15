package filelock

import (
	"os"
	"sync"
	"testing"

	"github.com/spf13/cast"
)

func TestFileLock(t *testing.T) {
	wg := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			flock, _ := New("testdata/config")
			_ = flock.Lock()
			defer func(flock *FileLock) {
				_ = flock.Unlock()
			}(flock)

			f, _ := os.OpenFile("testdata/config", os.O_APPEND|os.O_WRONLY, 0o666)
			defer f.Close()
			_, _ = f.WriteString("test" + cast.ToString(num) + "\n")
		}(i)
	}
	wg.Wait()
}
