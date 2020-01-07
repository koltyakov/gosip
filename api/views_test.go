package api

import (
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

	t.Run("Conf", func(t *testing.T) {
		v := web.GetList(listURI).Views()
		hs := map[string]*RequestConfig{
			"nometadata":      HeadersPresets.Nometadata,
			"minimalmetadata": HeadersPresets.Minimalmetadata,
			"verbose":         HeadersPresets.Verbose,
		}
		for key, preset := range hs {
			g := v.Conf(preset)
			if g.config != preset {
				t.Errorf("can't %v config", key)
			}
		}
	})

	t.Run("Modifiers", func(t *testing.T) {
		mods := web.GetList(listURI).Views().
			Select("*").Expand("*").Filter("*").OrderBy("*", true).
			modifiers
		if mods == nil || len(mods.mods) != 4 {
			t.Error("can't add modifiers")
		}
	})

	t.Run("Get", func(t *testing.T) {
		data, err := web.GetList(listURI).Views().Get()
		if err != nil {
			t.Error(err)
		}
		if data.Data()[0].Data().ID == "" {
			t.Error("can't unmarshal data")
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

	// ToDo:
	// Add

}

func getAnyView() (*ViewInfo, error) {
	web := NewSP(spClient).Web()
	listURI := getRelativeURL(spClient.AuthCnfg.GetSiteURL()) + "/Shared%20Documents"
	data, err := web.GetList(listURI).Views().Top(1).Get()
	if err != nil {
		return nil, err
	}
	return data.Data()[0].Data(), nil
}
