package saml

import (
	"os"
	"testing"

	h "github.com/koltyakov/gosip/test/helpers"
	u "github.com/koltyakov/gosip/test/utils"
)

var (
	cnfgPaths = []string{
		"./config/private.spo-user.json",
		"./config/private.spo-adfs.json",
	}
	ci bool
)

func init() {
	ci = os.Getenv("SPAUTH_CI") == "true"

	if ci { // In CI mode
		cnfgPath := "./config/private.spo-user.ci.json"
		auth := &AuthCnfg{
			SiteURL:  os.Getenv("SPAUTH_SITEURL"),
			Username: os.Getenv("SPAUTH_USERNAME"),
			Password: os.Getenv("SPAUTH_PASSWORD"),
		}
		auth.WriteConfig(u.ResolveCnfgPath(cnfgPath))
		cnfgPaths = []string{cnfgPath}
	}
}

func TestGettingAuthToken(t *testing.T) {
	configsChecked := 0
	for _, cnfgPath := range cnfgPaths {
		if !h.ConfigExists(cnfgPath) {
			continue
		}
		configsChecked++
		err := h.CheckAuth(
			&AuthCnfg{},
			cnfgPath,
			[]string{"SiteURL", "Username", "Password"},
		)
		if err != nil {
			t.Error(err)
		}
	}
	if configsChecked == 0 {
		t.Skip("No auth config(s) provided")
	}
}

func TestGettingDigest(t *testing.T) {
	configsChecked := 0
	for _, cnfgPath := range cnfgPaths {
		if !h.ConfigExists(cnfgPath) {
			continue
		}
		configsChecked++
		err := h.CheckDigest(&AuthCnfg{}, cnfgPath)
		if err != nil {
			t.Error(err)
		}
	}
	if configsChecked == 0 {
		t.Skip("No auth config(s) provided")
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
		os.MkdirAll(folderPath, os.ModePerm)
		if err := cnfg.WriteConfig(filePath); err != nil {
			t.Error(err)
		}
		os.RemoveAll(filePath)
	})

	t.Run("SetMasterkey", func(t *testing.T) {
		cnfg := &AuthCnfg{}
		cnfg.SetMasterkey("key")
		if cnfg.masterKey != "key" {
			t.Error("unable to set master key")
		}
	})

}
