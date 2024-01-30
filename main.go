package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go_rest_app/main/clients/telegram"
	"log"
	"net/http"
	"os"
)

// Add configurable structure
type App struct {
	Port string
}

type RequestBody struct {
	Message string `json:"message"`
}

const (
	tgBotHost = "api.telegram.org"
)

func (a *App) Start(client *telegram.Client) {
	http.Handle("/ping", logreq(okPage, client))
	addr := fmt.Sprintf(":%s", a.Port)
	log.Printf("Starting app on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func env(key, defaultValue string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return val
}

func logreq(f func(w http.ResponseWriter, r *http.Request), client *telegram.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case "GET":
				http.Error(w, "404 not found.", http.StatusNotFound)
				return
			case "POST":
				log.Printf("path %s", r.URL.Path)

				if err := handlePost(w, r, client); err != nil {
					http.Error(w, "400 not found.", http.StatusBadRequest)
					return
				}

				f(w, r)
			default:
				fmt.Fprintf(w, "Sorry, only POST methods are supported.")	
			}
	})
}

func handlePost(w http.ResponseWriter, r *http.Request, client *telegram.Client) error {
	var requestBody RequestBody
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return err
	}

	messageFromBody := requestBody.Message

	if err := client.SendMessage(564138790, messageFromBody); err != nil {
		log.Fatal("error: ", err)
		return err
	}

	return nil
}

func okPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func mustToken() string {
	token := flag.String(
		"token",
		"default_token",
		"access to telegram bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}

func main() {
	server := App{
			Port: env("PORT", "7878"),
	}
	server.Start(telegram.New(tgBotHost, mustToken()))
}
