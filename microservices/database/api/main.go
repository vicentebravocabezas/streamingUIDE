package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/vicentebravocabezas/streamingUIDE/microservices/database"
	"github.com/vicentebravocabezas/streamingUIDE/microservices/shared"
)

func main() {
	e := echo.New()

	// mostrar logs en la consola
	e.Use(middleware.Logger())

	e.POST("/execute", func(c echo.Context) error {
		var queryParams shared.DatabaseQuery
		if err := c.Bind(&queryParams); err != nil {
			return err
		}

		result := make(map[string]any)
		ok, err := database.ExecuteQueryNoRows(&queryParams)
		result["success"] = ok
		if err != nil {
			result["detail"] = err.Error()
			return c.JSON(http.StatusBadRequest, result)
		}

		return c.JSON(http.StatusOK, result)
	})

	e.POST("/query", func(c echo.Context) error {
		var queryParams shared.DatabaseQuery
		if err := c.Bind(&queryParams); err != nil {
			return err
		}

		result, err := database.ExecuteQuery(&queryParams)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, result)
	})

	e.Logger.Fatal(e.Start(shared.DatabaseAddr.DomainPort()))
}
