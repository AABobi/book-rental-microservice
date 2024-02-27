package main

import (
	"auth-db/db"
	"fmt"
	"log"
	"net/http"
)

const webPort = "91"

func main() {
	db.InitGDB()
	fmt.Println("test")
	log.Printf("Starting broker service on port %s\n", webPort)

	// define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: routes(),
	}

	// start the server
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
