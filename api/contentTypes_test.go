package api

import (
	"encoding/json"
	"testing"
)

func TestContentTypes(t *testing.T) {
	checkClient(t)

	web := NewSP(spClient).Web()
	listURI := getRelativeURL(spClient.AuthCnfg.GetSiteURL()) + "/Shared%20Documents"
	contentType := &ContentTypeInfo{}

	t.Run("GetFromWeb", func(t *testing.T) {
		data, err := web.ContentTypes().Conf(headers.verbose).Select("StringId").Top(1).Get()
		if err != nil {
			t.Error(err)
		}
		res := &struct {
			D struct {
				Results []*ContentTypeInfo `json:"results"`
			} `json:"d"`
		}{}
		if err := json.Unmarshal(data, &res); err != nil {
			t.Error(err)
		}
		contentType = res.D.Results[0]
	})

	t.Run("GetFromList", func(t *testing.T) {
		if _, err := web.GetList(listURI).ContentTypes().Top(1).Get(); err != nil {
			t.Error(err)
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		if _, err := web.ContentTypes().GetByID(contentType.ID).Get(); err != nil {
			t.Error(err)
		}
	})

}
