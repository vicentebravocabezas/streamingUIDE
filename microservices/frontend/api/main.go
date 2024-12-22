package main

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/vicentebravocabezas/streamingUIDE/microservices/frontend"
	"github.com/vicentebravocabezas/streamingUIDE/microservices/shared"
)

func main() {
	godotenv.Load()

	if os.Getenv("SESSIONKEY") == "" {
		log.Fatal(errors.New("SESSIONKEY environment variable not configured"))
	}

	e := echo.New()

	frontend.RegisterMiddleware(e)

	frontend.RegisterRoutes(e)

	// mostrar logs en la consola
	e.Use(middleware.Logger())

	e.Logger.Fatal(e.Start(shared.FrontendAddr.DomainPort()))
}
