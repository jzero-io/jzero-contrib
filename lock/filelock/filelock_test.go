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
			err := flock.Lock()
			if err != nil {
				t.Error(err)
			}
			defer flock.Unlock()

			f, _ := os.OpenFile("testdata/config", os.O_APPEND|os.O_WRONLY, 0666)
			defer f.Close()
			f.WriteString("test" + cast.ToString(num) + "\n")
		}(i)
	}
	wg.Wait()
}
