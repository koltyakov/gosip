package cnfg

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/koltyakov/gosip/cpass"
)

// AuthCnfg - auth config structure
type AuthCnfg struct {
	SiteURL  string `json:"siteUrl"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// InitAuthConfig constructs auth config
func InitAuthConfig(privateFile, masterKey string) (config *AuthCnfg, err error) {
	jsonFile, err := os.Open(privateFile)
	if err != nil {
		return
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	config = &AuthCnfg{}
	json.Unmarshal(byteValue, &config)

	c := cpass.Cpass(masterKey)
	config.Password, _ = c.Decode(config.Password)

	return
}
