package auth

import (
	"testing"
)

func TestAuthResolver(t *testing.T) {
	strategies := []string{
		"azurecert",
		"azurecreds",
		"device",
		"addin",
		"adfs",
		"fba",
		"ntlm",
		"saml",
		"tmg",
	}

	for _, strategy := range strategies {
		t.Run(strategy, func(t *testing.T) {
			cnfg, err := NewAuthCnfg(strategy, []byte("{}"))
			if err != nil {
				t.Error(err)
			}

			if cnfg.GetStrategy() != strategy {
				t.Errorf("strategy should be %s, but %s", strategy, cnfg.GetStrategy())
			}
		})
	}
}
