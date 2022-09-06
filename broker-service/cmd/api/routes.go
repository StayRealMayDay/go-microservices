package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *Application) routers() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Heartbeat("/ping"))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})
	router.Post("/borker", app.Borker)

	router.Post("/handle", app.HandleSubmission)

	return router
}
