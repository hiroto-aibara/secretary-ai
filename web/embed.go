package web

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed all:dist
var distFS embed.FS

func SPAHandler() http.HandlerFunc {
	distContent, _ := fs.Sub(distFS, "dist")
	fileServer := http.FileServer(http.FS(distContent))

	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/" {
			path = "index.html"
		} else {
			path = path[1:] // strip leading /
		}

		if _, err := distContent.Open(path); err != nil {
			// SPA fallback: serve index.html for non-file routes
			r.URL.Path = "/"
		}
		fileServer.ServeHTTP(w, r)
	}
}
