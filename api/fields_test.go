package api

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestFields(t *testing.T) {
	checkClient(t)

	web := NewSP(spClient).Web()
	listURI := getRelativeURL(spClient.AuthCnfg.GetSiteURL()) + "/Shared%20Documents"
	field, err := getAnyField()
	if err != nil {
		t.Error(err)
	}

	t.Run("GetFromWeb", func(t *testing.T) {
		data, err := web.Fields().Top(1).Get(context.Background())
		if err != nil {
			t.Error(err)
		}
		if bytes.Compare(data, data.Normalized()) == -1 {
			t.Error("response normalization error")
		}
	})

	t.Run("GetFromList", func(t *testing.T) {
		if _, err := web.GetList(listURI).Fields().Top(1).Get(context.Background()); err != nil {
			t.Error(err)
		}
	})

	// t.Run("GetFromContentType", func(t *testing.T) {
	// 	contentType, err := getAnyContentType()
	// 	if err != nil {
	// 		t.Error(err)
	// 	}
	// 	if _, err := web.ContentTypes().GetByID(contentType.ID).Fields().Get(); err != nil {
	// 		t.Error(err)
	// 	}
	// })

	t.Run("GetByID", func(t *testing.T) {
		if _, err := web.Fields().GetByID(field.ID).Get(context.Background()); err != nil {
			t.Error(err)
		}
	})

	t.Run("GetByTitle", func(t *testing.T) {
		if _, err := web.Fields().GetByTitle(field.Title).Get(context.Background()); err != nil {
			t.Error(err)
		}
	})

	t.Run("GetByInternalNameOrTitle", func(t *testing.T) {
		if _, err := web.Fields().GetByInternalNameOrTitle(field.InternalName).Get(context.Background()); err != nil {
			t.Error(err)
		}
		if _, err := web.Fields().GetByInternalNameOrTitle(field.Title).Get(context.Background()); err != nil {
			t.Error(err)
		}
	})

	t.Run("Add", func(t *testing.T) {
		title := strings.Replace(uuid.New().String(), "-", "", -1)
		fm := []byte(`{"__metadata":{"type":"SP.FieldText"},"Title":"` + title + `","FieldTypeKind":2,"MaxLength":255}`)
		if _, err := web.Fields().Add(context.Background(), fm); err != nil {
			t.Error(err)
		}
		if err := web.Fields().GetByInternalNameOrTitle(title).Delete(context.Background()); err != nil {
			t.Error(err)
		}
	})

	// CreateFieldAsXML/Web - in web_test.go
	// CreateFieldAsXML/List - in lists_test.go

}

func getAnyField() (*FieldInfo, error) {
	web := NewSP(spClient).Web()
	data, err := web.Fields().Top(1).Get(context.Background())
	if err != nil {
		return nil, err
	}
	if len(data.Data()) == 0 {
		return nil, fmt.Errorf("can't get random field")
	}
	return data.Data()[0].Data(), nil
}
