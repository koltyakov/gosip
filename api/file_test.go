package api

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestFile(t *testing.T) {
	checkClient(t)

	web := NewSP(spClient).Web()
	newFolderName := uuid.New().String()
	rootFolderURI := getRelativeURL(spClient.AuthCnfg.GetSiteURL()) + "/Shared%20Documents"
	newFolderURI := rootFolderURI + "/" + newFolderName
	if _, err := web.GetFolder(rootFolderURI).Folders().Add(newFolderName); err != nil {
		t.Error(err)
	}

	t.Run("AddSeries", func(t *testing.T) {
		for i := 1; i <= 5; i++ {
			fileName := fmt.Sprintf("File_%d.txt", i)
			fileData := []byte(fmt.Sprintf("File %d data", i))
			if _, err := web.GetFolder(newFolderURI).Files().Add(fileName, fileData, true); err != nil {
				t.Error(err)
			}
		}
	})

	t.Run("CheckOut", func(t *testing.T) {
		if _, err := web.GetFolder(newFolderURI).Files().GetByName("File_1.txt").CheckOut(); err != nil {
			t.Error(err)
		}
	})

	t.Run("UndoCheckOut", func(t *testing.T) {
		if _, err := web.GetFolder(newFolderURI).Files().GetByName("File_1.txt").UndoCheckOut(); err != nil {
			t.Error(err)
		}
	})

	t.Run("CheckIn", func(t *testing.T) {
		if _, err := web.GetFolder(newFolderURI).Files().GetByName("File_1.txt").CheckOut(); err != nil {
			t.Error(err)
		}
		if _, err := web.GetFolder(newFolderURI).Files().GetByName("File_1.txt").CheckIn("test", CheckInTypes.Major); err != nil {
			t.Error(err)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		if _, err := web.GetFolder(newFolderURI).Files().GetByName("File_1.txt").Delete(); err != nil {
			t.Error(err)
		}
	})

	t.Run("Recycle", func(t *testing.T) {
		if _, err := web.GetFile(newFolderURI + "/File_2.txt").Recycle(); err != nil {
			t.Error(err)
		}
	})

	t.Run("GetItem", func(t *testing.T) {
		if _, err := web.GetFile(newFolderURI + "/File_3.txt").Get(); err != nil {
			t.Error(err)
		}
	})

	if _, err := web.GetFolder(newFolderURI).Delete(); err != nil {
		t.Error(err)
	}
}
