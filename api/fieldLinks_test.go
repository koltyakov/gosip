package api

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestFieldLinks(t *testing.T) {
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

	t.Run("Conf", func(t *testing.T) {
		fl := web.ContentTypes().GetByID(ctID).FieldLinks()
		hs := map[string]*RequestConfig{
			"nometadata":      HeadersPresets.Nometadata,
			"minimalmetadata": HeadersPresets.Minimalmetadata,
			"verbose":         HeadersPresets.Verbose,
		}
		for key, preset := range hs {
			g := fl.Conf(preset)
			if g.config != preset {
				t.Errorf("can't %v config", key)
			}
		}
	})

	t.Run("Modifiers", func(t *testing.T) {
		fl := web.ContentTypes().GetByID(ctID).FieldLinks()
		mods := fl.Select("*").Filter("*").Top(1).modifiers
		if mods == nil || len(mods.mods) != 3 {
			t.Error("wrong number of modifiers")
		}
	})

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
		if resp.Data()[0].ID == "" {
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

	// t.Run("Add", func(t *testing.T) {
	// 	resp, err := web.ContentTypes().GetByID(ctID).FieldLinks().Add("Language", false, false)
	// 	if err != nil {
	// 		t.Error(err)
	// 	}
	// 	fmt.Println(resp)
	// })

	if err := web.ContentTypes().GetByID(ctID).Delete(); err != nil {
		t.Error(err)
	}

}
