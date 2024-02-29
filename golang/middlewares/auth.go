package middlewares

import (
	"book-rental/utils"
	"fmt"
	"net/http"
)

const servicePassword = "superSecretBooksPassword"

func Authorization(handler http.Handler) http.Handler {
	return handlerFunc(handler)
}

var handlerFunc = func(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("HANDLE GOLANG")
		key := r.Header.Get("key")

		compareKey := utils.CheckPasswordHash(servicePassword, key)

		if !compareKey {
			http.Error(w, "Incorrect key", 440)
			return
		}
		h.ServeHTTP(w, r)
	})
}

/*func Authorization(next http.Handler) http.Handler {
	return handler(next)
}

var handler = func(n http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		userId, err := utils.VerifyToken(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.Background()
		ctx = context.WithValue(ctx, "myKey", userId)

		n.ServeHTTP(w, r.WithContext(ctx))
	})
}*/
