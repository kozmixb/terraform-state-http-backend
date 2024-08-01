package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
)

func showConfig(c echo.Context) error {
	return c.JSON(200, map[string]interface{}{
		"message": "not available",
	})
}

func updateConfig(c echo.Context) error {
	return c.JSON(200, map[string]interface{}{
		"message": "not available",
	})
}

func lockConfig(c echo.Context) error {

	return c.JSON(200, map[string]interface{}{
		"message": "not available",
	})
}

func unlockConfig(c echo.Context) error {

	return c.JSON(200, map[string]interface{}{
		"message": "not available",
	})
}

func routes(e *echo.Echo) {
	e.GET("/:group/:key", showConfig)
	e.POST("/:group/:key", updateConfig)
	e.PUT("/:group/:key", lockConfig)
	e.DELETE("/:group/:key", unlockConfig)
}

func getStatePath(group string, key string) string {
	return fmt.Sprintf("storage/%s-%s.json", group, key)

}

func isStateExists(path string) bool {
	_, error := os.Stat(path)
	return !errors.Is(error, os.ErrNotExist)
}
