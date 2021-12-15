package http

import (
	"io/fs"
	"net/http"
)

func UIHandler(fs fs.FS) http.Handler {
	return http.FileServer(http.FS(fs))
}
