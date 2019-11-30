package api

import (
	"encoding/json"
	"testing"
)

func TestSite(t *testing.T) {
	checkClient(t)

	site := NewSP(spClient).Site()
	endpoint := spClient.AuthCnfg.GetSiteURL() + "/_api/Site"

	t.Run("Constructor", func(t *testing.T) {
		site := NewSite(spClient, endpoint, nil)
		if _, err := site.Select("Id").Get(); err != nil {
			t.Error(err)
		}
	})

	t.Run("ToURL", func(t *testing.T) {
		if site.ToURL() != endpoint {
			t.Errorf(
				"incorrect endpoint URL, expected \"%s\", received \"%s\"",
				endpoint,
				site.ToURL(),
			)
		}
	})

	t.Run("Conf", func(t *testing.T) {
		site.config = nil
		site.Conf(headers.verbose)
		if site.config != headers.verbose {
			t.Errorf("failed to apply config")
		}
	})

	t.Run("GetUrl", func(t *testing.T) {
		data, err := site.Select("Url").Conf(headers.verbose).Get()
		if err != nil {
			t.Error(err)
		}

		res := &struct {
			D struct {
				URL string `json:"Url"`
			} `json:"d"`
		}{}

		if err := json.Unmarshal(data, &res); err != nil {
			t.Error(err)
		}

		if res.D.URL == "" {
			t.Error("can't get site Url property")
		}
	})

	t.Run("RootWeb", func(t *testing.T) {
		data, err := site.RootWeb().Select("Title").Conf(headers.verbose).Get()
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
			t.Error("can't get root web title property")
		}
	})

	t.Run("OpenWebByID", func(t *testing.T) {
		data0, err := site.RootWeb().Select("Id").Conf(headers.verbose).Get()
		if err != nil {
			t.Error(err)
		}

		res := &struct {
			D struct {
				ID string `json:"ID"`
			} `json:"d"`
		}{}

		if err := json.Unmarshal(data0, &res); err != nil {
			t.Error(err)
		}

		if res.D.ID == "" {
			t.Error("can't get root web id property")
		}

		data, err := site.OpenWebByID(res.D.ID)
		if err != nil {
			t.Error(err)
		}

		if err := json.Unmarshal(data, &res); err != nil {
			t.Error(err)
		}

		if res.D.ID == "" {
			t.Error("can't open web by id property")
		}
	})

}
