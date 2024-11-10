package main

import (
	"github.com/labstack/echo/v4"
	"github.com/vicentebravocabezas/streamingUIDE/server"
)

func main() {
	e := echo.New()

	server.RegisterRoutes(e)

	e.Start(":8080")
}
