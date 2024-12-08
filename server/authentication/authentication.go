package authentication

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/vicentebravocabezas/streamingUIDE/database"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	Username       string
	Password       string
	HashedPassword []byte
	Email          string
}

func QueryUsers() []user {
	row, err := database.OpenDB().Query("SELECT * FROM users")
	if err != nil {
		log.Println(err)
	}

	var results []user

	for row.Next() {
		var username string
		var password string
		var email string
		row.Scan(&username, &password, &email)
		results = append(results, user{
			Username: username,
			Password: password,
			Email:    email,
		})
	}

	return results
}

func ReadUserFromCookies(c echo.Context) (*user, error) {
	sess, err := session.Get("session", c)
	if err != nil {
		return nil, err
	}

	username, ok := sess.Values["usuario"].(string)
	if !ok {
		return nil, errors.New("could not read valid user data in cookies")
	}

	return &user{
		Username: username,
	}, nil
}

// Constructor para crear usuario
func CreateUser(username, password, email string) *user {
	return &user{
		Username: username,
		Password: password,
		Email:    email,
	}
}

// metodo para borrar usuario
func (u *user) DeleteFromDB() {
	database.OpenDB().Exec("DELETE FROM users WHERE username = ?", u.Username)
}

// metodo para almacenar nuevo usuario, requiere el email del usuario
func (u *user) StoreInDB() error {
	var err error
	//hash de la contraseña para no almacenarla como texto
	u.HashedPassword, err = bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		fmt.Println(err)
	}

	if u.Email == "" {
		return errors.New("email not provided")
	}

	if _, err := database.OpenDB().Exec("INSERT INTO users VALUES (?, ?, ?)", u.Username, string(u.HashedPassword), u.Email); err != nil {
		return err
	}

	return nil
}

func (u *user) Login(c echo.Context) (bool, error) {
	hashed := u.queryHashedPassword()
	u.HashedPassword = []byte(hashed)

	if err := bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(u.Password)); err != nil {
		return false, nil
	}

	if err := u.storeCookie(c); err != nil {
		return false, err
	}

	return true, nil
}

func (u *user) queryHashedPassword() string {
	row := database.OpenDB().QueryRow(`SELECT password FROM users WHERE username = ?`, u.Username)

	var hashedPassword string

	row.Scan(&hashedPassword)

	return hashedPassword
}

func (u *user) storeCookie(c echo.Context) error {
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

	sess.Values["usuario"] = u.Username

	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return err
	}

	return nil
}
