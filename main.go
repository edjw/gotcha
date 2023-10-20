package main

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/edjw/gotcha/components"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/", templ.Handler(components.Home()).ServeHTTP)

	http.ListenAndServe(":3333", r)
}
