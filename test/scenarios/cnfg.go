package scenarios

import (
	"fmt"

	"github.com/koltyakov/gosip/cnfg"
)

// ConfigReaderTest : test scenario
func ConfigReaderTest() {
	config, _ := cnfg.InitAuthConfig("./config/private.saml.json", "")
	fmt.Println(config)
}

// ConfigReader2Test : test scenario
func ConfigReader2Test() {
	config, _ := cnfg.InitAuthConfigAddinOnly("./config/private.addinonly.json", "")
	fmt.Println(config)
}
