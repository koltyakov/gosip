package scenarios

import (
	"fmt"

	"github.com/koltyakov/gosip/auth/addin"
	"github.com/koltyakov/gosip/cnfg"
)

// ConfigReaderTest : test scenario
func ConfigReaderTest() {
	config, _ := cnfg.InitAuthConfig("./config/private.saml.json", "")
	fmt.Println(config)
}

// ConfigReaderSpoAddinOnlyTest : test scenario
func ConfigReaderSpoAddinOnlyTest() {
	config := &addin.AuthCnfg{}
	err := config.ReadConfig("./config/private.addin.json")
	if err != nil {
		fmt.Printf("Error reading config: %v", err)
		return
	}
	fmt.Println(config)
}
