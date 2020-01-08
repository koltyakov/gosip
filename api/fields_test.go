package api

import (
	"fmt"
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

	t.Run("Conf", func(t *testing.T) {
		f := web.Fields()
		hs := map[string]*RequestConfig{
			"nometadata":      HeadersPresets.Nometadata,
			"minimalmetadata": HeadersPresets.Minimalmetadata,
			"verbose":         HeadersPresets.Verbose,
		}
		for key, preset := range hs {
			g := f.Conf(preset)
			if g.config != preset {
				t.Errorf("can't %v config", key)
			}
		}
	})

	t.Run("Modifiers", func(t *testing.T) {
		f := web.Fields()
		mods := f.Select("*").Expand("*").Filter("*").Top(1).OrderBy("*", true).modifiers
		if mods == nil || len(mods.mods) != 5 {
			t.Error("can't add modifiers")
		}
	})

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

func getAnyField() (*GenericFieldInfo, error) {
	web := NewSP(spClient).Web()
	data, err := web.Fields().Top(1).Get()
	if err != nil {
		return nil, err
	}
	if len(data.Data()) == 0 {
		return nil, fmt.Errorf("can't get random field")
	}
	return data.Data()[0].Data(), nil
}
