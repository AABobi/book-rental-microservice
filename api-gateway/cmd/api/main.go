package main

import (
	"fmt"
	"log"
	"net/http"
)

const webPort = "8082"

func main() {

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
