package api

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
)

func TestWebs(t *testing.T) {
	checkClient(t)

	webs := NewSP(spClient).Web().Webs()
	endpoint := spClient.AuthCnfg.GetSiteURL() + "/_api/Web/Webs"
	newWebGUID := uuid.New().String()

	t.Run("Constructor", func(t *testing.T) {
		webs := NewWebs(spClient, endpoint, nil)
		if _, err := webs.Select("Id").Get(); err != nil {
			t.Error(err)
		}
	})

	t.Run("ToURL", func(t *testing.T) {
		if webs.ToURL() != endpoint {
			t.Errorf(
				"incorrect endpoint URL, expected \"%s\", received \"%s\"",
				endpoint,
				webs.ToURL(),
			)
		}
	})

	t.Run("Conf", func(t *testing.T) {
		webs.config = nil
		webs.Conf(headers.verbose)
		if webs.config != headers.verbose {
			t.Errorf("failed to apply config")
		}
	})

	t.Run("AddWeb", func(t *testing.T) {
		if !heavyTests {
			t.Skip("setup SPAPI_HEAVY_TESTS env var to \"true\" to run this test")
		}
		if _, err := webs.Add("CI: "+newWebGUID, "ci_"+newWebGUID, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("GetWebs", func(t *testing.T) {
		data, err := webs.Select("Id,Title").Conf(headers.verbose).Get()
		if err != nil {
			t.Error(err)
		}

		res := &struct {
			D struct {
				Results []struct {
					ID    string `json:"Id"`
					Title string `json:"Title"`
				} `json:"results"`
			} `json:"d"`
		}{}

		if err := json.Unmarshal(data, &res); err != nil {
			t.Error(err)
		}

		// if len(res.D.Results) == 0 {
		// 	t.Error("can't get webs")
		// }
	})

	t.Run("DeleteWeb", func(t *testing.T) {
		if !heavyTests {
			t.Skip("setup SPAPI_HEAVY_TESTS env var to \"true\" to run this test")
		}
		createdWebURL := spClient.AuthCnfg.GetSiteURL() + "/ci_" + newWebGUID
		subWeb := NewWeb(spClient, createdWebURL, nil)
		if _, err := subWeb.Delete(); err != nil {
			t.Error(err)
		}
	})

}
