package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/koltyakov/gosip"
	strategy "github.com/koltyakov/gosip/auth/ntlm"
)

func main() {
	configPath := "./config/private.ntlm.json"
	auth := &strategy.AuthCnfg{}

	err := auth.ReadConfig(configPath)
	if err != nil {
		fmt.Printf("Unable to get config: %v\n", err)
		return
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

	req.Header.Set("Accept", "application/json;odata=verbose")

	fmt.Printf("Requesting api endpoint: %s\n", endpoint)
	resp, err := client.Execute(req)
	if err != nil {
		fmt.Printf("Unable to request api: %v\n", err)
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Unable to read a response: %v\n", err)
		return
	}

	fmt.Printf("Raw data: %s", string(data))
}
