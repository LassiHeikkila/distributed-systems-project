package httputils

import (
	"net/http"
	"strings"

	"github.com/LassiHeikkila/flmnchll/account/accountdb"
)

const (
	authHeaderKey = "Authorization"
	bearerPrefix  = "Bearer "

	unauthorizedError = `{"ok":false,"msg":"you must provide authentication bearer token with Authorization header"}`
)

type AuthMw struct {
	next func(http.ResponseWriter, *http.Request)
}

func NewAuthMiddleware(next func(http.ResponseWriter, *http.Request)) *AuthMw {
	return &AuthMw{
		next: next,
	}
}

func (a *AuthMw) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	_, err := accountdb.AuthenticateToken(GetAuthToken(req))
	if err != nil {
		// log something?
		// since AuthenticateToken returned an error, the token is not valid
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(unauthorizedError))
		return
	}

	a.next(w, req)
}

func GetAuthToken(req *http.Request) string {
	return strings.TrimPrefix(req.Header.Get(authHeaderKey), bearerPrefix)
}
