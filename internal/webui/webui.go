package webui

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed static
var staticFiles embed.FS

// FS returns an http.FileSystem serving the embedded static files.
func FS() http.FileSystem {
	sub, err := fs.Sub(staticFiles, "static")
	if err != nil {
		panic(err)
	}
	return http.FS(sub)
}
