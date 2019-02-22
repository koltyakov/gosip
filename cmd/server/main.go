package main

import (
	"log"
	"net/http"

	"github.com/koltyakov/gosip/auth/basic"
	"github.com/koltyakov/gosip/cmd/server/handlers"
)

func main() {

	auth := &basic.AuthCnfg{}
	err := auth.ReadConfig("./config/private.basic.json")
	if err != nil {
		log.Fatalf("unable to get config: %v", err)
	}

	http.HandleFunc("/favicon.ico", handlers.GetFavicon(auth))
	http.HandleFunc("/web", handlers.GetWeb(auth))
	http.HandleFunc("/file", handlers.GetFile(auth))

	log.Fatal(http.ListenAndServe(":8081", nil))

}
