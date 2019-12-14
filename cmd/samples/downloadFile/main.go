package main

import (
	"log"
	"os"

	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/api"
	strategy "github.com/koltyakov/gosip/auth/saml"
)

func main() {
	configPath := "./config/integration/private.spo.json"
	auth := &strategy.AuthCnfg{}

	err := auth.ReadConfig(configPath)
	if err != nil {
		log.Fatalf("unable to get config: %v", err)
	}

	client := &gosip.SPClient{
		AuthCnfg: auth,
	}

	sp := api.NewSP(client)

	data, err := sp.Web().GetFile("Shared Documents/Sub Folder/My File.xlsx").Download()
	if err != nil {
		log.Fatalf("unable to download a file: %v", err)
	}

	file, err := os.Create("My File.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	file.Write(data)
	if err != nil {
		log.Fatalf("unable to write to file: %v", err)
	}

	log.Println("Done!")

}
