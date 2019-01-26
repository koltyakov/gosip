package cnfg

import (
	"fmt"
	"testing"
)

func TestConfigReader(t *testing.T) {
	config, err := InitAuthConfig("private.json", "DUMMY_KEY")
	if err != nil {
		t.Error(err)
	}
	if config.SiteURL == "" {
		t.Error("Got empty config property")
	}
	fmt.Println(config.SiteURL, config.Password)
	if config.Password != "secret" {
		t.Error("Got wrong password property")
	}
}
