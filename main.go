package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/edjw/gotcha/components"
	"github.com/edjw/gotcha/friendlyServer"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed public/*
var embeddedFiles embed.FS

func main() {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Serve the public folder.
	var fileServer http.Handler

	devEnv, devEnvExists := os.LookupEnv("GO_ENV")

	if devEnvExists && devEnv == "development" {
		fileServer = http.FileServer(http.Dir("./public"))
	} else {
		publicFS, err := fs.Sub(embeddedFiles, "public")
		if err != nil {
			log.Fatal(err)
		}

		fileServer = http.FileServer(http.FS(publicFS))
	}

	r.Handle("/public/*", http.StripPrefix("/public", fileServer))

	// Routes
	r.Get("/", templ.Handler(components.Index()).ServeHTTP)

	// Start the server
	friendlyServer.FriendlyServer(r)
}
