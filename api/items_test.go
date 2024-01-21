package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestItems(t *testing.T) {
	checkClient(t)

	web := NewSP(spClient).Web()
	newListTitle := strings.Replace(uuid.New().String(), "-", "", -1)
	if _, err := web.Lists().Add(context.Background(), newListTitle, nil); err != nil {
		t.Error(err)
	}
	list := web.Lists().GetByTitle(newListTitle)
	entType, err := list.GetEntityType(context.Background())
	if err != nil {
		t.Error(err)
	}

	t.Run("AddWithoutMetadataType", func(t *testing.T) {
		body := []byte(`{"Title":"Item"}`)
		if _, err := list.Items().Add(context.Background(), body); err != nil {
			t.Error(err)
		}
	})

	t.Run("AddResponse", func(t *testing.T) {
		body := []byte(`{"Title":"Item"}`)
		item, err := list.Items().Add(context.Background(), body)
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
			if _, err := list.Items().Add(context.Background(), body); err != nil {
				t.Error(err)
			}
		}
	})

	t.Run("Get", func(t *testing.T) {
		items, err := list.Items().Top(100).OrderBy("Title", false).Get(context.Background())
		if err != nil {
			t.Error(err)
		}
		if len(items.Data()) == 0 {
			t.Error("can't get items properly")
		}
		if items.Data()[0].Data().ID == 0 {
			t.Error("can't get items properly")
		}
		if bytes.Compare(items, items.Normalized()) == -1 {
			t.Error("wrong response normalization")
		}
		if len(items.ToMap()) == 0 {
			t.Error("can't map items properly")
		}
		if items.ToMap()[0]["ID"] == 0 {
			t.Error("can't map items properly")
		}
	})

	t.Run("GetPaged", func(t *testing.T) {
		paged, err := list.Items().Top(5).GetPaged(context.Background())
		if err != nil {
			t.Error(err)
		}
		if len(paged.Items.Data()) == 0 {
			t.Error("can't get items")
		}
		if !paged.HasNextPage() {
			t.Error("can't get next page")
		} else {
			if _, err := paged.GetNextPage(); err != nil {
				t.Error(err)
			}
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		item, err := list.Items().GetByID(1).Get(context.Background())
		if err != nil {
			t.Error(err)
		}
		if item.Data().ID == 0 {
			t.Error("can't get item properly")
		}
	})

	t.Run("Get/Unmarshal", func(t *testing.T) {
		item, err := list.Items().GetByID(1).Get(context.Background())
		if err != nil {
			t.Error(err)
		}
		if item.Data().ID == 0 {
			t.Error("can't get item ID property properly")
		}
		if item.Data().Title == "" {
			t.Error("can't get item Title property properly")
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
		data, err := list.Items().Select("Id").GetByCAML(context.Background(), caml)
		if err != nil {
			t.Error(err)
		}
		if len(data.Data()) != 1 {
			t.Error("incorrect number of items")
		}
		if data.Data()[0].Data().ID != 3 {
			t.Error("incorrect response")
		}
	})

	t.Run("GetByCAMLAdvanced", func(t *testing.T) {
		viewResp, err := list.Views().DefaultView().Select("ListViewXml").Get(context.Background())
		if err != nil {
			t.Error(err)
		}
		if _, err := list.Items().GetByCAML(context.Background(), viewResp.Data().ListViewXML); err != nil {
			t.Error(err)
		}
	})

	if err := list.Delete(context.Background()); err != nil {
		t.Error(err)
	}

}
