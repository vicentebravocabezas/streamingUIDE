package cookies

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/vicentebravocabezas/streamingUIDE/microservices/shared"
)

func GetUser(c echo.Context) (shared.User, error) {
	user := &shared.UserFromCookies{}
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

	user.Username = username
	user.Email = email

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

func StoreCookie(c echo.Context, u shared.User) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return fmt.Errorf("session data not obtained: %v", err)
	}

	sess.Options = &sessions.Options{
		Path: "/",
		//La cookie tiene una duracion de 7 dias
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	sess.Values["username"] = u.GetUsername()
	sess.Values["email"] = u.GetEmail()

	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return fmt.Errorf("error saving cookie: %v", err)
	}

	return nil
}
