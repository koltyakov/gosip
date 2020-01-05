package api

import (
	"strings"
	"testing"
	"time"
)

func TestProperties(t *testing.T) {
	checkClient(t)

	web := NewSP(spClient).Web()
	webProps := web.AllProps()
	endpoint := spClient.AuthCnfg.GetSiteURL() + "/_api/Web/AllProperties"

	t.Run("Constructor", func(t *testing.T) {
		webProps := NewProperties(spClient, endpoint, nil)
		if _, err := webProps.Select("vti_x005f_defaultlanguage").Get(); err != nil {
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

	t.Run("Conf", func(t *testing.T) {
		webProps.config = nil
		webProps.Conf(headers.verbose)
		if webProps.config != headers.verbose {
			t.Errorf("failed to apply config")
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

	t.Run("Set", func(t *testing.T) {
		if _, err := webProps.Set("test_gosip", time.Now().String()); err != nil {
			// By default is denied on Modern SPO sites, so ignore in tests
			if strings.Index(err.Error(), "System.UnauthorizedAccessException") == -1 {
				t.Error(err)
			}
		}
		if _, err := web.RootFolder().Props().Set("test_gosip", time.Now().String()); err != nil {
			// By default is denied on Modern SPO sites, so ignore in tests
			if strings.Index(err.Error(), "System.UnauthorizedAccessException") == -1 {
				t.Error(err)
			}
		}
	})

	t.Run("SetProps", func(t *testing.T) {
		if _, err := webProps.SetProps(map[string]string{
			"test_gosip_prop1": time.Now().String(),
			"test_gosip_prop2": time.Now().String(),
		}); err != nil {
			// By default is denied on Modern SPO sites, so ignore in tests
			if strings.Index(err.Error(), "System.UnauthorizedAccessException") == -1 {
				t.Error(err)
			}
		}
	})

}
