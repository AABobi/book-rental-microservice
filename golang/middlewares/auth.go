package middlewares

import (
	"book-rental/utils"
	"context"

	"net/http"
)

func Authorization(next http.Handler) http.Handler {
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
}
