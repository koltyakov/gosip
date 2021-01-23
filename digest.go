package gosip

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
)

var storage = cache.New(5*time.Minute, 10*time.Minute)

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
func GetDigest(context context.Context, client *SPClient) (string, error) {
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

	if context != nil {
		req = req.WithContext(context)
	}

	req.Header.Set("Accept", "application/json;odata=verbose")

	resp, err := client.Execute(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
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

	expiry := (results.D.GetContextWebInformation.FormDigestTimeoutSeconds - 60) * time.Second

	storage.Set(
		cacheKey,
		results.D.GetContextWebInformation.FormDigestValue,
		expiry,
	)

	return results.D.GetContextWebInformation.FormDigestValue, nil
}
