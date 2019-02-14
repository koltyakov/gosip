package fba

import (
	"testing"

	h "github.com/koltyakov/gosip/test/helpers"
)

func TestGettingAuthToken(t *testing.T) {
	err := h.CheckAuth(
		&AuthCnfg{},
		"./config/private.fba.json",
		[]string{"SiteURL", "Username", "Password"},
	)
	if err != nil {
		t.Error(err)
	}
}
