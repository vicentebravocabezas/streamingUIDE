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

	e.POST("/", func(c echo.Context) error {
		user := shared.GetUser()

		if err := c.Bind(&user); err != nil {
			return err
		}

		if user.GetUsername() == "" || user.GetPassword() == "" || user.GetEmail() == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "no field can be empty"})
		}

		if err := user.StoreInDB(); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}

		return c.JSON(http.StatusCreated, map[string]string{"message": "registration successful"})
	})

	e.DELETE("/", func(c echo.Context) error {
		user := shared.GetUser()

		if err := c.Bind(&user); err != nil {
			return err
		}

		if err := user.DeleteFromDB(); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "user deleted"})
	})

	e.Logger.Fatal(e.Start(shared.RegistrationAddr.DomainPort()))
}
