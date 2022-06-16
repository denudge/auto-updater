package api

import (
	"fmt"
	"github.com/denudge/auto-updater/cmd/catalog/app"
	"github.com/denudge/auto-updater/config"
	"log"
	"net/http"
)

type Api struct {
	mux     *http.ServeMux
	catalog *app.Catalog
}

func NewApi(catalog *app.Catalog) *Api {
	api := &Api{
		catalog: catalog,
		mux:     http.NewServeMux(),
	}

	api.setUpRoutes()

	return api
}

func (api *Api) setUpRoutes() {
	api.mux.Handle("/", http.HandlerFunc(api.homePage))
	api.mux.Handle("/register", http.HandlerFunc(api.register))
	// api.mux.Handle("/upgrade/check", http.HandlerFunc(api.findNextUpgrade))
	api.mux.Handle("/upgrade/step", http.HandlerFunc(api.findNextUpgrade))
}

func (api *Api) homePage(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Hello, world"))
}

func (api *Api) Serve() {
	port := config.Get("API_PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Serving HTTP API on port", port)
	log.Fatal(http.ListenAndServe(":"+port, api.mux))
}

// validateMethodIs checks if a given HTTP method is used. The error is written to the HTTP response
func (api *Api) validateMethodIs(w http.ResponseWriter, r *http.Request, method string) error {
	if r.Method != method {
		err := fmt.Errorf("method %s not allowed", r.Method)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	return nil
}
