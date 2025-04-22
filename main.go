package main

import (
	"quickprompt/config"
	"quickprompt/handlers"

	"github.com/labstack/echo/v4"
)

func main() {
	config.LoadEnv()

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Welcome to QuickPrompt API 🚀")
	})

	e.GET("/ping", func(c echo.Context) error {
		return c.String(200, "pong!")
	})

	// Optional: add GET for /api/prompt for clarity
	e.GET("/api/prompt", func(c echo.Context) error {
		return c.String(200, "Use POST method with a prompt payload.")
	})

	e.POST("/api/prompt", handlers.HandlePrompt)

	e.Logger.Fatal(e.Start(":8080"))
}
