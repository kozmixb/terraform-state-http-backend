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
	e.Use(middleware.Logger())

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
