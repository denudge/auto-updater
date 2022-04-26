package main

import (
	"github.com/denudge/auto-updater/config"
	"log"
	"net/http"
)

type Api struct {
	app *App
}

func NewApi(app *App) *Api {
	return &Api{
		app: app,
	}
}

func (api *Api) homePage(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Hello, world"))
}

func (api *Api) Serve() {
	http.Handle("/", http.HandlerFunc(api.homePage))

	port := config.Get("API_PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Serving HTTP API on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
