package cnfg

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// AuthCnfgAddinOnly - auth config structure
type AuthCnfgAddinOnly struct {
	SiteURL      string `json:"siteUrl"`
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	Realm        string `json:"realm"`
}

// InitAuthConfigAddinOnly constructs auth config
func InitAuthConfigAddinOnly(privateFile, masterKey string) (*AuthCnfgAddinOnly, error) {
	jsonFile, err := os.Open(privateFile)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	config := &AuthCnfgAddinOnly{}
	json.Unmarshal(byteValue, &config)

	return config, nil
}
