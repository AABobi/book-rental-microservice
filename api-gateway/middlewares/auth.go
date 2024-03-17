package middlewares

import (
	"api-gateway/data"
	"api-gateway/helpers"
	"api-gateway/utils"
	"context"
	"encoding/json"
	"net/http"
)

const authServicePassword = "superSecretAuthDbPassword"

func Authorization(handler http.Handler) http.Handler {
	return handlerFunc(handler)
}

var handlerFunc = func(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user data.User
		var req *http.Request
		var response data.ResponseMessage

		url := "http://localhost:91/authorization"

		req, err := http.NewRequest("GET", url, nil)

		if err != nil {
			http.Error(w, "Get request problem", 440)
			return
		}

		hashedPassword, err := utils.HashPassword(authServicePassword)
		if err != nil {
			helpers.ErrorJson(w, "Password hashing problem", 500)
			return
		}
		token := r.Header.Get("Authorization")

		req.Header.Set("key", hashedPassword)
		req.Header.Set("Authorization", token)
		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			http.Error(w, "Send auth request problem", 440)
			return
		}

		err = json.NewDecoder(resp.Body).Decode(&user)

		if err != nil {
			err = json.NewDecoder(resp.Body).Decode(&response)

			if err != nil {
				http.Error(w, "New request problem", 440)
			}
			return
		}

		if user.UserID == 0 {
			helpers.ErrorJson(w, "Incorrect token", 500)
			return
		}
		ctx := context.Background()
		ctx = context.WithValue(ctx, "myKey", user.UserID)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
