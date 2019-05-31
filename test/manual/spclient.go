package manual

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"runtime"
	"time"

	ntlmssp "github.com/Azure/go-ntlmssp"
	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/auth/adfs"
	"github.com/koltyakov/gosip/auth/fba"
	"github.com/koltyakov/gosip/auth/ntlm"
)

// SPClientTest : api call test
func SPClientTest() {
	configs := [][]string{
		{"../../config/private.adfs.json", "adfs"},
		{"../../config/private.fba.json", "fba"},
		{"../../config/private.ntlm.json", "ntlm"},
	}
	for _, c := range configs {
		authInitTest(c[0], c[1])
	}
}

func authInitTest(cnfgPath, strategy string) {
	_, filename, _, _ := runtime.Caller(1)
	configPath := path.Join(path.Dir(filename), cnfgPath)

	var auth gosip.AuthCnfg

	switch strategy {
	case "fba":
		auth = &fba.AuthCnfg{}
	case "ntlm":
		auth = &ntlm.AuthCnfg{}
	case "adfs":
		auth = &adfs.AuthCnfg{}
	default:
		fmt.Println("Can't resolve auth strategy")
		return
	}

	err := auth.ReadConfig(configPath)
	if err != nil {
		fmt.Printf("Unable to get config: %v\n", err)
		return
	}

	client := &gosip.SPClient{
		AuthCnfg: auth,
	}

	if strategy == "ntlm" {
		client.Transport = ntlmssp.Negotiator{
			RoundTripper: &http.Transport{},
		}
	}

	apiCallTest(client, auth.GetSiteURL())
}

func apiCallTest(client *gosip.SPClient, siteURL string) {
	fmt.Println("")
	startAt := time.Now()
	endpoint := siteURL + "/_api/web?$select=Title"
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

	type apiResponse struct {
		Result struct {
			Title string `json:"Title"`
		} `json:"d"`
	}

	results := &apiResponse{}

	err = json.Unmarshal(data, &results)
	if err != nil {
		fmt.Printf("Unable to parse a response: %v\n", err)
		return
	}

	// fmt.Println("=== Response from API ===")
	fmt.Printf("SiteURL: %s, Web title: %v\n", siteURL, results.Result.Title)
	fmt.Printf("API requested in, sec: %f\n\n", time.Since(startAt).Seconds())
}
