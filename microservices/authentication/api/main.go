package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/vicentebravocabezas/streamingUIDE/microservices/shared"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())

	e.POST("/check", func(c echo.Context) error {
		userData := shared.GetUserFromCookies()

		if err := c.Bind(&userData); err != nil {
			return err
		}

		verified, err := userData.CheckInDB()
		if err != nil {
			return err
		}

		if !verified {
			return c.JSON(http.StatusUnauthorized, map[string]bool{"authenticated": false})
		}

		return c.JSON(http.StatusOK, map[string]bool{"authenticated": true})
	})

	e.POST("/authenticate", func(c echo.Context) error {
		user := shared.GetUser()

		if err := c.Bind(&user); err != nil {
			return err
		}

		if err := user.QueryPasswordHash(); err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]bool{"authenticated": false})
		}

		resp, err := http.Post(shared.PasswordHashingAddr.WithSchemeAndPath("/verify"), "application/json", user.JSONReader())
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		var result map[string]bool

		if err := json.Unmarshal(body, &result); err != nil {
			return err
		}

		if !result["match"] {
			return c.JSON(http.StatusUnauthorized, map[string]bool{"authenticated": false})
		}

		return c.JSON(http.StatusOK, map[string]bool{"authenticated": true})
	})

	e.Logger.Fatal(e.Start(shared.AuthenticationAddr.DomainPort()))
}
