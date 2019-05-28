package main

import (
	"log"
	"net/http"

	"github.com/koltyakov/gosip/auth/ntlm"
	"github.com/koltyakov/gosip/cmd/server/handlers"
)

func main() {

	auth := &ntlm.AuthCnfg{}
	err := auth.ReadConfig("./config/private.onprem-ntlm.json")
	if err != nil {
		log.Fatalf("unable to get config: %v", err)
	}

	http.HandleFunc("/web", handlers.GetWeb(auth))
	http.HandleFunc("/file", handlers.GetFile(auth))
	http.HandleFunc("/", handlers.Proxy(auth))

	log.Fatal(http.ListenAndServe(":8081", nil))

}
