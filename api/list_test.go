package api

import (
	"testing"
)

func TestList(t *testing.T) {
	checkClient(t)

	web := NewSP(spClient).Web()
	listID, _, err := getAnyList()
	if err != nil {
		t.Error(err)
	}
	list := web.Lists().GetByID(listID)

	t.Run("GetEntityType", func(t *testing.T) {
		entType, err := list.GetEntityType()
		if err != nil {
			t.Error(err)
		}
		if entType == "" {
			t.Error("can't get entity type")
		}
	})

	t.Run("Items", func(t *testing.T) {
		if _, err := list.Items().Select("Id").Top(1).Get(); err != nil {
			t.Error(err)
		}
	})

	t.Run("CurrentChangeToken", func(t *testing.T) {
		token, err := list.GetChangeToken()
		if err != nil {
			t.Error(err)
		}
		if token == "" {
			t.Error("can't get current change token")
		}
	})

}
