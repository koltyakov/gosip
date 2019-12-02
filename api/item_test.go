package api

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestItem(t *testing.T) {
	checkClient(t)

	web := NewSP(spClient).Web()
	newListTitle := uuid.New().String()
	if _, err := web.Lists().Add(newListTitle, nil); err != nil {
		t.Error(err)
	}
	list := web.Lists().GetByTitle(newListTitle)
	entType, err := list.GetEntityType()
	if err != nil {
		t.Error(err)
	}

	t.Run("AddSeries", func(t *testing.T) {
		for i := 1; i < 5; i++ {
			metadata := make(map[string]interface{})
			metadata["__metadata"] = map[string]string{"type": entType}
			metadata["Title"] = fmt.Sprintf("Item %d", i)
			body, _ := json.Marshal(metadata)
			if _, err := list.Items().Add(body); err != nil {
				t.Error(err)
			}
		}
	})

	t.Run("Get", func(t *testing.T) {
		item, err := list.Items().GetByID(1).Get()
		if err != nil {
			t.Error(err)
		}
		if item.Data().ID == 0 {
			t.Error("can't get item properly")
		}
	})

	t.Run("Update", func(t *testing.T) {
		metadata := make(map[string]interface{})
		metadata["__metadata"] = map[string]string{"type": entType}
		metadata["Title"] = "Updated Item 1"
		body, _ := json.Marshal(metadata)
		if _, err := list.Items().GetByID(1).Update(body); err != nil {
			t.Error(err)
		}
	})

	t.Run("UpdateWithoutMetadataType", func(t *testing.T) {
		metadata := make(map[string]interface{})
		metadata["Title"] = "Updated Item 2"
		body, _ := json.Marshal(metadata)
		if _, err := list.Items().GetByID(2).Update(body); err != nil {
			t.Error(err)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		if _, err := list.Items().GetByID(1).Delete(); err != nil {
			t.Error(err)
		}
	})

	t.Run("Recycle", func(t *testing.T) {
		if _, err := list.Items().GetByID(2).Recycle(); err != nil {
			t.Error(err)
		}
	})

	if _, err := list.Delete(); err != nil {
		t.Error(err)
	}
}
