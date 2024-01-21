package api

import (
	"bytes"
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestRecycleBin(t *testing.T) {
	checkClient(t)

	sp := NewSP(spClient)
	newListTitle := uuid.New().String()
	if _, err := sp.Web().Lists().Add(context.Background(), newListTitle, nil); err != nil {
		t.Error(err)
	}
	list := sp.Web().Lists().GetByTitle(newListTitle)

	t.Run("Get/Site", func(t *testing.T) {
		data, err := sp.Site().RecycleBin().Top(1).Get(context.Background())
		if err != nil {
			t.Error(err)
		}
		if bytes.Compare(data, data.Normalized()) == -1 {
			t.Error("wrong response normalization")
		}
	})

	t.Run("Get/Web", func(t *testing.T) {
		_, err := sp.Web().RecycleBin().Top(1).Get(context.Background())
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Restore", func(t *testing.T) {
		data, err := list.Items().Add(context.Background(), []byte(`{"Title":"Item"}`))
		if err != nil {
			t.Error(err)
		}
		if err := list.Items().GetByID(data.Data().ID).Recycle(context.Background()); err != nil {
			t.Error(err)
		}
		items, err := sp.Web().RecycleBin().OrderBy("DeletedDate", false).Filter("LeafName eq '1_.000'").Top(1).Get(context.Background())
		if err != nil {
			t.Error(err)
		}
		itemID := items.Data()[0].Data().ID
		if _, err := sp.Web().RecycleBin().GetByID(itemID).Get(context.Background()); err != nil {
			t.Error(err)
		}
		if err := sp.Web().RecycleBin().GetByID(itemID).Restore(context.Background()); err != nil {
			t.Error(err)
		}
		data2, err := list.Items().Select("Id").Get(context.Background())
		if err != nil {
			t.Error(err)
		}
		if len(data2.Data()) != 1 {
			t.Error("can't restore an item from recycle bin")
		}
	})

	if err := list.Delete(context.Background()); err != nil {
		t.Error(err)
	}

}
