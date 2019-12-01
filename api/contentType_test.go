package api

import (
	"testing"
)

func TestContentType(t *testing.T) {
	checkClient(t)

	web := NewSP(spClient).Web()

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

}
