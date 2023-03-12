package main

import (
	_ "embed"
	"net/http"
)

//go:embed index.html
var indexHtml []byte

func ServeHTML(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	_, _ = w.Write(indexHtml)
}
