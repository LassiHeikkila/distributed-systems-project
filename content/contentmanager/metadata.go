package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/LassiHeikkila/flmnchll/content/contentdb"
	"github.com/gorilla/mux"
)

const (
	contentTypeJson = "application/json"
)

func VideoMetadataHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]

	v, err := contentdb.GetVideo(id)
	if errors.Is(err, contentdb.ErrVideoNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add(contentTypeHeaderKey, contentTypeJson)
	_, _ = w.Write(b)
	// TODO: can we do something about possible errors?
	// possible actions:
	// - logging
	// - ???
}

func VideoFileMetadataHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]

	v, err := contentdb.GetVideoFile(id)
	if errors.Is(err, contentdb.ErrVideoNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add(contentTypeHeaderKey, contentTypeJson)
	_, _ = w.Write(b)
}
