package ntlm

import (
	"testing"

	h "github.com/koltyakov/gosip/test/helpers"
)

func TestGettingAuthToken(t *testing.T) {
	err := h.CheckRequest(
		&AuthCnfg{},
		"./config/private.ntlm.json",
		[]string{"SiteURL", "Username", "Password"},
	)
	if err != nil {
		t.Error(err)
	}
}
