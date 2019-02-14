package scenarios

import (
	"fmt"

	"github.com/koltyakov/gosip/auth/addin"
)

// ConfigReaderSpoAddinOnlyTest : test scenario
func ConfigReaderSpoAddinOnlyTest() {
	config := &addin.AuthCnfg{}
	err := config.ReadConfig(resolveCnfgPath("./config/private.addin.json"))
	if err != nil {
		fmt.Printf("Error reading config: %v", err)
		return
	}
	fmt.Println(config)
}
