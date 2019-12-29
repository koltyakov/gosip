package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/api"
	strategy "github.com/koltyakov/gosip/auth/saml"
)

func main() {
	// Getting auth params and client
	client, err := getAuthClient()
	if err != nil {
		log.Fatalln(err)
	}

	// Binding SharePoint API
	sp := api.NewSP(client)

	// Custom headers
	headers := map[string]string{
		"Accept": "application/json;odata=minimalmetadata",
	}
	config := &api.RequestConfig{Headers: headers}

	// Chainable request sample
	data, err := sp.Conf(config).Web().Webs().Get()
	if err != nil {
		log.Fatalln(err)
	}

	res := &struct {
		Value []struct {
			ID    string `json:"Id"`
			Title string `json:"Title"`
		} `json:"value"`
	}{}

	if err := json.Unmarshal(data, &res); err != nil {
		log.Fatalf("unable to parse the response: %v", err)
	}

	for _, list := range res.Value {
		fmt.Printf("%+v\n", list)
	}

}

func getAuthClient() (*gosip.SPClient, error) {
	configPath := "./config/private.spo-user.json"
	auth := &strategy.AuthCnfg{}
	if err := auth.ReadConfig(configPath); err != nil {
		return nil, fmt.Errorf("unable to get config: %v", err)
	}
	return &gosip.SPClient{AuthCnfg: auth}, nil
}
