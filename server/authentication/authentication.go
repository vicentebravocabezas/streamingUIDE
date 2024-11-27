package authentication

import (
	"fmt"
	"log"

	"github.com/vicentebravocabezas/streamingUIDE/database"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string
	Password string
	Email    string
}

func QueryUsers() []User {
	row, err := database.OpenDB().Query("SELECT * FROM users")
	if err != nil {
		log.Println(err)
	}

	var results []User

	for row.Next() {
		var user string
		var password string
		var email string
		row.Scan(&user, &password, &email)
		results = append(results, User{
			Username: user,
			Password: password,
			Email:    email,
		})
	}

	return results
}

// Constructor para crear usuario
func CreateUser(username, password, email string) *User {
	//hash de la contraseña para no almacenarla como texto
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		fmt.Println(err)
	}

	return &User{
		Username: username,
		Password: string(hash),
		Email:    email,
	}
}

// metodo para borrar usuario
func (u *User) Delete() {
	database.OpenDB().Exec("DELETE FROM users WHERE username = ?", u.Username)
}

// metodo para almacenar nuevo usuario
func (u *User) Store() {
	database.OpenDB().Exec("INSERT INTO users VALUES (?, ?, ?)", u.Username, u.Password, u.Email)
}
