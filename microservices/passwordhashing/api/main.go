package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/vicentebravocabezas/streamingUIDE/microservices/shared"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	e := echo.New()

	// mostrar logs en la consola
	e.Use(middleware.Logger())

	e.POST("/create", func(c echo.Context) error {
		user := shared.GetUser()

		if err := c.Bind(&user); err != nil {
			return err
		}

		if user.GetPassword() == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "no data received"})
		}

		hashBytes, err := bcrypt.GenerateFromPassword([]byte(user.GetPassword()), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, map[string]string{"hash": string(hashBytes)})
	})

	e.POST("/verify", func(c echo.Context) error {
		user := shared.GetUser()

		if err := c.Bind(&user); err != nil {
			return err
		}

		if err := bcrypt.CompareHashAndPassword(user.GetPasswordHash(), []byte(user.GetPassword())); err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]bool{"match": false})
		}

		return c.JSON(http.StatusOK, map[string]bool{"match": true})
	})

	e.Logger.Fatal(e.Start(shared.PasswordHashingAddr.DomainPort()))
}
