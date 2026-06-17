package drivers

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type FileDriver struct{}

func (f FileDriver) Exists(group string, key string) (bool, error) {
	path := getStatePath(group, key)

	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return err == nil, err
}

func (f FileDriver) Show(group string, key string) ([]byte, error) {
	path := getStatePath(group, key)
	content, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrNotFound
	}

	return content, err
}

func (f FileDriver) Update(group string, key string, payload []byte) ([]byte, error) {
	path := getStatePath(group, key)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}

	if err := os.WriteFile(path, payload, 0666); err != nil {
		return nil, err
	}

	return payload, nil
}

func (f FileDriver) Lock(group string, key string, payload []byte) ([]byte, error) {
	path := getLockPath(group, key)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	if errors.Is(err, os.ErrExist) {
		lock, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil, readErr
		}

		return lock, ErrLocked
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if _, err := file.Write(payload); err != nil {
		return nil, err
	}

	return payload, nil
}

func (f FileDriver) Unlock(group string, key string, payload []byte) error {
	path := getLockPath(group, key)
	lock, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	if !sameLockID(lock, payload) {
		return ErrUnlockMismatch
	}

	return os.Remove(path)
}

func getStatePath(group string, key string) string {
	return filepath.Join("storage", fmt.Sprintf("%s-%s.json", group, key))
}

func getLockPath(group string, key string) string {
	return filepath.Join("storage", fmt.Sprintf("%s-%s.lock", group, key))
}
