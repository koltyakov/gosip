package api

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestChanges(t *testing.T) {
	checkClient(t)

	sp := NewSP(spClient)
	listTitle := uuid.New().String()

	if _, err := sp.Web().Lists().Add(context.Background(), listTitle, nil); err != nil {
		t.Error(err)
	}
	list := sp.Web().Lists().GetByTitle(listTitle)

	t.Run("GetCurrentToken", func(t *testing.T) {
		token, err := sp.Web().Changes().Conf(headers.verbose).GetCurrentToken(context.Background())
		if err != nil {
			t.Error(err)
		}
		if token == "" {
			t.Error("empty change token")
		}

		if envCode != "2013" {
			token, err := sp.Web().Changes().Conf(headers.minimalmetadata).GetCurrentToken(context.Background())
			if err != nil {
				t.Error(err)
			}
			if token == "" {
				t.Error("empty change token")
			}
			token, err = sp.Web().Changes().Conf(headers.nometadata).GetCurrentToken(context.Background())
			if err != nil {
				t.Error(err)
			}
			if token == "" {
				t.Error("empty change token")
			}
		}
	})

	t.Run("ListChanges", func(t *testing.T) {
		token, err := list.Changes().GetCurrentToken(context.Background())
		if err != nil {
			t.Error(err)
		}
		if token == "" {
			t.Error("empty change token")
		}
		if _, err := list.Items().Add(context.Background(), []byte(`{"Title":"Another item"}`)); err != nil {
			t.Error(err)
		}
		changes, err := list.Changes().GetChanges(context.Background(), &ChangeQuery{
			ChangeTokenStart: token,
			List:             true,
			Item:             true,
			Add:              true,
		})
		if err != nil {
			t.Error(err)
		}
		if len(changes.Data()) == 0 {
			t.Error("incorrect changes data")
		}
	})

	t.Run("WebChanges", func(t *testing.T) {
		token, err := sp.Web().Changes().GetCurrentToken(context.Background())
		if err != nil {
			t.Error(err)
		}
		if token == "" {
			t.Error("empty change token")
		}
		if _, err := list.Items().Add(context.Background(), []byte(`{"Title":"New item"}`)); err != nil {
			t.Error(err)
		}
		changes, err := sp.Web().Changes().GetChanges(context.Background(), &ChangeQuery{
			ChangeTokenStart: token,
			Web:              true,
			Item:             true,
			Add:              true,
		})
		if err != nil {
			t.Error(err)
		}
		if len(changes.Data()) == 0 {
			t.Error("incorrect changes data")
		}
	})

	t.Run("SiteChanges", func(t *testing.T) {
		token, err := sp.Site().Changes().GetCurrentToken(context.Background())
		if err != nil {
			t.Error(err)
		}
		if token == "" {
			t.Error("empty change token")
		}
		if _, err := list.Items().Add(context.Background(), []byte(`{"Title":"New item"}`)); err != nil {
			t.Error(err)
		}
		changes, err := sp.Site().Changes().GetChanges(context.Background(), &ChangeQuery{
			ChangeTokenStart: token,
			Site:             true,
			Item:             true,
			Add:              true,
		})
		if err != nil {
			t.Error(err)
		}
		if len(changes.Data()) == 0 {
			t.Error("incorrect changes data")
		}
	})

	t.Run("GetChangeType#Add", func(t *testing.T) {
		changeName := list.Changes().GetChangeType(1)
		if changeName != "Add" {
			t.Error("incorrect change type")
		}
	})

	t.Run("GetChangeType#Unknown", func(t *testing.T) {
		changeName := list.Changes().GetChangeType(41)
		if changeName != "Unknown" {
			t.Error("incorrect change type")
		}
	})

	if err := list.Delete(context.Background()); err != nil {
		t.Error(err)
	}
}

func TestChangesPagination(t *testing.T) {
	checkClient(t)

	sp := NewSP(spClient)
	listTitle := uuid.New().String()

	if _, err := sp.Web().Lists().Add(context.Background(), listTitle, nil); err != nil {
		t.Error(err)
	}
	list := sp.Web().Lists().GetByTitle(listTitle)

	t.Run("ListChanges", func(t *testing.T) {
		token, err := list.Changes().GetCurrentToken(context.Background())
		if err != nil {
			t.Error(err)
		}
		if token == "" {
			t.Error("empty change token")
		}
		for i := 1; i <= 5; i++ {
			if _, err := list.Items().Add(context.Background(), []byte(fmt.Sprintf(`{"Title":"Item %d"}`, i))); err != nil {
				t.Error(err)
			}
		}
		changesFirstPage, err := list.Changes().Top(1).GetChanges(context.Background(), &ChangeQuery{
			ChangeTokenStart: token,
			List:             true,
			Item:             true,
			Add:              true,
		})
		if err != nil {
			t.Error(err)
		}
		if len(changesFirstPage.Data()) == 0 {
			t.Error("incorrect changes data")
		}

		changesSecondPage, err := changesFirstPage.GetNextPage()
		if err != nil {
			t.Error(err)
		}
		if len(changesSecondPage.Data()) == 0 {
			t.Error("incorrect changes data")
		}

		if changesFirstPage.Data()[0].ChangeToken.StringValue == changesSecondPage.Data()[0].ChangeToken.StringValue {
			t.Error("should be different change tokens")
		}
	})

	if err := list.Delete(context.Background()); err != nil {
		t.Error(err)
	}
}
