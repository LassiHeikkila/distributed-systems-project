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

	"github.com/LassiHeikkila/flmnchll/content/contentdb"
	"github.com/LassiHeikkila/flmnchll/helpers/httputils"
)

// globals
var (
	videoFileDirectory     = ""
	videoFileTempDirectory = ""
	redisQueue             string
)

func main() {
	var (
		dbPath             string
		allowedCORSOrigins string
		httpPort           uint
		redisAddr          string
		redisPassword      string
	)

	flag.StringVar(&dbPath, "db", "content.db", "Path to database file")
	flag.StringVar(&videoFileDirectory, "contentDir", ".", "Directory where video files are stored")
	flag.StringVar(&videoFileTempDirectory, "tempDir", ".", "Directory where video files are temporarily stored while being downloaded")
	flag.StringVar(&allowedCORSOrigins, "cors", "*", "Comma-separated list of accepted CORS origins")
	flag.UintVar(&httpPort, "httpPort", 8080, "HTTP port")

	flag.Parse()

	redisAddr = os.Getenv("REDIS_ADDR")
	redisPassword = os.Getenv("REDIS_PASSWORD")
	redisQueue = os.Getenv("REDIS_JOB_QUEUE")

	if err := contentdb.Connect(dbPath); err != nil {
		log.Fatal("error connecting to database:", err)
	}
	defer contentdb.Disconnect()

	if err := contentdb.Init(); err != nil {
		log.Fatal("error initializing database:", err)
	}

	log.Println("connecting to redis...")
	if err := ConnectToRedis(context.Background(), redisAddr, redisPassword, 0); err != nil {
		log.Fatal("error connecting to redis: ", err)
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

	r.HandleFunc("/video/upload", VideoUploadHandler).
		Methods(http.MethodPost)
	r.HandleFunc("/video/download/{id}", VideoDownloadHandler).
		Methods(http.MethodGet)
	r.HandleFunc("/video/download/file/{id}", VideoFileDownloadHandler).
		Methods(http.MethodGet)
	r.HandleFunc("/video/info/{id}", VideoMetadataHandler).
		Methods(http.MethodGet)
	r.HandleFunc("/video/info/file/{id}", VideoFileMetadataHandler).
		Methods(http.MethodGet)
	r.HandleFunc("/video/subtitles/{id}/{language}", VideoSubtitleHandler).
		Methods(http.MethodGet)
	r.HandleFunc("/video/search", VideoSearchHandler).
		Methods(http.MethodGet)

	r.HandleFunc("/", ServeHTML).Methods(http.MethodGet)

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

	log.Printf("Open http://localhost:%d in your browser to add videos\n", httpPort)

	<-ctx.Done()

	log.Println("shutting down HTTP server")
	// give server 15s to shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()
	if err := s.Shutdown(shutdownCtx); err != nil {
		log.Println("error shutting down HTTP server:", err)
	}
}
