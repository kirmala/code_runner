package middleware

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/kirmala/code_runner/http_server/internal/service"
	"github.com/labstack/echo/v5"
	slogctx "github.com/veqryn/slog-context"
)

const UserIdKey = "user_id"


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

		id, err := a.Authenticator.Authenticate(c.Request().Context(), tokenStr)
		if err != nil {
			return err
		}

		ctx := c.Request().Context()
		ctx = slogctx.Prepend(ctx,
			slog.Any(UserIdKey, id),
		)
		
		c.SetRequest(c.Request().WithContext(ctx))
		c.Set(UserIdKey, id)

		return next(c)
	}
}