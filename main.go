package main

import (
	"fmt"
	"os"
	"terraform-state-http-backend/routes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${method} ${uri}, (${latency_human})\n",
	}))
	useBasicAuth(e)

	routes.Load(e)

	e.Logger.Fatal(e.Start(getPort()))
}

func getPort() string {
	port := os.Getenv("HTTP_PORT")

	if port == "" {
		return ":8080"
	}

	return fmt.Sprintf(":%s", port)
}

func useBasicAuth(e *echo.Echo) {
	username := os.Getenv("BASIC_AUTH_USERNAME")
	password := os.Getenv("BASIC_AUTH_PASSWORD")
	if username == "" || password == "" {
		return
	}

	e.Use(middleware.BasicAuth(func(inputUsername string, inputPassword string, c echo.Context) (bool, error) {
		return inputUsername == username && inputPassword == password, nil
	}))
}
