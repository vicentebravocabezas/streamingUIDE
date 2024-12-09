package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vicentebravocabezas/streamingUIDE/server/authentication"
	"github.com/vicentebravocabezas/streamingUIDE/server/media"
	"github.com/vicentebravocabezas/streamingUIDE/web/templates"
)

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

func signup(c echo.Context) error {
	return render(c, http.StatusOK, templates.Layout(templates.SignUp(nil)))
}

func logout(c echo.Context) error {
	if err := authentication.Logout(c); err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/login")
}

// control para registrar usuario
func registerUser(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	email := c.FormValue("email")

	user := authentication.CreateUser(username, password, email)
	if err := user.StoreInDB(); err != nil {
		return render(c, http.StatusOK, templates.Layout(templates.SignUp(fmt.Errorf("Error: %v", err))))
	}

	return c.Redirect(http.StatusSeeOther, "/login")
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
	list, err := media.MovieList()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, list)
}
