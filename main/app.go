package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// Add configurable structure
type App struct {
	Port string
}

// Start server with args
func (a *App) Start() {
	http.Handle("/ping", logreq(ping))
	addr := fmt.Sprintf(":%s", a.Port)
	log.Printf("Starting app on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// Load env
func env(key, defaultValue string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return val
}

// Add log for requrest as middleware
func logreq(f func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case "GET":
				http.Error(w, "404 not found.", http.StatusNotFound)
				return
			case "POST":
				log.Printf("pathL %s", r.URL.Path)
				f(w, r)
			default:
				fmt.Fprintf(w, "Sorry, only POST methods are supported.")	
			}
	})
}

// Ping page
func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Pong\n")
}

//  Main point
func main() {
	server := App{
			Port: env("PORT", "7878"),
	}
	server.Start()
}
