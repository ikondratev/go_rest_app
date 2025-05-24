package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"go_rest_app/main/clients/telegram"
	"go_rest_app/main/lib/constants"
	"go_rest_app/main/lib/settings"
)

type App struct {
	Port     string
	Secret   string
	TgClient *telegram.Client
	Log      *slog.Logger
}

type GlitchtipAlert struct {
	Attachments []struct {
		Title     string `json:"title"`
		TitleLink string `json:"title_link"`
	} `json:"attachments"`
}

type RequestBody struct {
	Message string `json:"message"`
}

const (
	forbidden  = "forbidden"
	badRequest = "bad request"
	okResponse = "OK"
)

func main() {
	settings := settings.MustLoad()
	log := setupLog(settings.Env)

	app := App{
		Port:     settings.Port,
		Secret:   settings.Secret,
		TgClient: telegram.New(settings.Telegram),
		Log:      log,
	}

	if err := app.Start(); err != nil {
		os.Exit(1)
	}
}

func (a *App) Start() error {
	router := mux.NewRouter()

	router.Handle("/send", a.withLogging(a.authenticated(a.handlePost))).Methods(http.MethodPost)
	router.HandleFunc("/send/{token}/alert", a.handleAlert).Methods(http.MethodPost)

	addr := fmt.Sprintf(":%s", a.Port)
	a.Log.Info(fmt.Sprintf("Starting server on %s", addr))
	return http.ListenAndServe(addr, router)
}

// POST handler (with header auth)
func (a *App) handlePost(w http.ResponseWriter, r *http.Request) {
	var body RequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		a.respondError(w, badRequest, http.StatusBadRequest)
		return
	}

	if err := a.TgClient.SendMessage(body.Message); err != nil {
		a.respondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	a.respondOK(w)
}

func (a *App) handleAlert(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]

	if token != a.Secret {
		a.respondError(w, forbidden, http.StatusForbidden)
		return
	}

	var alert GlitchtipAlert
	if err := json.NewDecoder(r.Body).Decode(&alert); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		a.Log.Error("Invalid JSON body", slog.String("error", err.Error()))
		return
	}

	message := "🚨 Exception occurred!"
	if len(alert.Attachments) > 0 {
		att := alert.Attachments[0]
		message += fmt.Sprintf("\n👾 %s\n🔗 %s", att.Title, att.TitleLink)
	}

	if err := a.TgClient.SendMessage(message); err != nil {
		a.respondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	a.respondOK(w)
}

func (a *App) withLogging(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.Log.Info("Received request", slog.String("method", r.Method), slog.String("path", r.URL.Path))
		h(w, r)
	}
}

func (a *App) authenticated(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authentication") != a.Secret {
			a.respondError(w, forbidden, http.StatusForbidden)
			return
		}
		h(w, r)
	}
}

func (a *App) respondOK(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(okResponse))
}

func (a *App) respondError(w http.ResponseWriter, msg string, status int) {
	a.Log.Error(msg)
	http.Error(w, msg, status)
}

func setupLog(env string) *slog.Logger {
	switch env {
	case constants.ProdEnv:
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}
}
