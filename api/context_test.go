package api

import (
	"testing"
)

func TestContextInfo(t *testing.T) {
	checkClient(t)

	sp := NewSP(spClient)

	t.Run("ContextInfo/SP", func(t *testing.T) {
		contextInfo, err := sp.ContextInfo()
		if err != nil {
			t.Error(err)
		}
		if contextInfo.WebFullURL != sp.ToURL() {
			t.Error("incorrect web url")
		}
	})

	t.Run("ContextInfo/Web", func(t *testing.T) {
		contextInfo, err := sp.Web().ContextInfo()
		if err != nil {
			t.Error(err)
		}
		if contextInfo.WebFullURL != sp.ToURL() {
			t.Error("incorrect web url")
		}
	})

	t.Run("ContextInfo/List", func(t *testing.T) {
		contextInfo, err := sp.Web().GetList("Shared Documents").ContextInfo()
		if err != nil {
			t.Error(err)
		}
		if contextInfo.WebFullURL != sp.ToURL() {
			t.Error("incorrect web url")
		}
	})

}
