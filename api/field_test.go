package api

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

var fieldID = ""

func TestField(t *testing.T) {
	t.Parallel()
	checkClient(t)

	web := NewSP(spClient).Web()
	field, err := getRandomField()
	if err != nil {
		t.Error(err)
	}

	t.Run("Get", func(t *testing.T) {
		data, err := field.Select("Id").Get()
		if err != nil {
			t.Error(err)
		}
		if data.Data().ID == "" {
			t.Error("can't unmarshal data")
		}
		if bytes.Compare(data, data.Normalized()) == -1 {
			t.Error("response normalization error")
		}
	})

	t.Run("UpdateDelete", func(t *testing.T) {
		guid := uuid.New().String()
		fm := []byte(`{"__metadata":{"type":"SP.FieldText"},"Title":"test-temp-` + guid + `","FieldTypeKind":2,"MaxLength":255}`)
		d, err := web.Fields().Add(fm)
		if err != nil {
			t.Error(err)
		}
		if _, err := web.Fields().GetByID(d.Data().ID).Update([]byte(`{"Description":"Test"}`)); err != nil {
			t.Error(err)
		}
		if err := web.Fields().GetByID(d.Data().ID).Delete(); err != nil {
			t.Error(err)
		}
	})

	// t.Run("Recycle", func(t *testing.T) {
	// 	guid := uuid.New().String()
	// 	fm := []byte(`{"__metadata":{"type":"SP.FieldText"},"Title":"test-temp-` + guid + `","FieldTypeKind":2,"MaxLength":255}`)
	// 	d, err := web.Fields().Add(fm)
	// 	if err != nil {
	// 		t.Error(err)
	// 	}
	// 	if err := web.Fields().GetByID(d.Data().ID).Recycle(); err != nil {
	// 		t.Error(err)
	// 	}
	// 	// ToDo: Empty Recycle Bin
	// })

}

func getRandomField() (*Field, error) {
	sp := NewSP(spClient)
	if fieldID == "" {
		resp, err := sp.Web().Fields().Top(1).Get()
		if err != nil {
			return nil, err
		}
		cts := resp.Data()
		if len(cts) != 1 {
			return nil, fmt.Errorf("wrong number of fields")
		}
		if cts[0].Data().ID == "" {
			return nil, fmt.Errorf("can't get field info")
		}
		fieldID = cts[0].Data().ID
	}
	return sp.Web().Fields().GetByID(fieldID), nil
}
