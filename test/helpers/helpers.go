package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"

	"github.com/koltyakov/gosip"
	u "github.com/koltyakov/gosip/test/utils"
)

// CheckAuth : common test case
func CheckAuth(auth gosip.AuthCnfg, cnfgPath string, required []string) error {
	err := auth.ReadConfig(u.ResolveCnfgPath(cnfgPath))
	if err != nil {
		return err
	}

	for _, prop := range required {
		v := getPropVal(auth, prop)
		if v == "" {
			return fmt.Errorf("doesn't contain required property value: %s", prop)
		}
	}

	token, err := auth.GetAuth()
	if err != nil {
		return err
	}
	if token == "" {
		return fmt.Errorf("accessToken is blank")
	}

	// Second auth should involve caching and be instant
	startAt := time.Now()
	token, err = auth.GetAuth()
	if err != nil {
		return err
	}
	if time.Since(startAt).Seconds() > 0.0001 {
		return fmt.Errorf("possible caching issue, too slow read: %f", time.Since(startAt).Seconds())
	}

	return nil
}

// CheckRequest : try sending basic request
func CheckRequest(auth gosip.AuthCnfg, cnfgPath string, required []string) error {
	err := auth.ReadConfig(u.ResolveCnfgPath(cnfgPath))
	if err != nil {
		return err
	}

	for _, prop := range required {
		v := getPropVal(auth, prop)
		if v == "" {
			return fmt.Errorf("doesn't contain required property value: %s", prop)
		}
	}

	client := &gosip.SPClient{
		AuthCnfg: auth,
	}

	apiEndpoint := auth.GetSiteURL() + "/_api/web?$select=Title"
	req, err := http.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return fmt.Errorf("unable to create a request: %v", err)
	}

	req.Header.Set("Accept", "application/json;odata=verbose")

	resp, err := client.Execute(req)
	if err != nil {
		return fmt.Errorf("unable to request api: %v", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unable to read a response: %v", err)
	}

	type apiResponse struct {
		D struct {
			Title string `json:"Title"`
		} `json:"d"`
		Error struct {
			Message struct {
				Value string `json:"value"`
			} `json:"message"`
		} `json:"error"`
	}
	results := &apiResponse{}

	err = json.Unmarshal(data, &results)
	if err != nil {
		return fmt.Errorf("unable to parse a response: %v", err)
	}

	if results.Error.Message.Value != "" {
		return fmt.Errorf(results.Error.Message.Value)
	}

	return nil
}

func getPropVal(v gosip.AuthCnfg, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	if !f.IsValid() {
		return ""
	}
	return string(f.String())
}
