package api

import (
	"testing"

	"github.com/google/uuid"
)

func TestChanges(t *testing.T) {
	checkClient(t)

	sp := NewSP(spClient)
	listTitle := uuid.New().String()

	if _, err := sp.Web().Lists().Add(listTitle, nil); err != nil {
		t.Error(err)
	}
	list := sp.Web().Lists().GetByTitle(listTitle)

	t.Run("GetCurrentToken", func(t *testing.T) {
		token, err := sp.Web().Changes().Conf(headers.verbose).GetCurrentToken()
		if err != nil {
			t.Error(err)
		}
		if token == "" {
			t.Error("empty change token")
		}

		if envCode != "2013" {
			token, err := sp.Web().Changes().Conf(headers.minimalmetadata).GetCurrentToken()
			if err != nil {
				t.Error(err)
			}
			if token == "" {
				t.Error("empty change token")
			}
			token, err = sp.Web().Changes().Conf(headers.nometadata).GetCurrentToken()
			if err != nil {
				t.Error(err)
			}
			if token == "" {
				t.Error("empty change token")
			}
		}
	})

	t.Run("ListChanges", func(t *testing.T) {
		token, err := list.Changes().GetCurrentToken()
		if err != nil {
			t.Error(err)
		}
		if token == "" {
			t.Error("empty change token")
		}
		if _, err := list.Items().Add([]byte(`{"Title":"Another item"}`)); err != nil {
			t.Error(err)
		}
		changes, err := list.Changes().GetChanges(&ChangeQuery{
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

	t.Run("WebChanges", func(t *testing.T) {
		token, err := sp.Web().Changes().GetCurrentToken()
		if err != nil {
			t.Error(err)
		}
		if token == "" {
			t.Error("empty change token")
		}
		if _, err := list.Items().Add([]byte(`{"Title":"New item"}`)); err != nil {
			t.Error(err)
		}
		changes, err := sp.Web().Changes().GetChanges(&ChangeQuery{
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

	t.Run("SiteChanges", func(t *testing.T) {
		token, err := sp.Site().Changes().GetCurrentToken()
		if err != nil {
			t.Error(err)
		}
		if token == "" {
			t.Error("empty change token")
		}
		if _, err := list.Items().Add([]byte(`{"Title":"New item"}`)); err != nil {
			t.Error(err)
		}
		changes, err := sp.Site().Changes().GetChanges(&ChangeQuery{
			ChangeTokenStart: token,
			Site:             true,
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

	if err := list.Delete(); err != nil {
		t.Error(err)
	}

}
