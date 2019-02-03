package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/koltyakov/gosip/auth/addin"
)

func main() {
	auth := &addin.AuthCnfg{
		SiteURL:      "https://contoso.sharepoint.com/sites/my_site",
		ClientID:     "[CLIENT_ID]",
		ClientSecret: "[CLIENT_SECRET]",
	}

	authToken, err := auth.GetAuth()
	if err != nil {
		fmt.Printf("unable to authenticate: %v", err)
		return
	}

	apiEndpoint := auth.SiteURL + "/_api/web?$select=Title"
	req, err := http.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		fmt.Printf("unable to create a request: %v", err)
		return
	}

	req.Header.Set("Accept", "application/json;odata=minimalmetadata")
	req.Header.Set("Authorization", "Bearer "+authToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("unable to request api: %v", err)
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("unable to read a response: %v", err)
		return
	}

	type apiResponse struct {
		Title string `json:"Title"`
	}

	results := &apiResponse{}

	err = json.Unmarshal(data, &results)
	if err != nil {
		fmt.Printf("unable to parse a response: %v", err)
		return
	}

	fmt.Println("=== Response from API ===")
	fmt.Printf("Web title: %v\n", results.Title)
}
