package middlewares

import (
	"crypto/subtle"
	"strings"

	"github.com/go-raptor/raptor/v4"
	"github.com/go-raptor/raptor/v4/errs"
)

// UpdateTokenMiddleware guards the manual schedule refresh, which logs in to Edupage and fetches
// every class. Callers send "Authorization: Bearer <app.update_token>"; with no token configured
// the endpoint stays closed.
type UpdateTokenMiddleware struct {
	raptor.Middleware

	token []byte
}

func (m *UpdateTokenMiddleware) Setup() error {
	m.token = []byte(m.Config.AppConfig["update_token"])
	if len(m.token) == 0 {
		m.Log.Warn("update_token is not set, manual schedule refresh is disabled")
	}
	return nil
}

func (m *UpdateTokenMiddleware) Handle(ctx *raptor.Context, next func(*raptor.Context) error) error {
	if !bearerTokenMatches(ctx.Request().Header.Get("Authorization"), m.token) {
		return errs.NewErrorUnauthorized("Missing or invalid update token")
	}
	return next(ctx)
}

func bearerTokenMatches(header string, token []byte) bool {
	got, ok := strings.CutPrefix(header, "Bearer ")
	return ok && len(token) > 0 && subtle.ConstantTimeCompare([]byte(got), token) == 1
}
