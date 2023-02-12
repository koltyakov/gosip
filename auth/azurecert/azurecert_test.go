package azurecert

import (
	"io/ioutil"
	"os"
	"testing"

	h "github.com/koltyakov/gosip/test/helpers"
	u "github.com/koltyakov/gosip/test/utils"
)

var cnfgPath = "./config/private.azurecert.json"

func TestGettingAuthToken(t *testing.T) {
	if !h.ConfigExists(cnfgPath) {
		t.Skip("No auth config provided")
	}
	err := h.CheckAuth(
		&AuthCnfg{},
		cnfgPath,
		[]string{"SiteURL", "TenantID", "ClientID", "CertPath", "CertPass"},
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

func TestRequest(t *testing.T) {
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
		folderPath := u.ResolveCnfgPath("./tmp")
		filePath := u.ResolveCnfgPath("./tmp/private.azurecert.malformed.json")
		_ = os.MkdirAll(folderPath, os.ModePerm)
		_ = ioutil.WriteFile(filePath, []byte("not a json"), 0644)
		if err := cnfg.ReadConfig(filePath); err == nil {
			t.Error("malformed config should not pass")
		}
		_ = os.RemoveAll(filePath)
	})

	t.Run("WriteConfig", func(t *testing.T) {
		folderPath := u.ResolveCnfgPath("./tmp")
		filePath := u.ResolveCnfgPath("./tmp/private.azurecert.json")
		cnfg := &AuthCnfg{
			SiteURL: "test",
		}
		_ = os.MkdirAll(folderPath, os.ModePerm)
		if err := cnfg.WriteConfig(filePath); err != nil {
			t.Error(err)
		}
		_ = os.RemoveAll(filePath)
	})
}

func TestCheckTransport(t *testing.T) {
	if !h.ConfigExists(cnfgPath) {
		t.Skip("No auth config provided")
	}
	err := h.CheckTransport(&AuthCnfg{}, cnfgPath)
	if err != nil {
		t.Error(err)
	}
}
