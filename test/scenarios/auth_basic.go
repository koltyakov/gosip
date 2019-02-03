package scenarios

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"runtime"
	"time"

	ntlmssp "github.com/Azure/go-ntlmssp"
	"github.com/koltyakov/gosip/auth/basic"
)

// GetBasicAuthTest : get auth test scenario
func GetBasicAuthTest() {
	startAtProc := time.Now()
	startAt := time.Now()

	_, filename, _, _ := runtime.Caller(1)
	configPath := path.Join(path.Dir(filename), "../../config/private.basic.json")

	auth := &basic.AuthCnfg{}
	err := auth.ReadConfig(configPath)
	if err != nil {
		fmt.Printf("unable to get config: %v", err)
		return
	}
	fmt.Printf("config: %v\n", auth)
	fmt.Printf("config read in, sec: %f\n", time.Since(startAt).Seconds())
	startAt = time.Now()

	client := &http.Client{
		Transport: ntlmssp.Negotiator{
			RoundTripper: &http.Transport{},
		},
	}

	apiEndpoint := auth.SiteURL + "/_api/web?$select=Title"
	req, err := http.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		fmt.Printf("unable to create a request: %v", err)
		return
	}

	req.Header.Set("Accept", "application/json;odata=verbose") // Just because it's 2013 in my case
	req.SetBasicAuth(auth.Username, auth.Password)

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

	// fmt.Println(resp.StatusCode)
	// fmt.Println(string(data))

	type apiResponse struct {
		Result struct {
			Title string `json:"Title"`
		} `json:"d"`
	}

	results := &apiResponse{}

	err = json.Unmarshal(data, &results)
	if err != nil {
		fmt.Printf("unable to parse a response: %v", err)
		return
	}

	fmt.Println("=== Response from API ===")
	fmt.Printf("Web title: %v\n", results.Result.Title)
	fmt.Printf("api requested in, sec: %f\n", time.Since(startAt).Seconds())
	fmt.Printf("summary time, sec: %f\n", time.Since(startAtProc).Seconds())

}
