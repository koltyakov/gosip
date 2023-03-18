package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/koltyakov/gosip"
	u "github.com/koltyakov/gosip/test/utils"
)

// CheckRequest : try sending basic request
func CheckRequest(auth gosip.AuthCnfg, cnfgPath string) error {
	err := auth.ReadConfig(u.ResolveCnfgPath(cnfgPath))
	if err != nil {
		return err
	}

	client := &gosip.SPClient{
		AuthCnfg: auth,
	}

	endpoint := auth.GetSiteURL() + "/_api/web?$select=Title"
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return fmt.Errorf("unable to create a request: %w", err)
	}

	req.Header.Set("Accept", "application/json;odata=verbose")

	resp, err := client.Execute(req)
	if err != nil {
		return fmt.Errorf("unable to request api: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unable to read a response: %w", err)
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
		return fmt.Errorf("unable to parse a response: %w", err)
	}

	if results.Error.Message.Value != "" {
		return fmt.Errorf(results.Error.Message.Value)
	}

	return nil
}
