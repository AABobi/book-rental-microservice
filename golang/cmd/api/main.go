package main

import (
	"book-rental/db"
	"fmt"
	"log"
	"net/http"
)

const webPort = "80"

func main() {
	db.InitGDB()
	log.Printf("Starting broker service on port %s\n", webPort)

	_ = StartServer()
}

func StartServer() error {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

	return srv.ListenAndServe()
}
