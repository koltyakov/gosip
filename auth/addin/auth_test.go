package addin

import (
	"os"
	"testing"

	h "github.com/koltyakov/gosip/test/helpers"
	u "github.com/koltyakov/gosip/test/utils"
)

var (
	cnfgPath = "./config/private.spo-addin.json"
	ci       bool
)

func init() {
	ci = os.Getenv("SPAUTH_CI") == "true"

	if ci { // In CI mode
		cnfgPath = "./config/private.spo-addin.ci.json"
		auth := &AuthCnfg{
			SiteURL:      os.Getenv("SPAUTH_SITEURL"),
			ClientID:     os.Getenv("SPAUTH_CLIENTID"),
			ClientSecret: os.Getenv("SPAUTH_CLIENTSECRET"),
		}
		_ = auth.WriteConfig(u.ResolveCnfgPath(cnfgPath))
	}
}

func TestGettingAuthToken(t *testing.T) {
	if !h.ConfigExists(cnfgPath) {
		t.Skip("No auth config provided")
	}
	err := h.CheckAuth(
		&AuthCnfg{},
		cnfgPath,
		[]string{"SiteURL", "ClientID", "ClientSecret"},
	)
	if err != nil {
		t.Error(err)
	}
}

func TestGettingDigest(t *testing.T) {
	if !h.ConfigExists(cnfgPath) {
		t.Skip("No auth config provided")
	}
	err := h.CheckDigest(&AuthCnfg{}, cnfgPath)
	if err != nil {
		t.Error(err)
	}
}

func TestCheckRequest(t *testing.T) {
	if !h.ConfigExists(cnfgPath) {
		t.Skip("No auth config provided")
	}
	err := h.CheckRequest(&AuthCnfg{}, cnfgPath)
	if err != nil {
		t.Error(err)
	}
}

func TestAuthEdgeCases(t *testing.T) {

	t.Run("ReadConfig/MissedConfig", func(t *testing.T) {
		cnfg := &AuthCnfg{}
		if err := cnfg.ReadConfig("wrong_path.json"); err == nil {
			t.Error("wrong_path config should not pass")
		}
	})

	t.Run("ReadConfig/MissedConfig", func(t *testing.T) {
		cnfg := &AuthCnfg{}
		if err := cnfg.ReadConfig(u.ResolveCnfgPath("./test/config/malformed.json")); err == nil {
			t.Error("malformed config should not pass")
		}
	})

	t.Run("WriteConfig", func(t *testing.T) {
		folderPath := u.ResolveCnfgPath("./test/tmp")
		filePath := u.ResolveCnfgPath("./test/tmp/addin.json")
		cnfg := &AuthCnfg{SiteURL: "test"}
		_ = os.MkdirAll(folderPath, os.ModePerm)
		if err := cnfg.WriteConfig(filePath); err != nil {
			t.Error(err)
		}
		_ = os.RemoveAll(filePath)
	})

	t.Run("SetMasterkey", func(t *testing.T) {
		cnfg := &AuthCnfg{}
		cnfg.SetMasterkey("key")
		if cnfg.masterKey != "key" {
			t.Error("unable to set master key")
		}
	})

}
