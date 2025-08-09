package api

import (
	"bytes"
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
			t.Fatal(err)
		}
		if _, err := ct.Select("*,Fields/*").Expand("Fields").Get(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Get", func(t *testing.T) {
		resp, err := web.ContentTypes().Top(5).Get()
		if err != nil {
			t.Fatal(err)
		}
		cts := resp.Data()
		if len(cts) != 5 {
			t.Error("wrong number of content types")
		}
		if cts[0].Data().ID == "" {
			t.Error("can't get content type info")
		}
		data, err := web.ContentTypes().GetByID(cts[0].Data().ID).Get()
		if err != nil {
			t.Fatal(err)
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
		ctResp, err := web.ContentTypes().Add(ct)
		if err != nil {
			t.Fatal(err)
		}
		ctID = ctResp.Data().ID // content type ID can't be set in REST API https://github.com/pnp/pnpjs/issues/457
		if _, err := web.ContentTypes().GetByID(ctID).Update([]byte(`{"Description":"Test"}`)); err != nil {
			t.Fatal(err)
		}
		if err := web.ContentTypes().GetByID(ctID).Delete(); err != nil {
			t.Fatal(err)
		}
	})

}

func getRandomCT() (*ContentType, error) {
	sp := NewSP(spClient)
	if ctID == "" {
		resp, err := sp.Web().ContentTypes().Top(1).Get()
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
