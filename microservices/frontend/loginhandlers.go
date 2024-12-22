package frontend

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vicentebravocabezas/streamingUIDE/microservices/frontend/cookies"
	"github.com/vicentebravocabezas/streamingUIDE/microservices/frontend/web/templates"
	"github.com/vicentebravocabezas/streamingUIDE/microservices/shared"
)

func index(c echo.Context) error {
	user, err := cookies.GetUser(c)
	if err != nil {
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

	if !decoded["authenticated"] {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	return c.Redirect(http.StatusSeeOther, "/stream")
}

func loginScreen(c echo.Context) error {
	return render(c, http.StatusOK, templates.Layout(templates.Login(nil)))
}

func login(c echo.Context) error {
	user := shared.Newuser{}

	if err := c.Bind(&user); err != nil {
		return err
	}

	resp, err := http.Post(shared.AuthenticationAddr.WithSchemeAndPath("/authenticate"), "application/json", user.JSONReader())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var decoded map[string]bool

	if err := json.Unmarshal(body, &decoded); err != nil {
		return err
	}

	if !decoded["authenticated"] {
		return render(c, http.StatusOK, templates.Layout(templates.Login(errors.New("Las credenciales no son correctas. Intente de nuevo."))))
	}

	if err := cookies.StoreCookie(c, &user); err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/stream")
}

func signup(c echo.Context) error {
	return render(c, http.StatusOK, templates.Layout(templates.SignUp(nil)))
}

func logout(c echo.Context) error {
	if err := cookies.Logout(c); err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/login")
}

// control para registrar usuario
func registerUser(c echo.Context) error {
	user := shared.Newuser{}
	if err := c.Bind(&user); err != nil {
		return err
	}

	if err := user.StoreInDB(); err != nil {
		return render(c, http.StatusOK, templates.Layout(templates.SignUp(fmt.Errorf("Error: %v", err))))
	}

	return c.Redirect(http.StatusSeeOther, "/login")
}
