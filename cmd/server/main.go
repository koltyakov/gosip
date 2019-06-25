package main

import (
	"log"
	"net/http"

	strategy "github.com/koltyakov/gosip/auth/saml"
	"github.com/koltyakov/gosip/cmd/server/handlers"
)

func main() {

	auth := &strategy.AuthCnfg{}
	err := auth.ReadConfig("./config/private.spo-user.json")
	if err != nil {
		log.Fatalf("unable to get config: %v", err)
	}

	http.HandleFunc("/digest", handlers.GetDigest(auth))
	http.HandleFunc("/web", handlers.GetWeb(auth))
	http.HandleFunc("/file", handlers.GetFile(auth))
	http.HandleFunc("/", handlers.Proxy(auth))

	log.Fatal(http.ListenAndServe(":8081", nil))

}
