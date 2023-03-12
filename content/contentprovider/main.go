package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/LassiHeikkila/flmnchll/content/contentdb"
	"github.com/LassiHeikkila/flmnchll/helpers/httputils"
)

// globals
var (
	videoFileDirectory = ""
)

func main() {
	var (
		dbPath             string
		allowedCORSOrigins string
		httpPort           uint
	)

	flag.StringVar(&dbPath, "db", "content.db", "Path to database file")
	flag.StringVar(&videoFileDirectory, "contentDir", ".", "Directory where video files are stored")
	flag.StringVar(&allowedCORSOrigins, "cors", "", "Comma-separated list of accepted CORS origins")
	flag.UintVar(&httpPort, "httpPort", 8080, "HTTP port")

	flag.Parse()

	if err := contentdb.Connect(dbPath); err != nil {
		fmt.Println("error connecting to database:", err)
		return
	}
	defer contentdb.Disconnect()

	if err := contentdb.Init(); err != nil {
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

	r.HandleFunc("/video/download/{id}", VideoDownloadHandler).
		Methods(http.MethodGet)
	r.HandleFunc("/video/info/{id}", VideoMetadataHandler).
		Methods(http.MethodGet)
	r.HandleFunc("/video/subtitles/{id}/{language}", VideoSubtitleHandler).
		Methods(http.MethodGet)
	r.HandleFunc("/video/search", VideoSearchHandler).
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

func parseAllowedOrigins(s string) []string {
	return strings.Split(s, ",")
}
