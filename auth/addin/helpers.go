package addin

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	cache "github.com/patrickmn/go-cache"
)

var (
	storage = cache.New(5*time.Minute, 10*time.Minute)
)

// GetAuth : get auth
func GetAuth(creds *AuthCnfg) (string, error) {
	parsedURL, err := url.Parse(creds.SiteURL)
	if err != nil {
		return "", err
	}

	cacheKey := parsedURL.Host + "@addinonly@" + creds.ClientID + "@" + creds.ClientSecret
	if accessToken, found := storage.Get(cacheKey); found {
		return accessToken.(string), nil
	}

	realm, err := getRealm(creds)
	if err != nil {
		return "", err
	}
	creds.Realm = realm

	authURL, err := getAuthURL(creds.Realm)
	if err != nil {
		return "", err
	}

	servicePrincipal := "00000003-0000-0ff1-ce00-000000000000" // TODO: move to constants
	resource := fmt.Sprintf("%s/%s@%s", servicePrincipal, parsedURL.Host, creds.Realm)
	fullClientID := fmt.Sprintf("%s@%s", creds.ClientID, creds.Realm)

	// type getAuthForm struct {
	// 	GrantType    string `json:"grant_type"`
	// 	ClientID     string `json:"client_id"`
	// 	ClientSecret string `json:"client_secret"`
	// 	Resource     string `json:"resource"`
	// }

	// type getAuthRequest struct {
	// 	JSON bool        `json:"json"`
	// 	Form getAuthForm `json:"form"`
	// }

	// reqBody := &getAuthRequest{
	// 	JSON: true,
	// 	Form: getAuthForm{
	// 		GrantType:    "client_credentials",
	// 		ClientID:     fullClientID,
	// 		ClientSecret: creds.ClientSecret,
	// 		Resource:     resource,
	// 	},
	// }

	params := url.Values{}
	params.Set("grant_type", "client_credentials")
	params.Set("client_id", fullClientID)
	params.Set("client_secret", creds.ClientSecret)
	params.Set("resource", resource)

	// reqBodyJSON, err := json.Marshal(reqBody)
	// if err != nil {
	// 	return "", err
	// }

	// proxyURL, _ := url.Parse("http://127.0.0.1:8888")
	// http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyURL), TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}

	resp, err := http.Post(authURL, "application/x-www-form-urlencoded", strings.NewReader(params.Encode())) // bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	type getAuthResponse struct {
		// ExpiresOn   int32  `json:"expires_on,string"`
		// NotBefore   int32  `json:"not_before,string"`
		// Resource    string `json:"resource"`
		AccessToken string        `json:"access_token"`
		TokenType   string        `json:"token_type"`
		ExpiresIn   time.Duration `json:"expires_in,string"`
		Error       string        `json:"error_description"`
	}

	results := &getAuthResponse{}

	err = json.Unmarshal(data, &results)
	if err != nil {
		return "", err
	}

	if results.Error != "" {
		return "", fmt.Errorf("%s", results.Error)
	}

	storage.Set(cacheKey, results.AccessToken, (results.ExpiresIn-60)*time.Second)

	return results.AccessToken, nil

}

func getAuthURL(realm string) (string, error) {
	endpoint := fmt.Sprintf("https://%s/metadata/json/1?realm=%s", "accounts.accesscontrol.windows.net", realm) // TODO: Add endpoint mapping

	cacheKey := endpoint
	if authURL, found := storage.Get(cacheKey); found {
		return authURL.(string), nil
	}

	resp, err := http.Get(endpoint)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	type getAuthURLResponse struct {
		Endpoints []struct {
			Protocol string `json:"protocol"`
			Location string `json:"location"`
		} `json:"endpoints"`
	}

	results := &getAuthURLResponse{}

	err = json.Unmarshal(data, &results)
	if err != nil {
		return "", err
	}

	for _, endpoint := range results.Endpoints {
		if endpoint.Protocol == "OAuth2" {
			storage.Set(cacheKey, endpoint.Location, 60*time.Minute)
			return endpoint.Location, nil
		}
	}

	return "", errors.New("no OAuth2 protocol location found")
}

func getRealm(creds *AuthCnfg) (string, error) {
	if creds.Realm != "" {
		return creds.Realm, nil
	}

	parsedURL, err := url.Parse(creds.SiteURL)
	if err != nil {
		return "", err
	}

	cacheKey := parsedURL.Host + "@realm"
	if realm, found := storage.Get(cacheKey); found {
		return realm.(string), nil
	}

	endpoint := creds.SiteURL + "/_vti_bin/client.svc"
	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer ")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	authHeader := resp.Header.Get("www-authenticate")

	for _, part := range strings.Split(authHeader, `",`) {
		p := strings.Split(part, `="`)
		if p[0] == "Bearer realm" {
			storage.Set(cacheKey, p[0], 60*time.Minute)
			return p[1], nil
		}
	}

	return "", errors.New("wasn't able to get Realm")
}
