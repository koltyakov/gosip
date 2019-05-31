package manual

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/koltyakov/gosip"
)

// CheckBasicPost : try creating an item
func CheckBasicPost(client *gosip.SPClient) (string, error) {
	endpoint := client.AuthCnfg.GetSiteURL() + "/_api/web/lists/getByTitle('Custom')/items"
	req, err := http.NewRequest(
		"POST",
		endpoint,
		bytes.NewBuffer([]byte(`{"__metadata":{"type":"SP.Data.CustomListItem"},"Title":"Test"}`)),
	)
	if err != nil {
		return "", fmt.Errorf("unable to create a request: %v", err)
	}

	req.Header.Set("Accept", "application/json;odata=verbose")
	req.Header.Set("Content-Type", "application/json;odata=verbose")

	resp, err := client.Execute(req)
	if err != nil {
		return "", fmt.Errorf("unable to request api: %v", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read a response: %v", err)
	}

	return string(data), nil
}
