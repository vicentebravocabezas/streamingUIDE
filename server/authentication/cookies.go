package authentication

import (
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/vicentebravocabezas/streamingUIDE/database"
)

type userFromCookies struct {
	username string
	email    string
}

func (u *userFromCookies) Username() string {
	return u.username
}

func (u *userFromCookies) Email() string {
	return u.email
}

func (u *userFromCookies) Login(c echo.Context) (bool, error) {
	return u.checkUserInDB()
}

func (u *userFromCookies) DeleteFromDB() error {
	return errors.New("cannot execute method, must be a new user not read from browser cookies")
}
func (u *userFromCookies) StoreInDB() error {
	return errors.New("cannot execute method, must be a new user not read from browser cookies")
}

func (u *userFromCookies) checkUserInDB() (bool, error) {
	row := database.OpenDB().QueryRow("SELECT username FROM users WHERE username = ?", u.username)
	var username string
	row.Scan(&username)

	return username != "", nil
}

func (u *userFromCookies) storeCookie(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return err
	}

	sess.Options = &sessions.Options{
		Path: "/",
		//La cookie tiene una duracion de 7 dias
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	sess.Values["username"] = u.username
	sess.Values["email"] = u.email

	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return err
	}

	return nil
}
