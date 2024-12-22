package shared

import "fmt"

type Address string

var scheme string = "http"

var (
	FrontendAddr        Address = "127.0.0.1:8080"
	AuthenticationAddr  Address = "127.0.0.1:8082"
	DatabaseAddr        Address = "127.0.0.1:8083"
	MediaListAddr       Address = "127.0.0.1:8084"
	MovieAddr           Address = "127.0.0.1:8085"
	PasswordHashingAddr Address = "127.0.0.1:8086"
	RegistrationAddr    Address = "127.0.0.1:8087"
	SongAddr            Address = "127.0.0.1:8088"
)

// El Path debe empezar con un "/", por ejemplo, "/query" all llamar DatabaseAddr para realizar una consulta a la base de datos
func (a *Address) WithSchemeAndPath(path string) string {
	return fmt.Sprintf("%v://%v%v", scheme, *a, path)
}

func (a *Address) DomainPort() string {
	return string(*a)
}
