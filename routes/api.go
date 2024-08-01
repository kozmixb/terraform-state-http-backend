package routes

import (
	"terraform-state-http-backend/http"

	"github.com/labstack/echo/v4"
)

func Load(e *echo.Echo) {
	e.GET("/:group/:key", http.Show)
	e.POST("/:group/:key", http.Update)
	e.PUT("/:group/:key", http.Lock)
	e.DELETE("/:group/:key", http.Unlock)
}
