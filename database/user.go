package database

type user struct {
	username string
	password string
	email    string
}

func User() *user {
	return &user{}
}

func (u *user) Test() {

}
