package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/a-h/templ"
	"github.com/edjw/gotcha/functions/fetchRandomUser"
	"github.com/edjw/gotcha/functions/friendlyServer"
	"github.com/edjw/gotcha/html/pages"
	"github.com/edjw/gotcha/html/partials"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/unrolled/secure"
)

//go:embed public/*
var embeddedFiles embed.FS

func main() {

	r := chi.NewRouter()
	partialsRouter := chi.NewRouter()

	// Apply middleware to the partials subrouter
	// If you set an environment variable called DEPLOYMENT_SITE_URL as the url of your app, then you can go someway towards making the partials routes only accessible by your site and not from others.

	protectPartials := true

	deploymentSiteURL, deploymentSiteURLExists := os.LookupEnv("DEPLOYMENT_SITE_URL")

	devEnv, devEnvExists := os.LookupEnv("GO_ENV")

	onlyInternal := func(next http.Handler) http.Handler {
		// This middleware checks that the request is coming from the same URL.
		// It's not foolproof, but it seems an ok start at keep partials internal and avoiding hotlinking
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			referer := r.Header.Get("Referer")

			var siteURL string

			if devEnvExists && devEnv == "development" {
				siteURL = "http://127.0.0.1:8080"
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
	if protectPartials && (deploymentSiteURLExists || (devEnvExists && devEnv == "development")) {
		partialsRouter.Use(onlyInternal)
	}

	// General Middleware
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

	// Page routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(pages.Home()).ServeHTTP(w, r)
	})

	r.Get("/about", func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(pages.About()).ServeHTTP(w, r)
	})

	// Partials routes
	partialHandlers := map[string]func() (templ.Component, error){

		"new_headline": func() (templ.Component, error) {
			return partials.NewHeadline(), nil
		},

		"random_name": func() (templ.Component, error) {
			userData, err := fetchRandomUser.FetchRandomUser()
			if err != nil {
				return nil, err
			}
			return partials.RandomName(*userData), nil
		},

		// Add other partial handlers here...
	}

	partialsRouter.Get("/{partialName}", func(w http.ResponseWriter, r *http.Request) {
		partialName := chi.URLParam(r, "partialName")
		handler, ok := partialHandlers[partialName]
		if !ok {
			http.Error(w, "Partial not found.", http.StatusNotFound)
			return
		}
		component, err := handler()
		if err != nil {
			http.Error(w, "Failed to handle request.", http.StatusInternalServerError)
			return
		}
		templ.Handler(component).ServeHTTP(w, r)
	})

	r.Mount("/partials", partialsRouter)

	// Start the server using the local Friendly Server package
	friendlyServer.FriendlyServer(r)
}
