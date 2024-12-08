package authentication

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/vicentebravocabezas/streamingUIDE/database"
	"golang.org/x/crypto/bcrypt"
)

type newuser struct {
	username       string
	password       string
	hashedPassword []byte
	email          string
}

func QueryUsers() []newuser {
	row, err := database.OpenDB().Query("SELECT * FROM users")
	if err != nil {
		log.Println(err)
	}

	var results []newuser

	for row.Next() {
		var username string
		var password string
		var email string
		row.Scan(&username, &password, &email)
		results = append(results, newuser{
			username: username,
			password: password,
			email:    email,
		})
	}

	return results
}

// metodo para borrar usuario
func (u *newuser) DeleteFromDB() error {
	result, err := database.OpenDB().Exec("DELETE FROM users WHERE username = ?", u.username)
	if err != nil {
		return fmt.Errorf("could not delete user: %v", err)
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected < 1 {
		return errors.New("no user deleted from the db")
	}

	return nil
}

// metodo para almacenar nuevo usuario
func (u *newuser) StoreInDB() error {
	var err error
	//hash de la contraseña para no almacenarla como texto
	u.hashedPassword, err = bcrypt.GenerateFromPassword([]byte(u.password), 10)
	if err != nil {
		fmt.Println(err)
	}

	if _, err := database.OpenDB().Exec("INSERT INTO users VALUES (?, ?, ?)", u.username, string(u.hashedPassword), u.email); err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return errors.New("usuario ya se encuentra registrado")
		}
		return err
	}

	return nil
}

func (u *newuser) Username() string {
	return u.username
}

func (u *newuser) Email() string {
	return u.email
}

func (u *newuser) Login(c echo.Context) (bool, error) {
	hashed := u.queryHashedPassword()
	u.hashedPassword = []byte(hashed)

	if err := bcrypt.CompareHashAndPassword([]byte(u.hashedPassword), []byte(u.password)); err != nil {
		return false, nil
	}

	if err := u.storeCookie(c); err != nil {
		return false, fmt.Errorf("could not send cookie to client: %v", err)
	}

	return true, nil
}

func (u *newuser) queryHashedPassword() string {
	row := database.OpenDB().QueryRow(`SELECT password FROM users WHERE username = ?`, u.username)

	var hashedPassword string

	row.Scan(&hashedPassword)

	return hashedPassword
}

func (u *newuser) storeCookie(c echo.Context) error {
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

	sess.Values["username"] = u.username
	sess.Values["email"] = u.email

	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return fmt.Errorf("error saving cookie: %v", err)
	}

	return nil
}
