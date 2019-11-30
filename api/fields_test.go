package api

import (
	"encoding/json"
	"testing"
)

func TestFields(t *testing.T) {
	checkClient(t)

	web := NewSP(spClient).Web()
	listURI := getRelativeURL(spClient.AuthCnfg.GetSiteURL()) + "/Shared%20Documents"
	field, err := getAnyField()
	if err != nil {
		t.Error(err)
	}

	t.Run("GetFromWeb", func(t *testing.T) {
		if _, err := web.Fields().Top(1).Get(); err != nil {
			t.Error(err)
		}
	})

	t.Run("GetFromList", func(t *testing.T) {
		if _, err := web.GetList(listURI).Fields().Top(1).Get(); err != nil {
			t.Error(err)
		}
	})

	t.Run("GetFromContentType", func(t *testing.T) {
		contentType, err := getAnyContentType()
		if err != nil {
			t.Error(err)
		}
		if _, err := web.ContentTypes().GetByID(contentType.ID).Fields().Get(); err != nil {
			t.Error(err)
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		if _, err := web.Fields().GetByID(field.ID).Get(); err != nil {
			t.Error(err)
		}
	})

	t.Run("GetByTitle", func(t *testing.T) {
		if _, err := web.Fields().GetByTitle(field.Title).Get(); err != nil {
			t.Error(err)
		}
	})

	t.Run("GetByInternalNameOrTitle", func(t *testing.T) {
		if _, err := web.Fields().GetByInternalNameOrTitle(field.InternalName).Get(); err != nil {
			t.Error(err)
		}
		if _, err := web.Fields().GetByInternalNameOrTitle(field.Title).Get(); err != nil {
			t.Error(err)
		}
	})

}

func getAnyField() (*FieldInfo, error) {
	web := NewSP(spClient).Web()
	data, err := web.Fields().Conf(headers.verbose).Top(1).Get()
	if err != nil {
		return nil, err
	}
	res := &struct {
		D struct {
			Results []*FieldInfo `json:"results"`
		} `json:"d"`
	}{}
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res.D.Results[0], nil
}
