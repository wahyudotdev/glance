package api

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

//go:embed all:dist
var staticAssets embed.FS

func (s *APIServer) registerStaticRoutes() {
	sub, err := fs.Sub(staticAssets, "dist")
	if err != nil {
		return
	}

	// Serve static files from embedded FS
	s.app.Use("/", filesystem.New(filesystem.Config{
		Root:         http.FS(sub),
		Index:        "index.html",
		NotFoundFile: "index.html", // SPA fallback
	}))
}
