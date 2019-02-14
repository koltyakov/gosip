package addin

import (
	"fmt"
	"testing"
	"time"
)

func TestGettingAuthToken(t *testing.T) {
	err := h.CheckAuth(
		&AuthCnfg{},
		"./config/private.addin.json",
		[]string{"SiteURL", "ClientID", "ClientSecret"},
	)
	if err != nil {
		t.Error(err)
	}
}
