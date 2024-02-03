package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"go_rest_app/main/clients/telegram"
	"go_rest_app/main/lib/settings"
)

type App struct {
	Port     string
	Secret   string
	TgClient *telegram.Client
}

type RequestBody struct {
	Message   string `json:"message"`
}

const (
	tgBotHost = "api.telegram.org"
	auth_token = "Authentication"
	channelID = 564138790
)

var (
	ErrUnknownEvent = errors.New("unknown event")
)

func (a *App) Start() {
	http.Handle("/send", a.logreq(ok))
	addr := fmt.Sprintf(":%s", a.Port)
	log.Printf("Starting app on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func (a *App) logreq(f func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost: 
				log.Printf("path %s", r.URL.Path)

				if r.Header.Get(auth_token) != a.Secret {
					http.Error(w, "forbidden.", http.StatusForbidden)
					return
				}

				if err := a.handlePost(w, r); err != nil {
					http.Error(w, "bad request.", http.StatusBadRequest)
					return
				}

				f(w, r)
			default:
				http.Error(w, "404 not found.", http.StatusNotFound)
				return
			}
	})
}

func (a *App) handlePost(w http.ResponseWriter, r *http.Request) error {
	var requestBody RequestBody
	
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return err
	}

	messageFromBody := requestBody.Message
	
	if err := a.TgClient.SendMessage(channelID, messageFromBody); err != nil {
		return err
	}

	return nil
}

func ok(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func main() {
	settings := settings.New(".env")
	server := App{
			Port: settings.Port,
			TgClient: telegram.New(tgBotHost, settings.TgToken),
			Secret: settings.Secret,
	}
	server.Start()
}
