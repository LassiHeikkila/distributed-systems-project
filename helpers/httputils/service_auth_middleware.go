package httputils

import (
	"net/http"
	"strings"
)

const (
	serviceSecretKeyPrefix = "Secret "
)

type ServiceAuthMw struct {
	// TODO: improve this so it's not just a static secret
	// secret should be some very long random string
	// it should only be used for internal service<->service API calls
	secret string
	next   func(http.ResponseWriter, *http.Request)
}

func NewServiceAuthMiddleware(secret string, next func(http.ResponseWriter, *http.Request)) *ServiceAuthMw {
	return &ServiceAuthMw{
		secret: secret,
		next:   next,
	}
}

func (a *ServiceAuthMw) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if a.secret != GetServiceAuthSecret(req) {
		// log something?
		// since AuthenticateToken returned an error, the token is not valid
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(unauthorizedError))
		return
	}

	a.next(w, req)
}

func GetServiceAuthSecret(req *http.Request) string {
	return strings.TrimPrefix(req.Header.Get(authHeaderKey), serviceSecretKeyPrefix)
}
