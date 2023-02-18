package main

import (
	"io"
	"net/http"
	"os"
	"path"

	"github.com/LassiHeikkila/flmnchll/content/contentdb"
)

const (
	contentTypeHeaderKey = "Content-Type"
	contentTypeMp4       = "video/mp4"
	contentTypeWebm      = "video/webm"
)

func SendVideo(w http.ResponseWriter, v *contentdb.Video) {
	// for now all files in one dir identified by FileID

	f, err := os.Open(path.Join(videoFileDirectory, v.FileID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch v.Encoding {
	case "mp4":
		w.Header().Add(contentTypeHeaderKey, contentTypeMp4)
	case "webm":
		w.Header().Add(contentTypeHeaderKey, contentTypeWebm)
	}

	_, _ = io.Copy(w, f)
	// TODO: can we do something about possible errors?
	// possible actions:
	// - logging
	// - ???
}
