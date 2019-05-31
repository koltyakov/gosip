package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/koltyakov/gosip"
	strategy "github.com/koltyakov/gosip/auth/ntlm"
)

func main() {
	configPath := "./config/private.onprem-ntlm.json"
	auth := &strategy.AuthCnfg{}

	err := auth.ReadConfig(configPath)
	if err != nil {
		fmt.Printf("unable to get config: %v\n", err)
		return
	}

	client := &gosip.SPClient{
		AuthCnfg: auth,
	}

	// Assumes you have Custom list created
	endpoint := client.AuthCnfg.GetSiteURL() + "/_api/web/lists/getByTitle('Custom')/items"
	req, err := http.NewRequest(
		"POST",
		endpoint,
		bytes.NewBuffer([]byte(`{"__metadata":{"type":"SP.Data.CustomListItem"},"Title":"Test"}`)),
	)
	if err != nil {
		fmt.Printf("unable to create a request: %v", err)
		return
	}

	req.Header.Set("Accept", "application/json;odata=verbose")
	req.Header.Set("Content-Type", "application/json;odata=verbose")

	resp, err := client.Execute(req)
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

	fmt.Printf("Results: %s", data)
}
