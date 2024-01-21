package api

import (
	"context"
	"testing"
)

func TestContextInfo(t *testing.T) {
	checkClient(t)

	sp := NewSP(spClient)

	t.Run("ContextInfo/SP", func(t *testing.T) {
		contextInfo, err := sp.ContextInfo(context.Background())
		if err != nil {
			t.Error(err)
		}
		if contextInfo.WebFullURL != sp.ToURL() {
			t.Error("incorrect web url")
		}
	})

	t.Run("ContextInfo/Web", func(t *testing.T) {
		contextInfo, err := sp.Web().ContextInfo(context.Background())
		if err != nil {
			t.Error(err)
		}
		if contextInfo.WebFullURL != sp.ToURL() {
			t.Error("incorrect web url")
		}
	})

	t.Run("ContextInfo/List", func(t *testing.T) {
		contextInfo, err := sp.Web().GetList("Shared Documents").ContextInfo(context.Background())
		if err != nil {
			t.Error(err)
		}
		if contextInfo.WebFullURL != sp.ToURL() {
			t.Error("incorrect web url")
		}
	})

}
