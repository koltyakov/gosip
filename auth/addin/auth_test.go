package addin

import (
	"fmt"
	"testing"
	"time"
)

func TestGettingAuthToken(t *testing.T) {
	auth := &AuthCnfg{}
	err := auth.ReadConfig("../../config/private.addinonly.json")

	if err != nil {
		t.Error(err)
	}
	if auth.SiteURL == "" {
		t.Error("Got empty config property")
	}
	if auth.ClientID == "" || auth.ClientSecret == "" {
		t.Error("Doesn't contain clientId and/or secretId properties")
	}

	token, err := auth.GetAuth()
	if err != nil {
		t.Error(err)
	}
	if token == "" {
		t.Error("AccessToken is blank")
	}

	// Second auth should involve caching and be instant
	startAt := time.Now()
	token, err = auth.GetAuth()
	if err != nil {
		t.Error(err)
	}
	if time.Since(startAt).Seconds() > 0.0001 {
		t.Error(fmt.Sprintf("Possible caching issue, too slow read: %f", time.Since(startAt).Seconds()))
	}

}
