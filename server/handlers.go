package server

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/vicentebravocabezas/streamingUIDE/server/authentication"
	"github.com/vicentebravocabezas/streamingUIDE/server/media"
	"github.com/vicentebravocabezas/streamingUIDE/web/templates"
)

func render(c echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(c.Request().Context(), buf); err != nil {
		return err
	}

	return c.HTML(statusCode, buf.String())
}

// TODO: implementar pagina de inicio
func index(c echo.Context) error {
	return render(c, http.StatusOK, templates.Index())
}

// control para registrar usuario
func insertUser(c echo.Context) error {
	username := c.QueryParam("username")
	password := c.QueryParam("password")
	email := c.QueryParam("email")

	user := authentication.CreateUser(username, password, email)
	user.Store()

	users := authentication.QueryUsers()
	return c.JSON(http.StatusOK, users)
}

// control para borrar usuario
func deleteUser(c echo.Context) error {
	username := c.QueryParam("username")
	password := c.QueryParam("password")
	email := c.QueryParam("email")

	user := authentication.CreateUser(username, password, email)
	user.Delete()

	users := authentication.QueryUsers()
	return c.JSON(http.StatusOK, users)
}

// control para lista de multimedia registrada
func mediaList(c echo.Context) error {
	list, err := media.MediaList()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, list)
}
