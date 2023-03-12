package main

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/LassiHeikkila/flmnchll/content/contentdb"
)

const (
	QueryResolutionKey = "res"
	QueryEncodingKey   = "enc"
)

func VideoDownloadHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	res := req.URL.Query().Get(QueryResolutionKey)
	enc := req.URL.Query().Get(QueryEncodingKey)

	v, err := contentdb.GetVideo(id)
	if errors.Is(err, contentdb.ErrVideoNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	f := pickEncodingAndResolution(v.Files, enc, res)
	if f == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	SendVideo(w, f)
}

func VideoFileDownloadHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]

	vf, err := contentdb.GetVideoFile(id)
	if errors.Is(err, contentdb.ErrVideoNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	SendVideo(w, vf)
}

func pickEncodingAndResolution(files []contentdb.VideoFile, enc string, res string) *contentdb.VideoFile {
	if res == "" {
		for _, r := range []string{resolution1080p, resolution720p, resolution480p} {
			for _, f := range files {
				if f.Encoding == enc && f.Resolution == r {
					return &f
				}
			}
		}
	}
	for _, f := range files {
		if f.Encoding == enc && f.Resolution == res {
			return &f
		}
	}

	return nil
}
