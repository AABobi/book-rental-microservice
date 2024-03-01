package middlewares

import (
	"book-rental/utils"
	"net/http"
)

const servicePassword = "superSecretBooksPassword"

func Authorization(handler http.Handler) http.Handler {
	return handlerFunc(handler)
}

var handlerFunc = func(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("key")

		compareKey := utils.CheckPasswordHash(servicePassword, key)

		if !compareKey {
			http.Error(w, "Incorrect key", http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	})
}
