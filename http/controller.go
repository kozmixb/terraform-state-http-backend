package http

import (
	"errors"
	"io"
	"net/http"
	"terraform-state-http-backend/app"
	"terraform-state-http-backend/app/drivers"

	"github.com/labstack/echo/v4"
)

func Show(c echo.Context) error {
	result, err := app.Show(c.Param("group"), c.Param("key"))

	if errors.Is(err, drivers.ErrNotFound) {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": "Not Found",
		})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.Blob(http.StatusOK, "application/json", result)
}

func Update(c echo.Context) error {
	var body []byte
	if c.Request().Body != nil {
		var err error
		body, err = io.ReadAll(c.Request().Body)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
			})
		}
	}

	result, err := app.Update(c.Param("group"), c.Param("key"), body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.Blob(http.StatusOK, "application/json", result)
}

func Lock(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	result, err := app.Lock(c.Param("group"), c.Param("key"), body)
	if errors.Is(err, drivers.ErrLocked) {
		return c.Blob(http.StatusLocked, "application/json", result)
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.Blob(http.StatusOK, "application/json", result)
}

func Unlock(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	err = app.Unlock(c.Param("group"), c.Param("key"), body)
	if errors.Is(err, drivers.ErrUnlockMismatch) {
		return c.JSON(http.StatusConflict, map[string]interface{}{
			"error": err.Error(),
		})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.NoContent(http.StatusOK)
}
