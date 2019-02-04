package scenarios

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"runtime"
	"time"

	"github.com/koltyakov/gosip/auth/fba"
)

// SPClient : SharePoint HTTP client struct
type SPClient struct {
	http.Client
	Auth       *fba.AuthCnfg
	ConfigPath string
}

// Execute : SharePoint HTTP client Do method
func (c *SPClient) Execute(req *http.Request) (*http.Response, error) {
	fmt.Println("Injecting auth call...")
	if c.ConfigPath != "" && c.Auth.SiteURL == "" {
		c.Auth.ReadConfig(c.ConfigPath)
	}
	authCookie, err := c.Auth.GetAuth()
	if err != nil {
		fmt.Printf("unable to get auth: %v", err)
		res := &http.Response{
			Status:     "401 Access Denied",
			StatusCode: 401,
			Request:    req,
		}
		return res, err
	}
	req.Header.Set("Cookie", authCookie)
	return c.Do(req)
}

// GetFbaAuthTest : get auth test scenario
func GetFbaAuthTest() {
	startAtProc := time.Now()
	startAt := time.Now()

	_, filename, _, _ := runtime.Caller(1)
	configPath := path.Join(path.Dir(filename), "../../config/private.fba.json")

	client := &SPClient{
		Auth:       &fba.AuthCnfg{},
		ConfigPath: configPath,
	} // &http.Client{}

	auth := &fba.AuthCnfg{}
	err := auth.ReadConfig(configPath)
	if err != nil {
		fmt.Printf("unable to get config: %v", err)
		return
	}
	// // fmt.Printf("config: %v\n", auth)
	// fmt.Printf("siteUrl: %s\n", auth.SiteURL)
	// fmt.Printf("config read in, sec: %f\n", time.Since(startAt).Seconds())
	// startAt = time.Now()

	// authCookie, err := auth.GetAuth()
	// if err != nil {
	// 	fmt.Printf("unable to get auth: %v", err)
	// 	return
	// }
	// fmt.Printf("auth cookie: %s\n", authCookie)
	// fmt.Printf("authenticated in, sec: %f\n", time.Since(startAt).Seconds())
	// startAt = time.Now()

	///
	// authCookie, err = auth.GetAuth()
	// if err != nil {
	// 	fmt.Printf("unable to get auth: %v", err)
	// 	return
	// }
	// // fmt.Printf("auth token: %s\n", authToken)
	// fmt.Printf("second auth (cached) in, sec: %f\n", time.Since(startAt).Seconds())
	// startAt = time.Now()
	///

	apiEndpoint := auth.SiteURL + "/_api/web?$select=Title"
	req, err := http.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		fmt.Printf("unable to create a request: %v", err)
		return
	}

	req.Header.Set("Accept", "application/json;odata=minimalmetadata")
	// req.Header.Set("Cookie", authCookie) // should be injected via SPClient

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
	fmt.Printf("api requested in, sec: %f\n", time.Since(startAt).Seconds())
	fmt.Printf("summary time, sec: %f\n", time.Since(startAtProc).Seconds())

}
