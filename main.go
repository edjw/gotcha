package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/a-h/templ"
	"github.com/edjw/gotcha/friendlyServer"
	"github.com/edjw/gotcha/pages"
	"github.com/edjw/gotcha/pages/partials"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/unrolled/secure"
)

//go:embed public/*
var embeddedFiles embed.FS

func pagesRouter(pagesMap map[string]func() templ.Component) *chi.Mux {

	r := chi.NewRouter()

	pathHandler := func(path string, w http.ResponseWriter, r *http.Request) {
		page, ok := pagesMap[path]
		if !ok {
			http.Error(w, "Page not found.", http.StatusNotFound)
			return
		}
		templ.Handler(page()).ServeHTTP(w, r)
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		pathHandler("/", w, r)

	})

	r.Get("/{pageName}", func(w http.ResponseWriter, r *http.Request) {

		pageName := chi.URLParam(r, "pageName")
		pathHandler("/"+pageName, w, r)
	})

	return r
}

func partialsRouter(partialsMap map[string]func() templ.Component) *chi.Mux {
	deploymentSiteURL, deploymentSiteURLExists := os.LookupEnv("DEPLOYMENT_SITE_URL")
	devEnv, devEnvExists := os.LookupEnv("GO_ENV")

	onlyInternal := func(next http.Handler) http.Handler {
		// This middleware checks that the request is coming from the same URL.
		// It's not foolproof, but it's a good start at keep partials internal.
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check the Referer header
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

	r := chi.NewRouter()

	if deploymentSiteURLExists || (devEnvExists && devEnv == "development") {
		r.Use(onlyInternal)
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

	// A map of page routes to pages written as templ components.
	pagesMap := map[string]func() templ.Component{
		"/":      pages.Home,
		"/about": pages.About,
	}

	// A map of partial routes to partials written as templ components.
	partialsMap := map[string]func() templ.Component{
		"new_headline": partials.NewHeadline,
	}

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

	// Page routes
	// r.Get("/", templ.Handler(pages.Index()).ServeHTTP)
	r.Mount("/", pagesRouter(pagesMap))

	// Partials / Fragments
	r.Mount("/partials", partialsRouter(partialsMap))

	// Start the server using the local Friendly Server package
	friendlyServer.FriendlyServer(r)
}
