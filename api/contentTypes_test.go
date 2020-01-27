package api

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/google/uuid"
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
		data, err := web.ContentTypes().Select("StringId").Top(1).Get()
		if err != nil {
			t.Error(err)
		}
		if bytes.Compare(data, data.Normalized()) == -1 {
			t.Error("response normalization error")
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

	t.Run("CreateUsingParentID", func(t *testing.T) {
		guid := uuid.New().String()
		newCTID, err := web.ContentTypes().Create(&ContentTypeCreationInfo{
			Name:                guid,
			Group:               "Test",
			ParentContentTypeID: "0x01",
		})
		if err != nil {
			t.Error(err)
		}
		if newCTID == "" {
			t.Error("can't parse CT creation response")
		}
		if err := web.ContentTypes().GetByID(newCTID).Delete(); err != nil {
			t.Error(err)
		}
	})

	t.Run("CreateUsingID", func(t *testing.T) {
		guid := uuid.New().String()
		newCTID, err := web.ContentTypes().Create(&ContentTypeCreationInfo{
			ID:    "0x0100" + strings.ToUpper(strings.Replace(guid, "-", "", -1)),
			Name:  guid,
			Group: "Test",
		})
		if err != nil {
			t.Error(err)
		}
		if newCTID == "" {
			t.Error("can't parse CT creation response")
		}
		if err := web.ContentTypes().GetByID(newCTID).Delete(); err != nil {
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
