package auth

import (
	"os"
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
			cnfg, err := NewAuthByStrategy(strategy)
			if err != nil {
				t.Error(err)
			}

			if cnfg.GetStrategy() != strategy {
				t.Errorf("strategy should be %s, but %s", strategy, cnfg.GetStrategy())
			}
		})
	}
}

func TestAuthResolverError(t *testing.T) {
	_, err := NewAuthByStrategy("unknown")
	if err == nil {
		t.Error("should return an error")
	}
}

func TestAuthFileResolver(t *testing.T) {
	// Create temp file with auth config
	file, err := os.CreateTemp("", "private.json")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(file.Name())

	// Write auth config to the file
	_, err = file.Write([]byte(`{
		"strategy": "saml",
		"siteUrl": "https://contoso.sharepoint.com",
		"username": "user@contoso.onmicrosoft.com",
		"password": "00000000-0000-0000-0000-000000000000"
	}`))
	if err != nil {
		t.Error(err)
	}

	// Resolve auth config from the file
	cnfg, err := NewAuthFromFile(file.Name())
	if err != nil {
		t.Error(err)
	}

	if cnfg.GetStrategy() != "saml" {
		t.Errorf("strategy should be saml, but %s", cnfg.GetStrategy())
	}
}
