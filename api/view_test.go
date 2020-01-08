package api

import (
	"bytes"
	"testing"
)

func TestView(t *testing.T) {
	checkClient(t)

	web := NewSP(spClient).Web()
	listURI := getRelativeURL(spClient.AuthCnfg.GetSiteURL()) + "/Shared%20Documents"

	t.Run("Conf", func(t *testing.T) {
		v := web.GetList(listURI).Views().DefaultView()
		hs := map[string]*RequestConfig{
			"nometadata":      HeadersPresets.Nometadata,
			"minimalmetadata": HeadersPresets.Minimalmetadata,
			"verbose":         HeadersPresets.Verbose,
		}
		for key, preset := range hs {
			lv := v.Conf(preset)
			if lv.config != preset {
				t.Errorf("can't %v config", key)
			}
		}
	})

	t.Run("Get", func(t *testing.T) {
		data, err := web.GetList(listURI).Views().DefaultView().Get()
		if err != nil {
			t.Error(err)
		}
		if data.Data().ID == "" {
			t.Error("can't unmarshal data")
		}
		if bytes.Compare(data, data.Normalized()) == -1 {
			t.Error("wrong response normalization")
		}
	})

	t.Run("Modifiers", func(t *testing.T) {
		mods := web.GetList(listURI).Views().DefaultView().
			Select("*").Expand("*").modifiers
		if mods == nil || len(mods.mods) != 2 {
			t.Error("can't add modifiers")
		}
	})

	// ToDo:
	// Update
	// Delete
	// Recycle

}
