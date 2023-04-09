package api

import (
	"bytes"
	"fmt"
	"io"
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

	// Create a temporary document library with minor versioning enabled
	tempDocLibName := uuid.New().String()
	if _, err := web.Lists().Add(tempDocLibName, map[string]interface{}{
		"BaseTemplate":        101,
		"EnableMinorVersions": true,
	}); err != nil {
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

	t.Run("PubUnPub", func(t *testing.T) {
		// Using temporary document library with minor versioning enabled
		lib := web.Lists().GetByTitle(tempDocLibName)
		if _, err := lib.RootFolder().Files().Add("File_1.txt", []byte("File 1 data"), true); err != nil {
			t.Error(err)
		}
		if _, err := lib.RootFolder().Files().GetByName("File_1.txt").Publish("test"); err != nil {
			t.Error(err)
		}
		if _, err := lib.RootFolder().Files().GetByName("File_1.txt").UnPublish("test"); err != nil {
			t.Error(err)
		}
	})

	t.Run("Properties", func(t *testing.T) {
		if envCode == "2013" {
			t.Skip("is not supported with SP 2013")
		}

		file := web.GetFolder(newFolderURI).Files().GetByName("File_1.txt")
		props, err := file.Props().Get()
		if err != nil {
			t.Error(err)
		}
		if len(props.Data()) == 0 {
			t.Error("can't get file properties")
		}
		if err := file.Props().Set("MyProp", "MyValue"); err != nil {
			t.Error("can't set file property")
		}
	})

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

	t.Run("GetFileByID", func(t *testing.T) {
		data, err := web.GetFile(newFolderURI + "/File_3.txt").Get()
		if err != nil {
			t.Error(err)
		}
		d, err := web.GetFileByID(data.Data().UniqueID).Get()
		if err != nil {
			t.Error(err)
		}

		if data.Data().ServerRelativeURL != d.Data().ServerRelativeURL {
			t.Error("can't get file by ID")
		}
	})

	t.Run("GetFileByPath", func(t *testing.T) {
		data, err := web.GetFile(newFolderURI + "/File_3.txt").Get()
		if err != nil {
			t.Error(err)
		}
		d, err := web.GetFileByPath(data.Data().ServerRelativeURL).Get()
		if err != nil {
			t.Error(err)
		}

		if data.Data().ServerRelativeURL != d.Data().ServerRelativeURL {
			t.Error("can't get file by ID")
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

	t.Run("GetReader", func(t *testing.T) {
		fileContent := []byte("file content")
		if _, err := web.GetFolder(newFolderURI).Files().Add("reader.txt", fileContent, true); err != nil {
			t.Error(err)
		}
		reader, err := web.GetFile(newFolderURI + "/reader.txt").GetReader()
		if err != nil {
			t.Error(err)
		}
		defer reader.Close()
		content, err := io.ReadAll(reader)
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(fileContent, content) {
			t.Error("incorrect file body")
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
		if !bytes.Equal(fileContent, content) {
			t.Error("incorrect file body")
		}
	})

	t.Run("MoveTo", func(t *testing.T) {
		if _, err := web.GetFile(newFolderURI+"/File_4.txt").MoveTo(newFolderURI+"/File_4_moved.txt", true); err != nil {
			t.Error(err)
		}
	})

	t.Run("CopyTo", func(t *testing.T) {
		if _, err := web.GetFile(newFolderURI+"/File_5.txt").CopyTo(newFolderURI+"/File_5_copyed.txt", true); err != nil {
			t.Error(err)
		}
	})

	// Clean up
	if err := web.GetFolder(newFolderURI).Delete(); err != nil {
		t.Error(err)
	}
	if err := web.Lists().GetByTitle(tempDocLibName).Delete(); err != nil {
		t.Error(err)
	}
}
