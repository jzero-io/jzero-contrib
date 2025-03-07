package filelock

import (
	"os"
	"syscall"
)

type FileLock struct {
	f *os.File
}

func New(fp string) (*FileLock, error) {
	f, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	return &FileLock{f: f}, nil
}

func (l *FileLock) Lock() error {
	err := syscall.Flock(int(l.f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		return err
	}
	return nil
}

func (l *FileLock) Unlock() error {
	defer l.f.Close()
	return syscall.Flock(int(l.f.Fd()), syscall.LOCK_UN)
}
