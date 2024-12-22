package shared

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
)

type User interface {
	JSONReader() *bytes.Reader
	GetUsername() string
	GetEmail() string
	GetPassword() string
	GetPasswordHash() []byte
	QueryPasswordHash() error
	CheckInDB() (bool, error)
	DeleteFromDB() error
	StoreInDB() error
}

func GetUserFromCookies() User {
	return &UserFromCookies{}
}

func GetUser() User {
	return &Newuser{}
}

type UserFromCookies struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (u *UserFromCookies) JSONReader() *bytes.Reader {
	marshalled, _ := json.Marshal(u)
	return bytes.NewReader(marshalled)
}

func (u *UserFromCookies) GetUsername() string {
	return u.Username
}

func (u *UserFromCookies) GetEmail() string {
	return u.Email
}

func (u *UserFromCookies) GetPassword() string {
	return ""
}

func (u *UserFromCookies) QueryPasswordHash() error {
	return errors.New("cannot call from this struct")
}

func (u *UserFromCookies) GetPasswordHash() []byte {
	return nil
}

func (u *UserFromCookies) StoreInDB() error {
	return errors.New("cannot call from this struct")
}

func (u *UserFromCookies) CheckInDB() (bool, error) {
	query := ConstructQuery("SELECT username FROM users WHERE username = ?", u.Username).JSONReader()

	resp, err := http.Post(DatabaseAddr.WithSchemeAndPath("/query"), "application/json", query)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var result []map[string]any

	body, _ := io.ReadAll(resp.Body)

	if err := json.Unmarshal(body, &result); err != nil {
		return false, err
	}

	if !slices.ContainsFunc(result, func(r map[string]any) bool {
		return r["username"] == u.Username
	}) {
		return false, nil
	}

	return true, nil
}

func (u *UserFromCookies) DeleteFromDB() error {
	return nil
}

type Newuser struct {
	Username       string `json:"username" form:"username"`
	Password       string `json:"password" form:"password"`
	HashedPassword []byte `json:"hashedPassword"`
	Email          string `json:"email" form:"email"`
}

func (u *Newuser) JSONReader() *bytes.Reader {
	marshalled, _ := json.Marshal(u)
	return bytes.NewReader(marshalled)
}

// metodo para borrar usuario
func (u *Newuser) DeleteFromDB() error {
	query := ConstructQuery("DELETE FROM users WHERE username = ?", u.Username).JSONReader()

	resp, err := http.Post(DatabaseAddr.WithSchemeAndPath("/execute"), "application/json", query)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	result := make(map[string]any)

	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	if !result["success"].(bool) {
		return fmt.Errorf("operation gave the following error: %v", result["detail"])
	}

	return nil
}

func (u *Newuser) CheckInDB() (bool, error) {
	return false, nil
}

func (u *Newuser) generateHashedPassword() error {
	resp, err := http.Post(PasswordHashingAddr.WithSchemeAndPath("/create"), "application/json", u.JSONReader())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var decoded map[string]string

	if err := json.Unmarshal(body, &decoded); err != nil {
		return err
	}

	hashedPasswordBytes, ok := decoded["hash"]
	if !ok {
		return errors.New("failed to obtain hashed password")
	}

	u.HashedPassword = []byte(hashedPasswordBytes)

	return nil
}

func (u *Newuser) StoreInDB() error {
	if err := u.generateHashedPassword(); err != nil {
		return err
	}

	query := ConstructQuery("INSERT INTO users VALUES (?, ?, ?)", u.Username, string(u.HashedPassword), u.Email).JSONReader()

	resp, err := http.Post(DatabaseAddr.WithSchemeAndPath("/execute"), "application/json", query)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var decoded map[string]any

	if err := json.Unmarshal(body, &decoded); err != nil {
		return err
	}

	if !decoded["success"].(bool) {
		if strings.Contains(decoded["detail"].(string), "UNIQUE constraint failed") {
			return errors.New("user already registered")
		}
		return fmt.Errorf("query gave the following error: %v", decoded["detail"])
	}

	return nil
}

func (u *Newuser) GetUsername() string {
	return u.Username
}

func (u *Newuser) GetEmail() string {
	return u.Email
}

func (u *Newuser) GetPassword() string {
	return u.Password
}

func (u *Newuser) GetPasswordHash() []byte {
	return u.HashedPassword
}

func (u *Newuser) QueryPasswordHash() error {
	if u.HashedPassword != nil {
		return nil
	}

	query := ConstructQuery("SELECT password FROM users WHERE username = ?", u.Username).JSONReader()

	resp, err := http.Post(DatabaseAddr.WithSchemeAndPath("/query"), "application/json", query)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var decoded []map[string]string

	if err := json.Unmarshal(body, &decoded); err != nil {
		return err
	}

	if len(decoded) == 0 {
		return errors.New("Failed to retrieve hash from database")
	}

	u.HashedPassword = []byte(decoded[0]["password"])

	return nil
}

// func (u *newuser) Login(c echo.Context) (bool, error) {
// 	hashed := u.queryHashedPassword()
// 	u.hashedPassword = []byte(hashed)
//
// 	if err := bcrypt.CompareHashAndPassword([]byte(u.hashedPassword), []byte(u.password)); err != nil {
// 		return false, nil
// 	}
//
// 	if err := u.storeCookie(c); err != nil {
// 		return false, fmt.Errorf("could not send cookie to client: %v", err)
// 	}
//
// 	return true, nil
// }
//
// func (u *newuser) queryHashedPassword() string {
// 	row := database.DB().QueryRow(`SELECT password FROM users WHERE username = ?`, u.username)
//
// 	var hashedPassword string
//
// 	row.Scan(&hashedPassword)
//
// 	return hashedPassword
// }
//
// func (u *newuser) storeCookie(c echo.Context) error {
// 	sess, err := session.Get("session", c)
// 	if err != nil {
// 		return fmt.Errorf("session data not obtained: %v", err)
// 	}
//
// 	sess.Options = &sessions.Options{
// 		Path: "/",
// 		//La cookie tiene una duracion de 7 dias
// 		MaxAge:   86400 * 7,
// 		HttpOnly: true,
// 		Secure:   true,
// 		SameSite: http.SameSiteStrictMode,
// 	}
//
// 	sess.Values["username"] = u.username
// 	sess.Values["email"] = u.email
//
// 	if err := sess.Save(c.Request(), c.Response()); err != nil {
// 		return fmt.Errorf("error saving cookie: %v", err)
// 	}
//
// 	return nil
// }
