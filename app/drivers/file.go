package drivers

import (
	"errors"
	"fmt"
	"os"
)

type FileDriver struct{}

func (f FileDriver) Exists(group string, key string) bool {
	path := getStatePath(group, key)

	_, error := os.Stat(path)
	return !errors.Is(error, os.ErrNotExist)
}

func (f FileDriver) Show(group string, key string) string {
	path := getStatePath(group, key)
	content, _ := os.ReadFile(path)

	return string(content)
}

func (f FileDriver) Update(group string, key string, payload string) string {
	path := getStatePath(group, key)
	os.WriteFile(path, []byte(payload), 0666)

	return payload
}

func getStatePath(group string, key string) string {
	return fmt.Sprintf("storage/%s-%s.json", group, key)

}
