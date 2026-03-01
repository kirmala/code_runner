package middleware

import (
	"code_processor/http_server/service"
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
)

const UserIdKey = "userId"


func getAuthToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", service.ErrUnauthenticated{Msg: "invalid Authorization token format need to include Bearer"}
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")

	return token, nil
}

type Auth struct {
	Authenticator service.Authenticator
}

func (a Auth) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		tokenStr, err := getAuthToken(c.Request()) 
		if err != nil {
			return err
		}

		id, err := a.Authenticator.Authenticate(tokenStr)
		if err != nil {
			return err
		}

		c.Set(UserIdKey, id)

		return next(c)
	}
}