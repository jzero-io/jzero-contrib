package lock

type Lock interface {
	Lock() error
	Unlock() error
}
