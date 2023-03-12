package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/LassiHeikkila/flmnchll/account/accountdb"
	"github.com/LassiHeikkila/flmnchll/helpers/httputils"
)

const (
	serviceSecretEnvKey = "account_service_internal_api_secret"
)

func main() {
	var (
		dbPath             string
		allowedCORSOrigins string
		httpPort           uint
		serviceSecret      string = os.Getenv(serviceSecretEnvKey)
	)

	if serviceSecret == "" {
		fmt.Println("internal API secret not defined, generating...")
		serviceSecret = accountdb.GenerateUUID()
		fmt.Println("internal API secret is:", serviceSecret)
	}

	flag.StringVar(&dbPath, "db", "account.db", "Path to database file")
	flag.StringVar(&allowedCORSOrigins, "cors", "*", "Comma-separated list of accepted CORS origins")
	flag.UintVar(&httpPort, "httpPort", 8080, "HTTP port")

	flag.Parse()

	if err := accountdb.Connect(dbPath); err != nil {
		fmt.Println("error connecting to database:", err)
		return
	}

	if err := accountdb.Init(); err != nil {
		fmt.Println("error initializing database:", err)
		return
	}

	r := mux.NewRouter()

	// CORS handling courtesy of:
	// https://stackoverflow.com/a/40987389/13580269
	headersOK := handlers.AllowedHeaders([]string{
		"Authorization",
		"Content-Type",
		"Access-Control-Request-Headers",
		"Access-Control-Request-Method",
	})
	originsOK := handlers.AllowedOrigins(httputils.ParseAllowedOrigins(allowedCORSOrigins))
	methodsOK := handlers.AllowedMethods([]string{http.MethodPost, http.MethodGet, http.MethodPut, http.MethodDelete})

	// signup cannot require authentication (unless its some static API key known by the frontend)
	r.HandleFunc("/account/signup", SignupHandler).
		Methods(http.MethodPost)

	r.HandleFunc("/account/login", LoginHandler).
		Methods(http.MethodPost)

	// look up account details based on token, will be used by other flmnchll services
	// not sure if this is the best way to do it, but it ought to work for now
	r.HandleFunc("/account/lookup/{token}", AccountLookupHandler).
		Methods(http.MethodGet)

	// the rest need to be authenticated
	r.Handle("/account/info/{id}", httputils.NewAuthMiddleware(AccountInfoHandler)).
		Methods(http.MethodGet)
	r.Handle("/account/info/{id}", httputils.NewAuthMiddleware(AccountInfoUpdateHandler)).
		Methods(http.MethodPut)
	r.Handle("/account/info/{id}", httputils.NewAuthMiddleware(AccountInfoDeleteHandler)).
		Methods(http.MethodDelete)
	r.Handle("/auth/validate", httputils.NewAuthMiddleware(TokenAuthenticateHandler)).
		Methods(http.MethodGet)
	r.Handle("/auth/invalidate", httputils.NewAuthMiddleware(TokenDeauthenticateHandler)).
		Methods(http.MethodPost)

	r.Handle("/internal/token/validate/{token}", httputils.NewServiceAuthMiddleware(serviceSecret, ServiceHandlerValidateToken)).
		Methods(http.MethodGet)
	r.Handle("/internal/account/info/{id}", httputils.NewServiceAuthMiddleware(serviceSecret, ServiceHandlerAccountLookup)).
		Methods(http.MethodGet)

	s := &http.Server{
		Handler: handlers.CombinedLoggingHandler(
			log.Writer(),
			handlers.CORS(
				originsOK,
				headersOK,
				methodsOK,
			)(r)),
		Addr:         fmt.Sprintf(":%d", httpPort),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			log.Println("HTTP server returned error:", err)
		}
	}()

	<-ctx.Done()

	log.Println("shutting down HTTP server")
	// give server 15s to shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()
	if err := s.Shutdown(shutdownCtx); err != nil {
		log.Println("error shutting down HTTP server:", err)
	}
}
