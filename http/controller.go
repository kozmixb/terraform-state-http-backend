package http

import (
	"encoding/json"
	"io"
	"net/http"
	"terraform-state-http-backend/app"

	"github.com/labstack/echo/v4"
)

func Show(c echo.Context) error {
	result := app.Show(c.Param("group"), c.Param("key"))

	if result == "" {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": "Not Found",
		})
	}

	return toJson(result, c)
}

func Update(c echo.Context) error {
	var body []byte
	if c.Request().Body != nil {
		body, _ = io.ReadAll(c.Request().Body)
	}

	result := app.Update(c.Param("group"), c.Param("key"), string(body))

	return toJson(result, c)
}

func Lock(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "not available",
	})
}

func Unlock(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "not available",
	})
}

func toJson(data string, c echo.Context) error {
	var response map[string]interface{}
	err := json.Unmarshal([]byte(data), &response)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":     data,
			"technical": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response)
}
