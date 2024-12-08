package server

import (
	"errors"
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/vicentebravocabezas/streamingUIDE/server/authentication"
	"github.com/vicentebravocabezas/streamingUIDE/server/media"
	"github.com/vicentebravocabezas/streamingUIDE/web/templates"
	"github.com/vicentebravocabezas/streamingUIDE/web/templates/stream"
)

func render(c echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(c.Request().Context(), buf); err != nil {
		return err
	}

	return c.HTML(statusCode, buf.String())
}

func index(c echo.Context) error {
	user, err := authentication.ReadUserFromCookies(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	if authorized, _ := user.Login(c); !authorized {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	return c.Redirect(http.StatusSeeOther, "/stream")
}

func loginScreen(c echo.Context) error {
	return render(c, http.StatusOK, templates.Layout(templates.Login(nil)))
}

func login(c echo.Context) error {
	user := authentication.CreateUser(c.FormValue("username"), c.FormValue("password"), "")
	authorized, err := user.Login(c)
	if err != nil {
		return err
	}

	if !authorized {
		return render(c, http.StatusOK, templates.Layout(templates.Login(errors.New("Las credenciales no son correctas. Intente de nuevo."))))
	}

	return c.Redirect(http.StatusSeeOther, "/stream")
}

func streamPage(c echo.Context) error {
	return render(c, http.StatusOK, templates.Layout(stream.Stream()))
}

func signup(c echo.Context) error {
	return render(c, http.StatusOK, templates.Layout(templates.SignUp()))
}

// control para registrar usuario
func insertUser(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	email := c.FormValue("email")

	user := authentication.CreateUser(username, password, email)
	user.StoreInDB()

	users := authentication.QueryUsers()
	return c.JSON(http.StatusOK, users)
}

// control para borrar usuario
func deleteUser(c echo.Context) error {
	username := c.QueryParam("username")
	password := c.QueryParam("password")
	email := c.QueryParam("email")

	user := authentication.CreateUser(username, password, email)
	user.DeleteFromDB()

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
