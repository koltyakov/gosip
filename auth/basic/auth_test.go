package basic

import (
	"testing"

	h "github.com/koltyakov/gosip/test/helpers"
)

func TestGettingAuthToken(t *testing.T) {
	err := h.CheckRequest(
		&AuthCnfg{},
		"./config/private.basic.json",
		[]string{"SiteURL", "Username", "Password"},
	)
	if err != nil {
		t.Error(err)
	}
}
