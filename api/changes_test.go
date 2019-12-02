package api

import (
	"testing"

	"github.com/google/uuid"
)

func TestChanges(t *testing.T) {
	checkClient(t)

	web := NewSP(spClient).Web()
	listTitle := uuid.New().String()

	if _, err := web.Lists().Add(listTitle, nil); err != nil {
		t.Error(err)
	}
	list := web.Lists().GetByTitle(listTitle)

	t.Run("WebChanges", func(t *testing.T) {
		token, err := web.GetChangeToken()
		if err != nil {
			t.Error(err)
		}
		if token == "" {
			t.Error("empty change token")
		}
		if _, err := list.Items().Add([]byte(`{"Title":"New item"}`)); err != nil {
			t.Error(err)
		}
		changes, err := web.GetChanges(&ChangeQuery{
			ChangeTokenStart: token,
			Web:              true,
			Item:             true,
			Add:              true,
		})
		if err != nil {
			t.Error(err)
		}
		if len(changes) == 0 {
			t.Error("incorrect changes data")
		}
	})

	t.Run("ListChanges", func(t *testing.T) {
		token, err := list.GetChangeToken()
		if err != nil {
			t.Error(err)
		}
		if token == "" {
			t.Error("empty change token")
		}
		if _, err := list.Items().Add([]byte(`{"Title":"Another item"}`)); err != nil {
			t.Error(err)
		}
		changes, err := list.GetChanges(&ChangeQuery{
			ChangeTokenStart: token,
			List:             true,
			Item:             true,
			Add:              true,
		})
		if err != nil {
			t.Error(err)
		}
		if len(changes) == 0 {
			t.Error("incorrect changes data")
		}
	})

	if _, err := list.Delete(); err != nil {
		t.Error(err)
	}

}
