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

func (a *App) Start() {
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

//  Main point
func main() {
	server := App{
			Port: env("PORT", "7878"),
	}
	server.Start()
}
