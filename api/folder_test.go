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

	t.Run("GetItem", func(t *testing.T) {
		if _, err := web.GetFolder(rootFolderURI + "/" + newFolderName).GetItem(); err != nil {
			t.Error(err)
		}
	})

	t.Run("Recycle", func(t *testing.T) {
		if _, err := web.GetFolder(rootFolderURI + "/" + newFolderName).Delete(); err != nil {
			t.Error(err)
		}
		// ToDo: Restore and delete
	})

}
