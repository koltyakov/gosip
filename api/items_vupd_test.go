package api

import (
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestItemsVUpd(t *testing.T) {
	checkClient(t)

	web := NewSP(spClient).Web()
	newListTitle := strings.Replace(uuid.New().String(), "-", "", -1)
	if _, err := web.Lists().Add(newListTitle, nil); err != nil {
		t.Error(err)
	}
	list := web.Lists().GetByTitle(newListTitle)

	t.Run("UpdateValidate", func(t *testing.T) {
		i, err := list.Items().Add([]byte(`{"Title":"Item"}`))
		if err != nil {
			t.Error(err)
		}

		options := &ValidateUpdateOptions{NewDocumentUpdate: true, CheckInComment: "test"}
		data := map[string]string{"Title": "New item"}
		res, err := list.Items().GetByID(i.Data().ID).UpdateValidate(data, options)
		if err != nil {
			t.Error(err)
		}

		if res.Value("Title") != "New item" {
			t.Error("unexpected field value")
		}
	})

	t.Run("UpdateValidateWrongPayload", func(t *testing.T) {
		i, err := list.Items().Add([]byte(`{"Title":"Item"}`))
		if err != nil {
			t.Error(err)
		}

		options := &ValidateUpdateOptions{NewDocumentUpdate: true, CheckInComment: "test"}
		data := map[string]string{"Modified": "wrong"}
		if _, err := list.Items().GetByID(i.Data().ID).UpdateValidate(data, options); err == nil {
			t.Error("failed to detect wrong payload")
		}
	})

	if err := list.Delete(); err != nil {
		t.Error(err)
	}

}
