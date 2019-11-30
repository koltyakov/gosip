package api

import (
	"encoding/json"
	"testing"
)

func TestWebProps(t *testing.T) {
	checkClient(t)

	webProps := NewSP(spClient).Web().Props()
	endpoint := spClient.AuthCnfg.GetSiteURL() + "/_api/Web/AllProperties"

	t.Run("Constructor", func(t *testing.T) {
		webProps := NewWebProps(spClient, endpoint, nil)
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
		data, err := webProps.Select("vti_x005f_defaultlanguage").Conf(headers.verbose).Get()
		if err != nil {
			t.Error(err)
		}

		res := &struct {
			D struct {
				Lang string `json:"vti_x005f_defaultlanguage"`
			} `json:"d"`
		}{}

		if err := json.Unmarshal(data, &res); err != nil {
			t.Error(err)
		}

		if res.D.Lang == "" {
			t.Error("can't get web prop")
		}
	})

}
