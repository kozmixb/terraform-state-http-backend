package app

import (
	"os"
	"terraform-state-http-backend/app/drivers"
)

func Show(group string, key string) string {
	driver := getDriver()

	if !driver.Exists(group, key) {
		return ""
	}

	return driver.Show(group, key)
}

func Update(group string, key string, payload string) string {
	driver := getDriver()

	return driver.Update(group, key, payload)

}

func Lock(group string, key string) {
	driver := getDriver()
	driver.Exists(group, key)
}

func Unlock(group string, key string) {
	driver := getDriver()
	driver.Exists(group, key)
}

func getDriver() drivers.Driver {
	switch getDefault() {
	// case "mysql":
	// 	return drivers.MysqlDriver{}

	case "sqlite":
		return drivers.NewSqLiteDriver()
	}

	return drivers.FileDriver{}
}

func getDefault() string {
	driver := os.Getenv("DRIVER")

	if driver == "" {
		return "sqlite"
	}

	return driver
}
