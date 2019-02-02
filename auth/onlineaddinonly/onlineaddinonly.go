package onlineaddinonly

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/koltyakov/gosip/cnfg"
)

// GetAuth : get auth
func GetAuth(creds *cnfg.AuthCnfgAddinOnly) (string, error) {
	parsedURL, err := url.Parse(creds.SiteURL)
	if err != nil {
		return "", err
	}

	authURL, err := getAuthURL(creds.Realm)
	if err != nil {
		return "", err
	}

	servicePrincipal := "00000003-0000-0ff1-ce00-000000000000" // TODO: mode to constance
	resource := fmt.Sprintf("%s/%s@%s", servicePrincipal, parsedURL.Hostname(), creds.Realm)
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

	fmt.Println(authURL)

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
		ExpiresOn   int32  `json:"expires_on,string"`
		Resource    string `json:"resource"`
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int32  `json:"expires_in,string"`
		NotBefore   int32  `json:"not_before,string"`
	}

	results := &getAuthResponse{}

	err = json.Unmarshal(data, &results)
	if err != nil {
		return "", err
	}

	// TODO: cache

	// fmt.Printf("Token: %s\n", results.AccessToken)

	return results.AccessToken, nil

}

func getAuthURL(realm string) (string, error) {
	endpoint := getAcsRealmURL(realm)
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
			return endpoint.Location, nil
		}
	}

	return "", errors.New("No OAuth2 protocol location found")
}

func getAcsRealmURL(realm string) string {
	return fmt.Sprintf("https://%s/metadata/json/1?realm=%s", "accounts.accesscontrol.windows.net", realm) // TODO: Add endpoint mapping
}
