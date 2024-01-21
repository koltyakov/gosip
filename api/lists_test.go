package api

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestLists(t *testing.T) {
	checkClient(t)

	web := NewSP(spClient).Web()
	newListTitle := uuid.New().String()

	t.Run("Get", func(t *testing.T) {
		data, err := web.Lists().Select("Id,Title").Conf(headers.verbose).Get(context.Background())
		if err != nil {
			t.Error(err)
		}
		if len(data.Data()) == 0 {
			t.Error("can't get webs")
		}
		if bytes.Compare(data, data.Normalized()) == -1 {
			t.Error("wrong response normalization")
		}
	})

	t.Run("Add", func(t *testing.T) {
		if _, err := web.Lists().Add(context.Background(), newListTitle, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		listInfo, err := getAnyList()
		if err != nil {
			t.Error(err)
		}
		if _, err := web.Lists().GetByID(listInfo.ID).Get(context.Background()); err != nil {
			t.Error(err)
		}
	})

	t.Run("GetByTitle", func(t *testing.T) {
		listInfo, err := getAnyList()
		if err != nil {
			t.Error(err)
		}
		if _, err := web.Lists().GetByTitle(listInfo.Title).Get(context.Background()); err != nil {
			t.Error(err)
		}
	})

	t.Run("GetListByURI", func(t *testing.T) {
		listURI := getRelativeURL(spClient.AuthCnfg.GetSiteURL()) +
			"/Lists/" + strings.Replace(newListTitle, "-", "", -1)
		if _, err := web.GetList(listURI).Get(context.Background()); err != nil {
			t.Error(err)
		}
	})

	t.Run("AddWithURI", func(t *testing.T) {
		if envCode == "2013" {
			t.Skip("is not supported with SP 2013")
		}

		listTitle := uuid.New().String()
		listURI := uuid.New().String()
		if _, err := web.Lists().AddWithURI(context.Background(), listTitle, listURI, nil); err != nil {
			t.Error(err)
		}
		if err := web.Lists().GetByTitle(listTitle).Delete(context.Background()); err != nil {
			t.Error(err)
		}
	})

	t.Run("CreateFieldAsXML", func(t *testing.T) {
		title := strings.Replace(uuid.New().String(), "-", "", -1)
		schemaXML := `<Field Type="Text" DisplayName="` + title + `" MaxLength="255" Name="` + title + `" Title="` + title + `"></Field>`
		if _, err := web.Lists().GetByTitle(newListTitle).Fields().CreateFieldAsXML(context.Background(), schemaXML, 12); err != nil {
			t.Error(err)
		}
		if err := web.Lists().GetByTitle(newListTitle).Fields().GetByInternalNameOrTitle(title).Delete(context.Background()); err != nil {
			t.Error(err)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		if err := web.Lists().GetByTitle(newListTitle).Delete(context.Background()); err != nil {
			t.Error(err)
		}
	})

}

func getAnyList() (*ListInfo, error) {
	web := NewSP(spClient).Web()
	data, err := web.Lists().Select("Id,Title").OrderBy("Created", true).Top(1).Conf(headers.verbose).Get(context.Background())
	if err != nil {
		return nil, err
	}
	if len(data.Data()) == 0 {
		return nil, fmt.Errorf("can't get webs")
	}
	return data.Data()[0].Data(), nil
}
