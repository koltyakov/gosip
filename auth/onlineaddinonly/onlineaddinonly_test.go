package onlineaddinonly

import (
	"fmt"
	"testing"
	"time"

	"github.com/koltyakov/gosip/cnfg"
)

func TestGettingAuthToken(t *testing.T) {
	config, err := cnfg.InitAuthConfigAddinOnly("../../config/private.addinonly.json", "")
	if err != nil {
		t.Error(err)
	}
	if config.SiteURL == "" {
		t.Error("Got empty config property")
	}
	if config.ClientID == "" || config.ClientSecret == "" {
		t.Error("Doesn't contain clientId and/or secretId properties")
	}

	token, err := GetAuth(config)
	if err != nil {
		t.Error(err)
	}
	if token == "" {
		t.Error("AccessToken is blank")
	}

	// Second auth should involve caching and be instant
	startAt := time.Now()
	token, err = GetAuth(config)
	if err != nil {
		t.Error(err)
	}
	if time.Since(startAt).Seconds() > 0.0001 {
		t.Error(fmt.Sprintf("Possible caching issue, too slow read: %f", time.Since(startAt).Seconds()))
	}

}
