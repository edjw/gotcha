package friendlyServer

// This is a simple server to help development on a Mac. It also works in production on Fly and Render.
// It works around the firewall issue on Macs where you have to manually allow incoming connections to your Go app on every rebuild.

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
)

func FriendlyServer(router chi.Router) {
	port, portExists := os.LookupEnv("PORT")
	devEnv, devEnvExists := os.LookupEnv("GO_ENV")

	var ip string = "0.0.0.0" // Standard IP address on Fly and Render

	// For development, this relies on setting the GO_ENV variable
	// to "development" in your .zshrc or .bashrc file with
	// export GO_ENV=development
	// Then run source ~/.zshrc or source ~/.bashrc
	// Then close the terminal and reopen it.

	if devEnvExists && devEnv == "development" {
		ip = "127.0.0.1"
	}

	if !portExists {
		port = "8080"
	}

	server := &http.Server{
		Addr:         ip + ":" + port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("\n\nListening on:\nhttp://%s:%v\n\n", ip, port)

	// TODO: Except it'll actually be http://localhost:3000 if we're using browser-sync. Make this not print?

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
