package adfs

import (
	"fmt"
	"testing"
	"time"
)

func TestGettingAuthToken(t *testing.T) {
	auth := &AuthCnfg{}
	err := auth.ReadConfig("../../config/private.adfs.json")

	if err != nil {
		t.Error(err)
	}
	if auth.SiteURL == "" {
		t.Error("Got empty config property")
	}
	if auth.Username == "" || auth.Password == "" {
		t.Error("Doesn't contain username and/or password properties")
	}

	cookie, err := auth.GetAuth()
	if err != nil {
		t.Error(err)
	}
	if cookie == "" {
		t.Error("AuthCookie is blank")
	}

	// Second auth should involve caching and be instant
	startAt := time.Now()
	cookie, err = auth.GetAuth()
	if err != nil {
		t.Error(err)
	}
	if time.Since(startAt).Seconds() > 0.0001 {
		t.Error(fmt.Sprintf("Possible caching issue, too slow read: %f", time.Since(startAt).Seconds()))
	}

}
