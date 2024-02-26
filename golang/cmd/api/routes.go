package main

import (
	"book-rental/handlers"
	"book-rental/middlewares"
	"net/http"

	"github.com/go-chi/chi/v5"
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

	//mux.Use(middleware.Heartbeat("/ping"))

	mux.Get("/get-users", handlers.GetUsers)
	mux.Get("/get-all-books", handlers.GetAllBooks)

	mux.With(middlewares.Authorization).Get("/get", handlers.GetUsers)

	mux.Post("/signup", handlers.CreateNewUser)
	mux.Post("/login", handlers.Login)
	mux.With(middlewares.Authorization).Get("/get-rented-books", handlers.GetRentedBooks)
	mux.With(middlewares.Authorization).Post("/rent", handlers.RentBook)
	mux.With(middlewares.Authorization).Get("/get-available-books", handlers.GetAvailableBooks)
	mux.With(middlewares.Authorization).Post("/return-books", handlers.ReturnTheBook)

	return mux
}
