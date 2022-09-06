package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *Application)routes() http.Handler{

	mux := chi.NewRouter()

	mux.Use(middleware.Heartbeat("/ping"))
	mux.Post("/authenticate", app.Authenticate)

	return mux
}