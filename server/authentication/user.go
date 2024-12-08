package authentication

import (
	"errors"
	"fmt"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type user interface {
	Username() string
	Email() string
	storeCookie(c echo.Context) error
	Login(c echo.Context) (bool, error)

	// Si se llama desde un usuario leido de las cookies, devolvera un error
	DeleteFromDB() error

	// Si se llama desde un usuario leido de las cookies, devolvera un error
	StoreInDB() error
}

// Constructor para crear nuevo usuario
func CreateUser(username, password, email string) user {
	return &newuser{
		username: username,
		password: password,
		email:    email,
	}
}

// Constructor para leer usuario de las cookies del cliente. Devuelve error si no existe usuario en cookies o si no se han podido leer
func ReadUserFromCookies(c echo.Context) (user, error) {
	user := &userFromCookies{}
	sess, err := session.Get("session", c)
	if err != nil {
		return user, fmt.Errorf("cookie session not obtainable: %v", err)
	}

	username, ok := sess.Values["username"].(string)
	if !ok {
		return user, errors.New("could not read valid user data in cookies")
	}

	email, ok := sess.Values["email"].(string)
	if !ok {
		return user, errors.New("could not read valid user data in cookies")
	}

	user.username = username
	user.email = email

	return user, nil
}

func Logout(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return fmt.Errorf("cookie session not obtainable: %v", err)
	}

	sess.Options.MaxAge = -1

	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return err
	}

	return nil
}
