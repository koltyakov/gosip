package api

import (
	"testing"
)

func TestList(t *testing.T) {
	checkClient(t)

	web := NewSP(spClient).Web()
	listInfo, err := getAnyList()
	if err != nil {
		t.Error(err)
	}
	list := web.Lists().GetByID(listInfo.ID)

	t.Run("GetEntityType", func(t *testing.T) {
		entType, err := list.GetEntityType()
		if err != nil {
			t.Error(err)
		}
		if entType == "" {
			t.Error("can't get entity type")
		}
	})

	t.Run("Get", func(t *testing.T) {
		l, err := list.Get() // .Select("*")
		if err != nil {
			t.Error(err)
		}
		if l.Data().Title == "" {
			t.Error("can't unmarshal list info")
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
