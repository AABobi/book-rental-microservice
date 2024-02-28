package main

import (
	"api-gateway/data"
	"api-gateway/helpers"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LOGIN")
	var user data.User
	err := helpers.ReadJSON(w, r, &user)

	if err != nil {
		http.Error(w, "Cannot read user json", 440)
		return
	}
	a := data.User{
		Email:    "f11",
		Password: "f11",
	}
	jsonData, err := json.Marshal(a)
	req, err := http.NewRequest("POST", "http://localhost:91/authenticate", bytes.NewBuffer(jsonData))

	if err != nil {
		http.Error(w, "New request problem", 440)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	respBytes := buf.String()

	respString := string(respBytes)

	fmt.Println("Response Status:", respString)
	/*body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Print the response body
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(body))*/
}
