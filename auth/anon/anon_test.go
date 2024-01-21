package anon

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/koltyakov/gosip"
	h "github.com/koltyakov/gosip/test/helpers"
	u "github.com/koltyakov/gosip/test/utils"
)

var cnfgPath = "./config/private.anon.json"

func TestGettingAuthToken(t *testing.T) {
	if !h.ConfigExists(cnfgPath) {
		t.Skip("No auth config provided")
	}
	err := h.CheckAuth(
		&AuthCnfg{},
		cnfgPath,
		[]string{"SiteURL"},
	)
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
		filePath := u.ResolveCnfgPath("./test/tmp/anon.json")
		cnfg := &AuthCnfg{SiteURL: "test"}
		_ = os.MkdirAll(folderPath, os.ModePerm)
		if err := cnfg.WriteConfig(filePath); err != nil {
			t.Error(err)
		}
		_ = os.RemoveAll(filePath)
	})

	t.Run("GetStrategy", func(t *testing.T) {
		cnfg := &AuthCnfg{SiteURL: "test"}
		if cnfg.GetStrategy() != "anonymous" {
			t.Errorf(`wrong strategy name, expected "anonymous" got "%s"`, cnfg.GetStrategy())
		}
	})

	t.Run("GetStrategy", func(t *testing.T) {
		cnfg := &AuthCnfg{SiteURL: "http://test"}
		if cnfg.GetStrategy() != "anonymous" {
			t.Errorf(`wrong strategy name, expected "anonymous" got "%s"`, cnfg.GetStrategy())
		}
	})

	t.Run("GetSiteURL", func(t *testing.T) {
		cnfg := &AuthCnfg{SiteURL: "http://test"}
		if cnfg.GetSiteURL() != "http://test" {
			t.Errorf(`wrong strategy name, expected "http://test" got "%s"`, cnfg.GetSiteURL())
		}
	})

	t.Run("GetAuth", func(t *testing.T) {
		cnfg := &AuthCnfg{SiteURL: "http://test"}
		if _, _, err := cnfg.GetAuth(context.Background()); err != nil {
			t.Error(err)
		}
	})

	t.Run("SetAuth", func(t *testing.T) {
		cnfg := &AuthCnfg{SiteURL: "http://test"}
		req := &http.Request{}
		client := &gosip.SPClient{AuthCnfg: cnfg}
		if err := cnfg.SetAuth(req, client); err != nil {
			t.Error(err)
		}
	})

}
