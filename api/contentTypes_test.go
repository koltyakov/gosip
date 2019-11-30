package api

import (
	"encoding/json"
	"testing"
)

func TestContentTypes(t *testing.T) {
	checkClient(t)

	web := NewSP(spClient).Web()
	listURI := getRelativeURL(spClient.AuthCnfg.GetSiteURL()) + "/Shared%20Documents"
	contentType, err := getAnyContentType()
	if err != nil {
		t.Error(err)
	}

	t.Run("GetFromWeb", func(t *testing.T) {
		if _, err := web.ContentTypes().Conf(headers.verbose).Select("StringId").Top(1).Get(); err != nil {
			t.Error(err)
		}
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

func getAnyContentType() (*ContentTypeInfo, error) {
	web := NewSP(spClient).Web()
	data, err := web.ContentTypes().Conf(headers.verbose).Top(1).Get()
	if err != nil {
		return nil, err
	}
	res := &struct {
		D struct {
			Results []*ContentTypeInfo `json:"results"`
		} `json:"d"`
	}{}
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res.D.Results[0], nil
}
