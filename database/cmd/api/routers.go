package main

import (
	"auth-db/authenticate"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"

	"github.com/go-chi/cors"
)

func routes() http.Handler {
	mux := chi.NewRouter()

	// specify who is allowed to connect
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:81"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	//mux.Use(middleware.Heartbeat("/ping"))

	mux.Post("/authenticate", FindUser)
	mux.With(authenticate.Authenticate).Post("/add-new-user", CreateUser)
	mux.Get("/get-all-users", GetAllUsers)
	return mux
}
func allowOnlyLocalhostMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request's Origin header matches "http://localhost:80"
		origin := r.Header.Get("Origin")
		fmt.Println(origin)
		fmt.Println("asd")
		/*if origin == "http://localhost:80" {
			// Allow the request if it's from the specified origin
			w.Header().Set("Access-Control-Allow-Origin", origin)
			next.ServeHTTP(w, r)
			return
		}

		// Reject the request if it's not from the specified origin
		http.Error(w, "Forbidden: Origin not allowed", http.StatusForbidden)*/
	})
}
