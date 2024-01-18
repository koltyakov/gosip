package api

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
)

var ctID = ""

func TestContentType(t *testing.T) {
	checkClient(t)

	web := NewSP(spClient).Web()

	t.Run("Modifiers", func(t *testing.T) {
		ct, err := getRandomCT()
		if err != nil {
			t.Error(err)
		}
		if _, err := ct.Select("*,Fields/*").Expand("Fields").Get(context.Background()); err != nil {
			t.Error(err)
		}
	})

	t.Run("Get", func(t *testing.T) {
		resp, err := web.ContentTypes().Top(5).Get(context.Background())
		if err != nil {
			t.Error(err)
		}
		cts := resp.Data()
		if len(cts) != 5 {
			t.Error("wrong number of content types")
		}
		if cts[0].Data().ID == "" {
			t.Error("can't get content type info")
		}
		data, err := web.ContentTypes().GetByID(cts[0].Data().ID).Get(context.Background())
		if err != nil {
			t.Error(err)
		}
		if bytes.Compare(data, data.Normalized()) == -1 {
			t.Error("response normalization error")
		}
	})

	t.Run("UpdateDelete", func(t *testing.T) {
		guid := uuid.New().String()
		ctID := "0x0100" + strings.ToUpper(strings.Replace(guid, "-", "", -1))
		ct := []byte(TrimMultiline(`{
			"Group":"Custom Content Types",
			"Id":{"StringValue":"` + ctID + `"},
			"Name":"test-temp-ct ` + guid + `"
		}`))
		ctResp, err := web.ContentTypes().Add(context.Background(), ct)
		if err != nil {
			t.Error(err)
		}
		ctID = ctResp.Data().ID // content type ID can't be set in REST API https://github.com/pnp/pnpjs/issues/457
		if _, err := web.ContentTypes().GetByID(ctID).Update(context.Background(), []byte(`{"Description":"Test"}`)); err != nil {
			t.Error(err)
		}
		if err := web.ContentTypes().GetByID(ctID).Delete(context.Background()); err != nil {
			t.Error(err)
		}
	})

}

func getRandomCT() (*ContentType, error) {
	sp := NewSP(spClient)
	if ctID == "" {
		resp, err := sp.Web().ContentTypes().Top(1).Get(context.Background())
		if err != nil {
			return nil, err
		}
		cts := resp.Data()
		if len(cts) != 1 {
			return nil, fmt.Errorf("wrong number of content types")
		}
		if cts[0].Data().ID == "" {
			return nil, fmt.Errorf("can't get content type info")
		}
		ctID = cts[0].Data().ID
	}
	return sp.Web().ContentTypes().GetByID(ctID), nil
}
