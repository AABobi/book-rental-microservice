package main

import (
	"auth-db/authenticate"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"

	"github.com/go-chi/cors"
)

func routes() http.Handler {
	mux := chi.NewRouter()

	// specify who is allowed to connect
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))

	mux.Get("/test", TestAuth)
	mux.With(authenticate.Authenticate).Post("/authenticate", FindUser)
	mux.With(authenticate.Authenticate).Get("/authorization", Auth)
	mux.With(authenticate.Authenticate).Post("/add-new-user", CreateUser)
	mux.With(authenticate.Authenticate).Get("/get-all-users", GetAllUsers)
	return mux
}
