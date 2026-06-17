package app

import (
	"os"
	"sync"
	"terraform-state-http-backend/app/drivers"
)

var (
	driver     drivers.Driver
	driverLock sync.Mutex
)

func Show(group string, key string) ([]byte, error) {
	driver := getDriver()

	exists, err := driver.Exists(group, key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, drivers.ErrNotFound
	}

	return driver.Show(group, key)
}

func Update(group string, key string, payload []byte) ([]byte, error) {
	driver := getDriver()

	return driver.Update(group, key, payload)
}

func Lock(group string, key string, payload []byte) ([]byte, error) {
	driver := getDriver()
	return driver.Lock(group, key, payload)
}

func Unlock(group string, key string, payload []byte) error {
	driver := getDriver()
	return driver.Unlock(group, key, payload)
}

func getDriver() drivers.Driver {
	driverLock.Lock()
	defer driverLock.Unlock()

	if driver != nil {
		return driver
	}

	switch getDefault() {
	// case "mysql":
	// 	return drivers.MysqlDriver{}

	case "sqlite":
		driver = drivers.NewSqLiteDriver()
	default:
		driver = drivers.FileDriver{}
	}

	return driver
}

func getDefault() string {
	driver := os.Getenv("DRIVER")

	if driver == "" {
		return "file"
	}

	return driver
}
