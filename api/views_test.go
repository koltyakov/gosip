package api

import (
	"encoding/json"
	"testing"
)

func TestViews(t *testing.T) {
	checkClient(t)

	web := NewSP(spClient).Web()
	listURI := getRelativeURL(spClient.AuthCnfg.GetSiteURL()) + "/Shared%20Documents"
	view, err := getAnyView()
	if err != nil {
		t.Error(err)
	}

	t.Run("Get", func(t *testing.T) {
		if _, err := web.GetList(listURI).Views().Get(); err != nil {
			t.Error(err)
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		if _, err := web.GetList(listURI).Views().GetByID(view.ID).Get(); err != nil {
			t.Error(err)
		}
	})

	t.Run("DefaultView", func(t *testing.T) {
		if _, err := web.GetList(listURI).Views().DefaultView().Get(); err != nil {
			t.Error(err)
		}
	})

	t.Run("GetByTitle", func(t *testing.T) {
		if _, err := web.GetList(listURI).Views().GetByTitle(view.Title).Get(); err != nil {
			t.Error(err)
		}
	})

}

func getAnyView() (*ViewInfo, error) {
	web := NewSP(spClient).Web()
	listURI := getRelativeURL(spClient.AuthCnfg.GetSiteURL()) + "/Shared%20Documents"
	data, err := web.GetList(listURI).Views().Top(1).Get()
	if err != nil {
		return nil, err
	}
	res := &struct {
		D struct {
			Results []*ViewInfo `json:"results"`
		} `json:"d"`
	}{}
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res.D.Results[0], nil
}
