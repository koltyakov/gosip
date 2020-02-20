package adfs

import (
	"os"
	"testing"

	h "github.com/koltyakov/gosip/test/helpers"
	u "github.com/koltyakov/gosip/test/utils"
)

var (
	cnfgPaths = []string{
		"./config/private.onprem-adfs.json",
		"./config/private.onprem-wap.json",
		"./config/private.onprem-wap-adfs.json",
	}
)

func TestGettingAuthToken(t *testing.T) {
	for _, cnfgPath := range cnfgPaths {
		if !h.ConfigExists(cnfgPath) {
			continue
		}
		err := h.CheckAuth(
			&AuthCnfg{},
			cnfgPath,
			[]string{"SiteURL", "Username", "Password"},
		)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestGettingDigest(t *testing.T) {
	for _, cnfgPath := range cnfgPaths {
		if !h.ConfigExists(cnfgPath) {
			continue
		}
		err := h.CheckDigest(&AuthCnfg{}, cnfgPath)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestCheckRequest(t *testing.T) {
	for _, cnfgPath := range cnfgPaths {
		if !h.ConfigExists(cnfgPath) {
			continue
		}
		err := h.CheckRequest(&AuthCnfg{}, cnfgPath)
		if err != nil {
			t.Error(err)
		}
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
		filePath := u.ResolveCnfgPath("./test/tmp/adfs.json")
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
