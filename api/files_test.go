package api

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestFiles(t *testing.T) {
	checkClient(t)

	web := NewSP(spClient).Web()
	newFolderName := uuid.New().String()
	rootFolderURI := getRelativeURL(spClient.AuthCnfg.GetSiteURL()) + "/Shared%20Documents"
	newFolderURI := rootFolderURI + "/" + newFolderName
	if _, err := web.GetFolder(rootFolderURI).Folders().Add(newFolderName); err != nil {
		t.Error(err)
	}

	// Add preset files
	for i := 1; i <= 5; i++ {
		fileName := fmt.Sprintf("File_%d.txt", i)
		fileData := []byte(fmt.Sprintf("File %d data", i))
		if _, err := web.GetFolder(newFolderURI).Files().Add(fileName, fileData, true); err != nil {
			t.Error(err)
		}
	}

	t.Run("Get", func(t *testing.T) {
		data, err := web.GetFolder(newFolderURI).Files().Get()
		if err != nil {
			t.Error(err)
		}
		if len(data.Data()) == 0 {
			t.Error("can't get files")
		}
		if bytes.Compare(data, data.Normalized()) == -1 {
			t.Error("response normalization error")
		}
	})

	t.Run("GetByName", func(t *testing.T) {
		data, err := web.GetFolder(newFolderURI).Files().GetByName("File_1.txt").Get()
		if err != nil {
			t.Error(err)
		}
		if data.Data().Length == 0 {
			t.Error("can't get file props")
		}
	})

	t.Run("GetFile", func(t *testing.T) {
		data, err := web.GetFile(newFolderURI + "/File_2.txt").Get()
		if err != nil {
			t.Error(err)
		}
		if data.Data().Name == "" {
			t.Error("can't get file props")
		}
	})

	t.Run("GetFileByPath", func(t *testing.T) {
		if envCode != "spo" {
			t.Skip("is not supported with legacy SP")
		}

		data, err := web.GetFileByPath(newFolderURI + "/File_2.txt").Get()
		if err != nil {
			t.Error(err)
		}
		if data.Data().Name == "" {
			t.Error("can't get file props")
		}
	})

	t.Run("GetFileByID", func(t *testing.T) {
		if envCode != "spo" {
			t.Skip("is not supported with legacy SP")
		}

		data, err := web.GetFile(newFolderURI + "/File_2.txt").Select("UniqueId").Get()
		if err != nil {
			t.Error(err)
		}
		if _, err := web.GetFileByID(data.Data().UniqueID).Get(); err != nil {
			t.Error(err)
		}
	})

	// Clean up
	if err := web.GetFolder(newFolderURI).Delete(); err != nil {
		t.Error(err)
	}
}
