package main

import (
	"api-gateway/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"net/http"
)

func routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Post("/login", Login)
	mux.Post("/signup", SingUp)
	mux.Get("/test", DockerTest)
	mux.With(middlewares.Authorization).Get("/get", GetAllBooks)
	mux.With(middlewares.Authorization).Get("/get-available-books", GetAvailableBooks)
	mux.With(middlewares.Authorization).Get("/get-rented-books", GetRentedBooks)
	mux.With(middlewares.Authorization).Post("/rent-book", RentBook)
	mux.With(middlewares.Authorization).Post("/return-books", ReturnBook)
	return mux
}
