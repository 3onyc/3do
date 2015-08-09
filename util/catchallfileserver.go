package util

import (
	"log"
	"net/http"
	"path"
	"strings"
)

type catchAllFileHandler struct {
	FS           http.FileSystem
	InnerHandler http.Handler
	serveInstead string
}

func CatchAllFileServer(root http.FileSystem, serveInstead string) http.Handler {
	return &catchAllFileHandler{
		FS:           root,
		InnerHandler: http.FileServer(root),
		serveInstead: serveInstead,
	}
}

func (f *catchAllFileHandler) ServeFallback(w http.ResponseWriter, r *http.Request) {
	file, err := f.FS.Open(f.serveInstead)
	if err != nil {
		log.Printf("ServeFallback Error: %s\n", err)
	}

	d, err := file.Stat()
	if err != nil {
		log.Printf("ServeFallback Error: %s\n", err)
	}

	http.ServeContent(w, r, f.serveInstead, d.ModTime(), file)
}

func (f *catchAllFileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/") {
		r.URL.Path = "/" + r.URL.Path
	}

	findPath := path.Clean(r.URL.Path)
	file, err := f.FS.Open(findPath)
	if err != nil {
		f.ServeFallback(w, r)
		return
	}

	if _, err := file.Stat(); err != nil {
		f.ServeFallback(w, r)
		return
	}

	f.InnerHandler.ServeHTTP(w, r)
}
