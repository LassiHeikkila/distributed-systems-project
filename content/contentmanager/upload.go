package main

import (
	"io"
	"net/http"
	"os"
	"path"
)

const (
	kB = 1024
	MB = 1024 * kB
)

func VideoUploadHandler(w http.ResponseWriter, req *http.Request) {
	err := req.ParseMultipartForm(64 * MB)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for k, fileHeaders := range req.MultipartForm.File {
		_ = k
		for i, h := range fileHeaders {
			_ = i
			f, err := h.Open()
			if err != nil {
				/* ??? */
			}
			wf, err := os.CreateTemp("/tmp", "")
			if err != nil {
				/* ??? */
			}
			_, _ = io.Copy(wf, f)

			os.Rename(wf.Name(), path.Join(videoFileDirectory, h.Filename))
		}
	}
}
