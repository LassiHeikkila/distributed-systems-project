package main

import (
	"net/http"

	_ "github.com/LassiHeikkila/flmnchll/content/contentdb"
)

func VideoSubtitleHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
