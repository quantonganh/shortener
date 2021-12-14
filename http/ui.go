package http

import (
	"io/fs"
	"log"
	"net/http"

	"github.com/quantonganh/shortener/ui"
)

func UIHandler(dir string) http.Handler {
	publicFS, err := fs.Sub(ui.Public, dir)
	if err != nil {
		log.Fatal(err)
	}

	return http.FileServer(http.FS(publicFS))
}
