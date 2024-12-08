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

			if _, err := authentication.ReadUserFromCookies(c); err != nil && c.Path() == "/media" {
				return c.Redirect(http.StatusPermanentRedirect, "/login")
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
