package api

import (
	"fmt"
	"testing"
)

var ctID = ""

func TestContentType(t *testing.T) {
	checkClient(t)

	web := NewSP(spClient).Web()

	t.Run("Conf", func(t *testing.T) {
		ct := web.ContentTypes().GetByID("")
		hs := map[string]*RequestConfig{
			"nometadata":      HeadersPresets.Nometadata,
			"minimalmetadata": HeadersPresets.Minimalmetadata,
			"verbose":         HeadersPresets.Verbose,
		}
		for key, preset := range hs {
			g := ct.Conf(preset)
			if g.config != preset {
				t.Errorf("can't %v config", key)
			}
		}
	})

	t.Run("Modifiers", func(t *testing.T) {
		ct, err := getRandomCT()
		if err != nil {
			t.Error(err)
		}
		if _, err := ct.Select("*,Fields/*").Expand("Fields").Get(); err != nil {
			t.Error(err)
		}
	})

	t.Run("Get", func(t *testing.T) {
		resp, err := web.ContentTypes().Top(5).Get()
		if err != nil {
			t.Error(err)
		}
		cts := resp.Data()
		if len(cts) != 5 {
			t.Error("wrong number of content types")
		}
		if cts[0].Data().ID == "" {
			t.Error("can't get content type info")
		}
		if _, err := web.ContentTypes().GetByID(cts[0].Data().ID).Get(); err != nil {
			t.Error(err)
		}
	})

	t.Run("FromURL", func(t *testing.T) {
		ct, err := getRandomCT()
		if err != nil {
			t.Error(err)
		}
		ctR, err := ct.Get()
		if err != nil {
			t.Error(err)
		}
		entityURL := ExtractEntityURI(ctR)
		if entityURL == "" {
			t.Error("can't extract entity URL")
		}
		ct1, err := web.ContentTypes().GetByID("").FromURL(entityURL).Get()
		if err != nil {
			t.Error(err)
		}
		if ctR.Data().ID != ct1.Data().ID {
			t.Error("can't get CT from entity URL")
		}
	})

	// ToDo:
	// Update
	// Delete
	// Recycle

}

func getRandomCT() (*ContentType, error) {
	sp := NewSP(spClient)
	if ctID == "" {
		resp, err := sp.Web().ContentTypes().Top(1).Get()
		if err != nil {
			return nil, err
		}
		cts := resp.Data()
		if len(cts) != 1 {
			return nil, fmt.Errorf("wrong number of content types")
		}
		if cts[0].Data().ID == "" {
			return nil, fmt.Errorf("can't get content type info")
		}
		ctID = cts[0].Data().ID
	}
	return sp.Web().ContentTypes().GetByID(ctID), nil
}
