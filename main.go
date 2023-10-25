package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/a-h/templ"
	"github.com/edjw/gotcha/components"
	"github.com/edjw/gotcha/components/partials"
	"github.com/edjw/gotcha/friendlyServer"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/unrolled/secure"
)

//go:embed public/*
var embeddedFiles embed.FS

func partialsRouter() *chi.Mux {
	deploymentSiteURL, deploymentSiteURLExists := os.LookupEnv("DEPLOYMENT_SITE_URL")
	devEnv, devEnvExists := os.LookupEnv("GO_ENV")

	onlyInternal := func(next http.Handler) http.Handler {
		// This middleware checks that the request is coming from the same server.
		// It's not foolproof, but it's a good start
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check the Referer header
			referer := r.Header.Get("Referer")

			var siteURL string

			if devEnvExists && devEnv == "development" {
				siteURL = "http://127.0.0.1"
			} else if deploymentSiteURLExists {
				siteURL = deploymentSiteURL
			}

			if !strings.HasPrefix(referer, siteURL) {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	r := chi.NewRouter()

	if deploymentSiteURLExists || (devEnvExists && devEnv == "development") {
		r.Use(onlyInternal)
	}

	// A map of partial routes to templ component partials.
	partialsMap := map[string]func() templ.Component{
		"new_headline": partials.NewHeadline,
	}

	r.Get("/{partialName}", func(w http.ResponseWriter, r *http.Request) {
		partialName := chi.URLParam(r, "partialName")

		partialComponent, ok := partialsMap[partialName]
		if !ok {
			http.Error(w, "Partial not found.", http.StatusNotFound)
			return
		}
		templ.Handler(partialComponent()).ServeHTTP(w, r)
	})

	return r
}

func main() {
	devEnv, devEnvExists := os.LookupEnv("GO_ENV")

	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	secureMiddleware := secure.New(secure.Options{
		ReferrerPolicy:     "same-origin",
		ContentTypeNosniff: true,
		FrameDeny:          true,
		BrowserXssFilter:   true,
		IsDevelopment:      devEnvExists && devEnv == "development",
	})
	r.Use(secureMiddleware.Handler)

	// Serve the public folder.
	var fileServer http.Handler

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

	// Partials / Fragments
	r.Mount("/partials", partialsRouter())

	// Start the server
	friendlyServer.FriendlyServer(r)
}
