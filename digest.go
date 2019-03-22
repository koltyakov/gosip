package gosip

import (
	"encoding/json"
	"errors"
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

	cacheKey := siteURL + "@digest@" + "" // TODO: add unique identity for the client
	if digestValue, found := storage.Get(cacheKey); found {
		return digestValue.(string), nil
	}

	contextInfoURL := siteURL + "/_api/contextinfo"
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

	storage.Set(
		cacheKey,
		results.D.GetContextWebInformation.FormDigestValue,
		results.D.GetContextWebInformation.FormDigestTimeoutSeconds*time.Second,
	)

	return results.D.GetContextWebInformation.FormDigestValue, nil
}
