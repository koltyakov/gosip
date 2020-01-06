package api

import (
	"testing"

	"github.com/google/uuid"
)

func TestRecycleBin(t *testing.T) {
	checkClient(t)

	sp := NewSP(spClient)
	newListTitle := uuid.New().String()
	if _, err := sp.Web().Lists().Add(newListTitle, nil); err != nil {
		t.Error(err)
	}
	list := sp.Web().Lists().GetByTitle(newListTitle)

	t.Run("Get/Site", func(t *testing.T) {
		if _, err := sp.Site().RecycleBin().Top(1).Get(); err != nil {
			t.Error(err)
		}
	})

	t.Run("Get/Web", func(t *testing.T) {
		_, err := sp.Web().RecycleBin().Top(1).Get()
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Restore", func(t *testing.T) {
		data, err := list.Items().Add([]byte(`{"Title":"Item"}`))
		if err != nil {
			t.Error(err)
		}
		if err := list.Items().GetByID(data.Data().ID).Recycle(); err != nil {
			t.Error(err)
		}
		items, err := sp.Web().RecycleBin().OrderBy("DeletedDate", false).Filter("LeafName eq '1_.000'").Top(1).Get()
		if err != nil {
			t.Error(err)
		}
		if err := sp.Web().RecycleBin().GetByID(items.Data()[0].Data().ID).Restore(); err != nil {
			t.Error(err)
		}
		data2, err := list.Items().Select("Id").Get()
		if err != nil {
			t.Error(err)
		}
		if len(data2.Data()) != 1 {
			t.Error("can't restore an item from recycle bin")
		}
	})

	if err := list.Delete(); err != nil {
		t.Error(err)
	}

}
