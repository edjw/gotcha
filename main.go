package main

import (
	"github.com/a-h/templ"
	"github.com/edjw/gotcha/components"
	"github.com/edjw/gotcha/friendlyServer"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", templ.Handler(components.Home()).ServeHTTP)

	// r.Get("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("hello world"))
	// })

	friendlyServer.FriendlyServer(r)

}

// func FriendlyServer() {
// 	port, portExists := os.LookupEnv("PORT")
// 	devEnv, devEnvExists := os.LookupEnv("GO_ENV")

// 	var ip string = "0.0.0.0" // Standard IP address on Fly and Render

// 	// For development, this relies on setting the GO_ENV variable
// 	// to "development" in your .zshrc or .bashrc file with
// 	// export GO_ENV=development
// 	// Then run source ~/.zshrc or source ~/.bashrc
// 	// Then close the terminal and reopen it.

// 	if devEnvExists && devEnv == "development" {
// 		ip = "127.0.0.1"
// 	}

// 	if !portExists {
// 		port = "8080"
// 	}

// 	server := &http.Server{
// 		Addr:         ip + ":" + port,
// 		Handler:      nil,
// 		ReadTimeout:  5 * time.Second,
// 		WriteTimeout: 10 * time.Second,
// 		IdleTimeout:  120 * time.Second,
// 	}

// 	log.Printf("\n\nListening on:\nhttp://%s:%v\n\n", ip, port)

// 	// TODO: Except it'll actually be http://localhost:3000 if we're using browser-sync. Make this not print?

// 	if err := server.ListenAndServe(); err != nil {
// 		log.Fatal(err)
// 	}
// }
