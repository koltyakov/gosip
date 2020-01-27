package api

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestFieldLinks(t *testing.T) {
	t.Parallel()
	checkClient(t)

	web := NewSP(spClient).Web()
	guid := uuid.New().String()
	ctID := "0x0100" + strings.ToUpper(strings.Replace(guid, "-", "", -1))
	ct := []byte(TrimMultiline(`{
		"Group":"Custom Content Types",
		"Id": {"StringValue":"` + ctID + `"},
		"Name":"test-temp-ct ` + guid + `"
	}`))
	ctResp, err := web.ContentTypes().Add(ct)
	if err != nil {
		t.Error(err)
	}
	ctID = ctResp.Data().ID // content type ID can't be set in REST API https://github.com/pnp/pnpjs/issues/457

	t.Run("Get", func(t *testing.T) {
		resp, err := web.ContentTypes().GetByID(ctID).FieldLinks().Get()
		if err != nil {
			t.Error(err)
		}
		if len(resp.Data()) == 0 {
			t.Error("can't get field links")
		}
		if bytes.Compare(resp, resp.Normalized()) == -1 {
			t.Error("response normalization error")
		}
		if resp.Data()[0].Data().ID == "" {
			t.Error("can't unmarshal field info")
		}
	})

	t.Run("GetFields", func(t *testing.T) {
		resp, err := web.ContentTypes().GetByID(ctID).FieldLinks().GetFields()
		if err != nil {
			t.Error(err)
		}
		if len(resp.Data()) == 0 {
			t.Error("can't get fields")
		}
		if bytes.Compare(resp, resp.Normalized()) == -1 {
			t.Error("response normalization error")
		}
		if resp.Data()[0].Data().SchemaXML == "" {
			t.Error("can't unmarshal field info")
		}
	})

	t.Run("Add", func(t *testing.T) {
		fl, err := web.ContentTypes().GetByID(ctID).FieldLinks().Add("Language")
		if err != nil {
			t.Error(err)
		}
		if fl == "" {
			t.Error("can't parse field link add response")
		}
		fls, _ := web.ContentTypes().GetByID(ctID).FieldLinks().Get()
		if len(fls.Data()) < 3 {
			t.Error("failed adding field link")
		}
		if err := web.ContentTypes().GetByID(ctID).FieldLinks().GetByID(fl).Delete(); err != nil {
			t.Error(err)
		}
	})

	if err := web.ContentTypes().GetByID(ctID).Delete(); err != nil {
		t.Error(err)
	}

}
