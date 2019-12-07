package api

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"

	"github.com/google/uuid"
)

func TestItemsPaged(t *testing.T) {
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
		var wg sync.WaitGroup
		for i := 1; i <= 25; i++ {
			wg.Add(1)
			go func(i int) {
				metadata := make(map[string]interface{})
				metadata["__metadata"] = map[string]string{"type": entType}
				metadata["Title"] = fmt.Sprintf("Item %d", i)
				body, _ := json.Marshal(metadata)
				if _, err := list.Items().Add(body); err != nil {
					t.Error(err)
				}
				wg.Done()
			}(i)
		}
		wg.Wait()
	})

	t.Run("HasNextPage", func(t *testing.T) {
		items, err := list.Items().Select("Id").Top(1).Get()
		if err != nil {
			t.Error(err)
		}
		if items.NextPageURL() == "" {
			t.Error("can't get next page URL")
		}

		if envCode == "2013" {
			t.Skip("is not supported with SP 2013")
		}

		items, err = list.Items().Conf(HeadersPresets.Minimalmetadata).Select("Id").Top(1).Get()
		if err != nil {
			t.Error(err)
		}
		if items.NextPageURL() == "" {
			t.Error("can't get next page URL")
		}

		items, err = list.Items().Conf(HeadersPresets.Nometadata).Select("Id").Top(1).Get()
		if err != nil {
			t.Error(err)
		}
		if items.NextPageURL() == "" {
			t.Error("can't get next page URL")
		}
	})

	t.Run("GetAll", func(t *testing.T) {
		items, err := list.Items().Select("Id").Top(10).GetAll()
		if err != nil {
			t.Error(err)
		}
		for _, i := range items {
			fmt.Printf("%v\n", i.Data())
		}
		if len(items) != 25 {
			t.Errorf("incorrect items number, extected 25, got %d", len(items))
		}
	})

	if _, err := list.Delete(); err != nil {
		t.Error(err)
	}

}
