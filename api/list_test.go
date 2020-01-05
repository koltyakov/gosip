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

	t.Run("ParentWeb", func(t *testing.T) {
		if _, err := list.ParentWeb().Get(); err != nil {
			t.Error(err)
		}
	})

	t.Run("ReserveListItemID", func(t *testing.T) {
		nextID, err := list.ReserveListItemID()
		if err != nil {
			t.Error(err)
		}
		if nextID == 0 {
			t.Error("can't reserve list item ID")
		}
	})

	t.Run("RenderListData", func(t *testing.T) {
		listData, err := list.RenderListData(`<View></View>`)
		if err != nil {
			t.Error(err)
		}
		if listData.Data().FolderPermissions == "" {
			t.Error("incorrect data")
		}
	})

}
