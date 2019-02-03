package scenarios

import (
	"fmt"

	"github.com/koltyakov/gosip/auth/spoaddinonly"
	"github.com/koltyakov/gosip/cnfg"
)

// ConfigReaderTest : test scenario
func ConfigReaderTest() {
	config, _ := cnfg.InitAuthConfig("./config/private.saml.json", "")
	fmt.Println(config)
}

// ConfigReaderSpoAddinOnlyTest : test scenario
func ConfigReaderSpoAddinOnlyTest() {
	config := &spoaddinonly.AuthCnfg{}
	err := config.ReadConfig("./config/private.addinonly.json")
	if err != nil {
		fmt.Printf("Error reading config: %v", err)
		return
	}
	fmt.Println(config)
}
