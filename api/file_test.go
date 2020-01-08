package api

import (
	"bytes"
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

	t.Run("Modifiers", func(t *testing.T) {
		f := web.GetFolder(newFolderURI).Files().GetByName("File_1.txt")
		mods := f.Select("*").Expand("*").modifiers
		if mods == nil || len(mods.mods) != 5 {
			t.Error("can't add modifiers")
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

	// t.Run("Update", func(t *testing.T) {
	// 	fm := []byte(`{"Name":"Test"}`)
	// 	if _, err := web.GetFolder(newFolderURI).Files().GetByName("File_1.txt").Update(fm); err != nil {
	// 		t.Error(err)
	// 	}
	// })

	t.Run("Delete", func(t *testing.T) {
		if err := web.GetFolder(newFolderURI).Files().GetByName("File_1.txt").Delete(); err != nil {
			t.Error(err)
		}
	})

	t.Run("Recycle", func(t *testing.T) {
		if err := web.GetFile(newFolderURI + "/File_2.txt").Recycle(); err != nil {
			t.Error(err)
		}
	})

	t.Run("Get", func(t *testing.T) {
		data, err := web.GetFile(newFolderURI + "/File_3.txt").Get()
		if err != nil {
			t.Error(err)
		}
		if bytes.Compare(data, data.Normalized()) == -1 {
			t.Error("response normalization error")
		}
	})

	t.Run("GetItem", func(t *testing.T) {
		item, err := web.GetFile(newFolderURI + "/File_3.txt").GetItem()
		if err != nil {
			t.Error(err)
		}
		data, err := item.Select("Id").Get()
		if err != nil {
			t.Error(err)
		}
		if data.Data().ID == 0 {
			t.Error("can't get file's item")
		}
	})

	t.Run("ContextInfo", func(t *testing.T) {
		if _, err := web.GetFile(newFolderURI + "/File_3.txt").ContextInfo(); err != nil {
			t.Error(err)
		}
	})

	t.Run("Download", func(t *testing.T) {
		fileContent := []byte("file content")
		if _, err := web.GetFolder(newFolderURI).Files().Add("download.txt", fileContent, true); err != nil {
			t.Error(err)
		}
		content, err := web.GetFile(newFolderURI + "/download.txt").Download()
		if err != nil {
			t.Error(err)
		}
		if bytes.Compare(fileContent, content) != 0 {
			t.Error("incorrect file body")
		}
	})

	t.Run("MoveTo", func(t *testing.T) {
		if envCode == "2013" {
			t.Skip("is not supported with SP 2013")
		}

		if _, err := web.GetFile(newFolderURI+"/File_4.txt").MoveTo(newFolderURI+"/File_4_moved.txt", true); err != nil {
			t.Error(err)
		}
	})

	t.Run("CopyTo", func(t *testing.T) {
		if envCode == "2013" {
			t.Skip("is not supported with SP 2013")
		}

		if _, err := web.GetFile(newFolderURI+"/File_5.txt").CopyTo(newFolderURI+"/File_5_copyed.txt", true); err != nil {
			t.Error(err)
		}
	})

	if err := web.GetFolder(newFolderURI).Delete(); err != nil {
		t.Error(err)
	}
}
