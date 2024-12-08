package tokens

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"

	"github.com/google/uuid"
)

type jwtKey struct {
	expirationDate time.Time
	token          string
}

func Generate(jwtSigningKey []byte) (jwtKey, error) {
	accessExpirationDate := time.Now().Add(1*time.Hour + time.Duration(rand.Intn(10)-5)*time.Minute)

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ID:        uuid.NewString(),
		ExpiresAt: jwt.NewNumericDate(accessExpirationDate),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
	})

	ss, err := tok.SignedString(jwtSigningKey)
	if err != nil {
		return jwtKey{}, err
	}

	key := jwtKey{
		expirationDate: accessExpirationDate,
		token:          ss,
	}

	return key, nil
}

func (k *jwtKey) SaveCookie(c echo.Context) error {
	sess, err := session.Get("SESSIONID", c)
	if err != nil {
		return err
	}

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   int(k.expirationDate.Unix() - time.Now().Unix()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	sess.Values["token"] = k.token

	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return err
	}

	return nil
}
