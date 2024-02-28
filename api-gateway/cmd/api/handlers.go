package main

import (
	"api-gateway/data"
	"api-gateway/helpers"
	"api-gateway/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const authServicePassword = "superSecretAuthDbPassword"

func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LOGIN")
	var user data.User
	var response string
	err := helpers.ReadJSON(w, r, &user)

	if err != nil {
		http.Error(w, "Cannot read user json", 440)
		return
	}

	err = sendRequest(user, &response, w)

	if err != nil {
		http.Error(w, "New request problem", 440)
	}
	fmt.Println("Response Status:", response)
}

func SingUp(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SIGN")
	var user data.User
	var response string
	err := helpers.ReadJSON(w, r, &user)

	if err != nil {
		http.Error(w, "Cannot read user json", 440)
		return
	}

	err = sendRequest(user, &response, w)

	if err != nil {
		http.Error(w, "New request problem", 440)
	}
	fmt.Println("Response Status:", response)
}

func sendRequest(data any, response *string, w http.ResponseWriter) error {
	fmt.Println("SENDREQ")
	fmt.Println(data)
	jsonData, err := json.Marshal(data)
	req, err := http.NewRequest("POST", "http://localhost:91/add-new-user", bytes.NewBuffer(jsonData))

	if err != nil {
		http.Error(w, "New request problem", 440)
		return err
	}

	hashedPassword, err := utils.HashPassword(authServicePassword)
	req.Header.Set("key", hashedPassword)
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)

	if err != nil {
		return err
	}

	*response = buf.String()

	defer resp.Body.Close()
	return nil
}
