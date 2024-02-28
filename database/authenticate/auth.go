package authenticate

import (
	"auth-db/utils"
	"net/http"
)

const servicePassword = "superSecretAuthDbPassword"

func Authenticate(handler http.Handler) http.Handler {
	return handlerFunc(handler)
}

var handlerFunc = func(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("key")

		compareKey := utils.CheckPasswordHash(servicePassword, key)

		if !compareKey {
			http.Error(w, "Incorrect key", 440)
			return
		}
		h.ServeHTTP(w, r)
	})
}
