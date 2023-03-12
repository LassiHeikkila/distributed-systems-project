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

	"github.com/LassiHeikkila/flmnchll/account/accountservice/accountclient"
	"github.com/LassiHeikkila/flmnchll/helpers/httputils"
	"github.com/LassiHeikkila/flmnchll/room/roomdb"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	var (
		peerServerAddr     string
		accountServiceAddr string
		dbPath             string
		allowedCORSOrigins string
		httpPort           uint
	)

	flag.StringVar(&peerServerAddr, "peerjs", "", "Address of PeerJS server")
	flag.StringVar(&accountServiceAddr, "accountServiceAddr", "", "Address of account-service")
	flag.StringVar(&dbPath, "db", "room.db", "Path to database file")
	flag.StringVar(&allowedCORSOrigins, "cors", "*", "Comma-separated list of accepted CORS origins")
	flag.UintVar(&httpPort, "httpPort", 8080, "HTTP port")

	if peerServerAddr != "" {
		RegisterPeerServer(peerServerAddr)
	}

	if accountServiceAddr != "" {
		accountclient.SetAccountServiceAddr(accountServiceAddr)
	}

	if err := roomdb.Connect(dbPath); err != nil {
		log.Fatal("error connecting to database:", err)
	}
	defer roomdb.Disconnect()

	if err := roomdb.Init(); err != nil {
		log.Fatal("error initializing database:", err)
	}

	// hardcode admin user
	adm, _ := roomdb.CreateUser("4dm1n", "admin")
	// hardcode room ab12 watching some video
	room, _ := roomdb.CreateRoom(*adm, peerServerAddr)
	room.ShortID = "ab12"
	_ = roomdb.UpdateRoom(*room)

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

	r.HandleFunc("/room/join/{id}/{username}", JoinRoomHandler).
		Methods(http.MethodPost)
	r.HandleFunc("/room/leave/{username}", LeaveRoomHandler).
		Methods(http.MethodPost)

	r.HandleFunc("/room/details/{id}", GetRoomDetailsHandler).
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

	log.Println("serving HTTP...")
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
