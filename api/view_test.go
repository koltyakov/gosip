package api

import (
	"testing"
)

func TestView(t *testing.T) {
	checkClient(t)

	web := NewSP(spClient).Web()
	listURI := getRelativeURL(spClient.AuthCnfg.GetSiteURL()) + "/Shared%20Documents"

	t.Run("Get", func(t *testing.T) {
		if _, err := web.GetList(listURI).Views().Get(); err != nil {
			t.Error(err)
		}
	})

}
