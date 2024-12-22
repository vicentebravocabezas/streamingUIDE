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

	e.GET("/:id", func(c echo.Context) error {
		movie, err := shared.GetMusic(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"message": "Not Found"})
		}

		return c.JSON(http.StatusOK, movie)
	})

	e.Logger.Fatal(e.Start(shared.SongAddr.DomainPort()))
}
