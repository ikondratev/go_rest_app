package main

import (
	"encoding/json"
	"fmt"
	// "log"
	"log/slog"
	"net/http"
	"os"

	"go_rest_app/main/clients/telegram"
	"go_rest_app/main/lib/constants"
	"go_rest_app/main/lib/settings"
)

type App struct {
	Port      string
	Secret    string
	TgClient  *telegram.Client
	Log 	  *slog.Logger
}

type RequestBody struct {
	Message string `json:"message"`
}

const (
	authToken = "Authentication"
	forbidden = "forbidden."
	badRequest = "bad request."
	notFound = "not found."
)

func main() {
	// Load settings
	settings := settings.MustLoad()

	// Configuraste logger
	log := setupLog(settings.Env)

	server := App{
			Port: settings.Port,
			TgClient: telegram.New(settings.Telegram),
			Secret: settings.Secret,
			Log: log,
	}

	if err := server.Start(); err != nil {
		os.Exit(1)
	}
}

func (a *App) Start() error {
	http.Handle("/send", a.logreq(ok))

	addr := fmt.Sprintf(":%s", a.Port)
	a.Log.Info(fmt.Sprintf("Starting app on port %s", addr))

	if err := http.ListenAndServe(addr, nil); err != nil {
		a.Log.Error(err.Error())
		return err
	}

	return nil
}


func (a *App) logreq(f func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get(authToken) != a.Secret {
				http.Error(w, forbidden, http.StatusForbidden)
				return
			}

			switch r.Method {
			case http.MethodPost: 
				a.Log.Info(fmt.Sprintf("path: %s", r.URL.Path))

				if err := a.handlePost(w, r); err != nil {
					http.Error(w, badRequest, http.StatusBadRequest)
					return
				}

				f(w, r)
			default:
				http.Error(w, notFound, http.StatusNotFound)
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
	
	if err := a.TgClient.SendMessage(messageFromBody); err != nil {
		return err
	}

	return nil
}

func ok(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func setupLog(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case constants.ProdEnv:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	return log
}
