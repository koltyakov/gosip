package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/koltyakov/gosip"
	strategy "github.com/koltyakov/gosip/auth/addin"
)

func main() {
	auth := &strategy.AuthCnfg{
		SiteURL:      os.Getenv("SPAUTH_SITEURL"),
		ClientID:     os.Getenv("SPAUTH_CLIENTID"),
		ClientSecret: os.Getenv("SPAUTH_CLIENTSECRET"),
	}

	client := &gosip.SPClient{
		AuthCnfg: auth,
	}

	endpoint := auth.GetSiteURL() + "/_api/web?$select=Title"
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		fmt.Printf("Unable to create a request: %v", err)
		return
	}

	req.Header.Set("Accept", "application/json;odata=minimalmetadata")

	resp, err := client.Execute(req)
	if err != nil {
		fmt.Printf("Unable to request api: %v", err)
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Unable to read a response: %v", err)
		return
	}

	type apiResponse struct {
		Title string `json:"Title"`
	}

	results := &apiResponse{}

	err = json.Unmarshal(data, &results)
	if err != nil {
		fmt.Printf("Unable to parse a response: %v", err)
		return
	}

	fmt.Printf("Web title: %v\n", results.Title)
}
