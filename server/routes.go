package server

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/vicentebravocabezas/streamingUIDE/server/authentication"
)

func RegisterMiddleware(e *echo.Echo) {
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("SESSIONKEY")))))

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Add("Cache-Control", "no-cache")
			c.Response().Header().Add("Cache-Control", "private")

			return next(c)
		}
	})

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			user, err := authentication.ReadUserFromCookies(c)
			if err != nil && c.Path() == "/stream" {
				return c.Redirect(http.StatusSeeOther, "/login")
			}

			authorized, _ := user.Login(c)
			if !authorized && c.Path() == "/stream" {
				return c.Redirect(http.StatusSeeOther, "/login")
			}

			return next(c)
		}
	})
}

func RegisterRoutes(e *echo.Echo) {
	e.GET("/", index)

	e.GET("/signup", signup)

	e.GET("/login", loginScreen)

	e.POST("/login", login)

	// ruta para insertar usuario
	e.GET("/insert-user", insertUser)

	//ruta para borrar usuario
	e.GET("/delete-user", deleteUser)

	e.GET("/stream", streamPage)

	//ruta para obtener lista de archivos multimedia
	e.GET("/media", mediaList)

	e.Static("/public", "web/assets/public")
}
