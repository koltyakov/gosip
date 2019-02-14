package tmg

import (
	"testing"

	h "github.com/koltyakov/gosip/test/helpers"
)

func TestGettingAuthToken(t *testing.T) {
	err := h.CheckAuth(
		&AuthCnfg{},
		"./config/private.tmg.json",
		[]string{"SiteURL", "Username", "Password"},
	)
	if err != nil {
		t.Error(err)
	}
}
