package frontend

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/vicentebravocabezas/streamingUIDE/microservices/frontend/cookies"
	"github.com/vicentebravocabezas/streamingUIDE/microservices/shared"
)

func RegisterMiddleware(e *echo.Echo) {
	//middleware para establecer manejo de cookies
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("SESSIONKEY")))))

	//middleware para control de cache
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Add("Cache-Control", "no-cache")
			c.Response().Header().Add("Cache-Control", "private")

			return next(c)
		}
	})

	//middleware para autenticacion por cookies
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, err := cookies.GetUser(c)
			if err != nil && c.Path() == "/stream" {
				return c.Redirect(http.StatusSeeOther, "/login")
			}

			resp, err := http.Post(shared.AuthenticationAddr.WithSchemeAndPath("/check"), "application/json", user.JSONReader())
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)

			var decoded map[string]bool

			if err := json.Unmarshal(body, &decoded); err != nil {
				return err
			}

			if !decoded["authenticated"] && c.Path() == "/stream" {
				return c.Redirect(http.StatusSeeOther, "/login")
			}

			return next(c)
		}
	})
}

func RegisterRoutes(e *echo.Echo) {
	e.GET("/", index)

	//rutas para autenticacion y registro de usuario
	e.GET("/signup", signup)

	e.POST("/signup", registerUser)

	e.GET("/login", loginScreen)

	e.POST("/login", login)

	e.GET("/logout", logout)

	e.GET("/stream", streamPage)

	e.Static("/public", "web/assets/public")
}

func render(c echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(c.Request().Context(), buf); err != nil {
		return err
	}

	return c.HTML(statusCode, buf.String())
}
