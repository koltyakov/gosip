package api

import (
	"fmt"
	"testing"
)

var fieldID = ""

func TestField(t *testing.T) {
	checkClient(t)

	web := NewSP(spClient).Web()
	field, err := getRandomField()
	if err != nil {
		t.Error(err)
	}

	t.Run("Conf", func(t *testing.T) {
		f := web.Fields().GetByID("")
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

	t.Run("Get", func(t *testing.T) {
		data, err := field.Select("Id").Get()
		if err != nil {
			t.Error(err)
		}
		if data.Data().ID == "" {
			t.Error("can't unmarshal data")
		}
	})

	t.Run("Modifiers", func(t *testing.T) {
		f, err := getRandomField()
		if err != nil {
			t.Error(err)
		}
		mods := f.Select("*").Expand("*").modifiers
		if mods == nil || len(mods.mods) != 2 {
			t.Error("can't add modifiers")
		}
	})

	t.Run("FromURL", func(t *testing.T) {
		fieldR, err := field.Get()
		if err != nil {
			t.Error(err)
		}
		entityURL := ExtractEntityURI(fieldR)
		if entityURL == "" {
			t.Error("can't extract entity URL")
		}
		field1, err := web.Fields().GetByID("").FromURL(entityURL).Get()
		if err != nil {
			t.Error(err)
		}
		if fieldR.Data().ID != field1.Data().ID {
			t.Error("can't get CT from entity URL")
		}
	})

	// ToDo:
	// Update
	// Delete
	// Recycle

}

func getRandomField() (*Field, error) {
	sp := NewSP(spClient)
	if fieldID == "" {
		resp, err := sp.Web().Fields().Top(1).Get()
		if err != nil {
			return nil, err
		}
		cts := resp.Data()
		if len(cts) != 1 {
			return nil, fmt.Errorf("wrong number of fields")
		}
		if cts[0].Data().ID == "" {
			return nil, fmt.Errorf("can't get field info")
		}
		fieldID = cts[0].Data().ID
	}
	return sp.Web().Fields().GetByID(fieldID), nil
}
