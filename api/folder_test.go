package api

import (
	"testing"

	"github.com/google/uuid"
)

func TestFolder(t *testing.T) {
	checkClient(t)

	web := NewSP(spClient).Web()
	newFolderName := uuid.New().String()
	rootFolderURI := getRelativeURL(spClient.AuthCnfg.GetSiteURL()) + "/Shared%20Documents"

	t.Run("Add", func(t *testing.T) {
		if _, err := web.GetFolder(rootFolderURI).Folders().Add(newFolderName); err != nil {
			t.Error(err)
		}
	})

	t.Run("Get", func(t *testing.T) {
		if _, err := web.GetFolder(rootFolderURI + "/" + newFolderName).Get(); err != nil {
			t.Error(err)
		}
	})

	t.Run("Props", func(t *testing.T) {
		props, err := web.GetFolder(rootFolderURI + "/" + newFolderName).Props().Get()
		if err != nil {
			t.Error(err)
		}
		if len(props.Data()) == 0 {
			t.Error("can't get property bags")
		}
	})

	t.Run("GetItem", func(t *testing.T) {
		item, err := web.GetFolder(rootFolderURI + "/" + newFolderName).GetItem()
		if err != nil {
			t.Error(err)
		}
		data, err := item.Get()
		if err != nil {
			t.Error(err)
		}
		if data.Data().ID == 0 {
			t.Error("can't get folder's item")
		}
	})

	t.Run("Recycle", func(t *testing.T) {
		if err := web.GetFolder(rootFolderURI + "/" + newFolderName).Delete(); err != nil {
			t.Error(err)
		}
		// ToDo: Restore and delete
	})

}
