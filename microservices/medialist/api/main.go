package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/vicentebravocabezas/streamingUIDE/microservices/shared"
)

func main() {
	e := echo.New()

	// mostrar logs en la consola
	e.Use(middleware.Logger())

	e.GET("/", func(c echo.Context) error {
		media, err := shared.MediaList()
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, media)
	})

	e.Logger.Fatal(e.Start(shared.MediaListAddr.DomainPort()))
}
