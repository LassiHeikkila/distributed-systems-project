package main

import (
	"net/http"
	"strings"

	"github.com/LassiHeikkila/flmnchll/account/accountdb"
)

const (
	authHeaderKey = "Authorization"
	bearerPrefix  = "Bearer "
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
	v := req.Header.Get(authHeaderKey)
	v = strings.TrimPrefix(v, bearerPrefix)

	_, err := accountdb.AuthenticateToken(v)
	if err != nil {
		// log something?
		// since AuthenticateToken returned an error, the token is not valid
		return
	}

	a.next(w, req)
}
