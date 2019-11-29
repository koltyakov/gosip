package api

import (
	"encoding/json"
	"testing"
)

func TestWeb(t *testing.T) {
	checkClient(t)

	web := NewSP(spClient).Web()
	endpoint := spClient.AuthCnfg.GetSiteURL() + "/_api/Web"

	t.Run("Constructor", func(t *testing.T) {
		web := NewWeb(spClient, endpoint, nil)
		if _, err := web.Get(); err != nil {
			t.Error(err)
		}
	})

	t.Run("ToURL", func(t *testing.T) {
		if web.ToURL() != endpoint {
			t.Errorf(
				"incorrect endpoint URL, expected \"%s\", received \"%s\"",
				endpoint,
				web.ToURL(),
			)
		}
	})

	t.Run("GetTitle", func(t *testing.T) {
		data, err := web.Select("Title").Conf(headers.verbose).Get()
		if err != nil {
			t.Error(err)
		}

		res := &struct {
			D struct {
				Title string `json:"Title"`
			} `json:"d"`
		}{}

		if err := json.Unmarshal(data, &res); err != nil {
			t.Error(err)
		}

		if res.D.Title == "" {
			t.Error("can't get web title property")
		}
	})

	t.Run("NoTitle", func(t *testing.T) {
		data, err := web.Select("Id").Conf(headers.verbose).Get()
		if err != nil {
			t.Error(err)
		}

		res := &struct {
			D struct {
				Title string `json:"Title"`
			} `json:"d"`
		}{}

		if err := json.Unmarshal(data, &res); err != nil {
			t.Error(err)
		}

		if res.D.Title != "" {
			t.Error("can't get web title property")
		}
	})

}
