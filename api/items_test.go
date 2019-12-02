package api

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestItems(t *testing.T) {
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

	t.Run("AddWithoutMetadataType", func(t *testing.T) {
		body := []byte(`{"Title":"Item"}`)
		if _, err := list.Items().Add(body); err != nil {
			t.Error(err)
		}
	})

	t.Run("AddResponse", func(t *testing.T) {
		body := []byte(`{"Title":"Item"}`)
		item, err := list.Items().Add(body)
		if err != nil {
			t.Error(err)
		}
		if item.Data().ID == 0 {
			t.Error("can't get item properly")
		}
	})

	t.Run("AddSeries", func(t *testing.T) {
		for i := 1; i < 10; i++ {
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
		items, err := list.Items().Top(100).OrderBy("Title", false).Get()
		if err != nil {
			t.Error(err)
		}
		if len(items.Data()) == 0 {
			t.Error("can't get items properly")
		}
		if items.Data()[0].Data().ID == 0 {
			t.Error("can't get items properly")
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		item, err := list.Items().GetByID(1).Get()
		if err != nil {
			t.Error(err)
		}
		if item.Data().ID == 0 {
			t.Error("can't get item properly")
		}
	})

	t.Run("GetByCAML", func(t *testing.T) {
		caml := `
			<View>
				<Query>
					<Where>
						<Eq>
							<FieldRef Name='ID' />
							<Value Type='Number'>3</Value>
						</Eq>
					</Where>
				</Query>
			</View>
		`
		data, err := list.Items().GetByCAML(caml)
		if err != nil {
			t.Error(err)
		}

		res := &struct {
			D struct {
				Results []interface{} `json:"results"`
			} `json:"d"`
		}{}

		if err := json.Unmarshal(data, &res); err != nil {
			t.Error(err)
		}

		if len(res.D.Results) != 1 {
			t.Error("incorrect number of items")
		}
	})

	if _, err := list.Delete(); err != nil {
		t.Error(err)
	}

}
