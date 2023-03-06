package main

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/LassiHeikkila/flmnchll/content/contentdb"
)

func VideoDownloadHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]

	v, err := contentdb.GetVideo(id)
	if errors.Is(err, contentdb.ErrVideoNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	SendVideo(w, v)
}
