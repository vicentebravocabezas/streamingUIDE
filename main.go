package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/vicentebravocabezas/streamingUIDE/server"
)

func main() {
	godotenv.Load()

	e := echo.New()

	server.RegisterMiddleware(e)

	server.RegisterRoutes(e)

	e.Use(middleware.Logger())

	e.Logger.Fatal(e.Start(":8090"))
}
