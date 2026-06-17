package drivers

import "errors"

var (
	ErrNotFound       = errors.New("state not found")
	ErrLocked         = errors.New("state is locked")
	ErrUnlockMismatch = errors.New("lock id does not match")
)

type Driver interface {
	Exists(group string, key string) (bool, error)
	Show(group string, key string) ([]byte, error)
	Update(group string, key string, payload []byte) ([]byte, error)
	Lock(group string, key string, payload []byte) ([]byte, error)
	Unlock(group string, key string, payload []byte) error
}
