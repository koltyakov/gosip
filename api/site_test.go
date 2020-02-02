package api

import (
	"bytes"
	"testing"
)

func TestSite(t *testing.T) {
	checkClient(t)

	sp := NewSP(spClient)
	site := sp.Site()
	endpoint := spClient.AuthCnfg.GetSiteURL() + "/_api/Site"

	t.Run("Constructor", func(t *testing.T) {
		s := NewSite(spClient, endpoint, nil)
		if _, err := s.Select("Id").Get(); err != nil {
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

	t.Run("GetURL", func(t *testing.T) {
		data, err := site.Select("Url").Conf(headers.verbose).Get()
		if err != nil {
			t.Error(err)
		}
		if data.Data().URL == "" {
			t.Error("can't get site Url property")
		}
		if bytes.Compare(data, data.Normalized()) == -1 {
			t.Error("wrong response normalization")
		}
	})

	t.Run("FromURL", func(t *testing.T) {
		s := site.FromURL("site_url")
		if s.endpoint != "site_url" {
			t.Error("can't get site from url")
		}
	})

	t.Run("RootWeb", func(t *testing.T) {
		data, err := site.RootWeb().Select("Title").Conf(headers.verbose).Get()
		if err != nil {
			t.Error(err)
		}
		if data.Data().Title == "" {
			t.Error("can't get root web title property")
		}
	})

	t.Run("OpenWebByID", func(t *testing.T) {
		data0, err := site.RootWeb().Select("Id").Conf(headers.verbose).Get()
		if err != nil {
			t.Error(err)
		}
		if data0.Data().ID == "" {
			t.Error("can't get root web id property")
		}

		data, err := site.OpenWebByID(data0.Data().ID)
		if err != nil {
			t.Error(err)
		}
		if data.Data().ID == "" {
			t.Error("can't open web by id property")
		}
	})

	t.Run("Owner", func(t *testing.T) {
		if _, err := site.Owner().Get(); err != nil {
			t.Error(err)
		}
	})

}
