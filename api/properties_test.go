package api

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestProperties(t *testing.T) {
	checkClient(t)

	sp := NewSP(spClient)
	web := sp.Web()
	webProps := web.AllProps()
	endpoint := spClient.AuthCnfg.GetSiteURL() + "/_api/Web/AllProperties"

	t.Run("Constructor", func(t *testing.T) {
		p := NewProperties(spClient, endpoint, nil, "web")
		if _, err := p.Select("vti_x005f_defaultlanguage").Get(); err != nil {
			t.Error(err)
		}
	})

	t.Run("ToURL", func(t *testing.T) {
		if webProps.ToURL() != endpoint {
			t.Errorf(
				"incorrect endpoint URL, expected \"%s\", received \"%s\"",
				endpoint,
				webProps.ToURL(),
			)
		}
	})

	t.Run("Get", func(t *testing.T) {
		if envCode == "2013" {
			t.Skip("is not supported with SP 2013")
		}

		data, err := webProps.Select("vti_x005f_defaultlanguage").Get()
		if err != nil {
			t.Error(err)
		}

		if data.Data()["vti_defaultlanguage"] == "" {
			t.Error("can't get web prop")
		}

		if bytes.Compare(data, data.Normalized()) == -1 {
			t.Error("wrong response normalization")
		}
	})

	t.Run("GetProps", func(t *testing.T) {
		data, err := webProps.GetProps([]string{"vti_defaultlanguage"})
		if err != nil {
			t.Error(err)
		}

		if data["vti_defaultlanguage"] == "" {
			t.Error("can't get web prop")
		}
	})

	t.Run("GetMultipleProps", func(t *testing.T) {
		data, err := webProps.GetProps([]string{"vti_defaultlanguage", "vti_associategroups"})
		if err != nil {
			t.Error(err)
		}

		if len(data) != 2 {
			t.Error("wrong props number")
		}
	})

	t.Run("GetNonExistongProps", func(t *testing.T) {
		_, err := webProps.GetProps([]string{"vti_defaultlanguage", "vti_associategroups_do_not_exist"})
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Set", func(t *testing.T) {
		if err := webProps.Set("test_gosip", time.Now().String()); err != nil {
			// By default is denied on Modern SPO sites, so ignore in tests
			if !strings.Contains(err.Error(), "System.UnauthorizedAccessException") {
				t.Error(err)
			}
		}
		if err := web.RootFolder().Props().Set("test_gosip", time.Now().String()); err != nil {
			// By default is denied on Modern SPO sites, so ignore in tests
			if !strings.Contains(err.Error(), "System.UnauthorizedAccessException") {
				t.Error(err)
			}
		}
	})

	t.Run("SetProps", func(t *testing.T) {
		if err := webProps.SetProps(map[string]string{
			"test_gosip_prop1": time.Now().String(),
			"test_gosip_prop2": time.Now().String(),
		}); err != nil {
			// By default is denied on Modern SPO sites, so ignore in tests
			if !strings.Contains(err.Error(), "System.UnauthorizedAccessException") {
				t.Error(err)
			}
		}
	})

	t.Run("PrintNoScriptWarning", func(t *testing.T) {
		printNoScriptWarning("https://contoso.sharepoint.com", fmt.Errorf("System.UnauthorizedAccessException"))
	})

}
