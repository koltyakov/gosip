// SharePoint REST API POST requests API Digest helper
package gosip

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	cache "github.com/patrickmn/go-cache"
)

var (
	storage = cache.New(5*time.Minute, 10*time.Minute)
)

type contextInfoResponse struct {
	D struct {
		GetContextWebInformation struct {
			FormDigestTimeoutSeconds time.Duration `json:"FormDigestTimeoutSeconds"`
			FormDigestValue          string        `json:"FormDigestValue"`
			LibraryVersion           string        `json:"LibraryVersion"`
		} `json:"GetContextWebInformation"`
	} `json:"d"`
}

// GetDigest retrieves and caches SharePoint API X-RequestDigest value
func GetDigest(client *SPClient) (string, error) {
	siteURL := client.AuthCnfg.GetSiteURL()

	cacheKey := siteURL + "@digest@" + fmt.Sprintf("%#v", client.AuthCnfg)
	if digestValue, found := storage.Get(cacheKey); found {
		return digestValue.(string), nil
	}

	contextInfoURL := siteURL + "/_api/ContextInfo"
	req, err := http.NewRequest("POST", contextInfoURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/json;odata=verbose")

	resp, err := client.Execute(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	results := &contextInfoResponse{}

	err = json.Unmarshal(data, &results)
	if err != nil {
		return "", err
	}

	if results.D.GetContextWebInformation.FormDigestValue == "" {
		return "", errors.New("received empty FormDigestValue")
	}

	expirity := (results.D.GetContextWebInformation.FormDigestTimeoutSeconds - 60) * time.Second

	storage.Set(
		cacheKey,
		results.D.GetContextWebInformation.FormDigestValue,
		expirity,
	)

	return results.D.GetContextWebInformation.FormDigestValue, nil
}
