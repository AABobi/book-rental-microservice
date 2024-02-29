package middlewares

import (
	"api-gateway/data"
	"api-gateway/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type ResponseMessage struct {
	Message string `json:"message"`
}

const authServicePassword = "superSecretAuthDbPassword"

func Authorization(handler http.Handler) http.Handler {
	return handlerFunc(handler)
}

var handlerFunc = func(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user data.User
		var req *http.Request
		var response ResponseMessage

		url := "http://localhost:91/authorization"

		req, err := http.NewRequest("GET", url, nil)

		if err != nil {
			http.Error(w, "Get request problem", 440)
			return
		}

		hashedPassword, err := utils.HashPassword(authServicePassword)
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
				fmt.Println(err)
				http.Error(w, "New request problem", 440)
			}
			return
		}

		fmt.Println("Response Status:", user)
		fmt.Println("1")
		ctx := context.Background()
		fmt.Println("1")
		ctx = context.WithValue(ctx, "myKey", user.UserID)
		fmt.Println("1")
		h.ServeHTTP(w, r.WithContext(ctx))
		fmt.Println("1")
	})
}
