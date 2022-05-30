package api

import (
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestItemsVAdd(t *testing.T) {
	checkClient(t)

	web := NewSP(spClient).Web()
	newListTitle := strings.Replace(uuid.New().String(), "-", "", -1)
	if _, err := web.Lists().Add(newListTitle, nil); err != nil {
		t.Error(err)
	}
	list := web.Lists().GetByTitle(newListTitle)

	t.Run("AddValidate", func(t *testing.T) {
		if envCode == "2013" {
			t.Skip("is not supported with SP 2013")
		}

		options := &ValidateAddOptions{NewDocumentUpdate: true, CheckInComment: "test"}
		data := map[string]string{"Title": "New item"}
		res, err := list.Items().AddValidate(data, options)
		if err != nil {
			t.Error(err)
		}

		if res.Value("Title") != "New item" {
			t.Error("unexpected field value")
		}

		if res.ID() == 0 {
			t.Error("unexpected item ID")
		}
	})

	t.Run("AddValidateWrongPayload", func(t *testing.T) {
		if envCode == "2013" {
			t.Skip("is not supported with SP 2013")
		}

		options := &ValidateAddOptions{NewDocumentUpdate: true, CheckInComment: "test"}
		data := map[string]string{"Modified": "wrong"}
		if _, err := list.Items().AddValidate(data, options); err == nil {
			t.Error("failed to detect wrong payload")
		}
	})

	t.Run("AddValidateWithPath", func(t *testing.T) {
		if envCode == "2013" {
			t.Skip("is not supported with SP 2013")
		}
		// doesn't work anymore in SPO, an item can't be created in a folder which is not item-folder
		// if _, err := list.RootFolder().Folders().Add("subfolder"); err != nil {
		// 	t.Error(err)
		// }

		folderName := "subfolder"

		if _, err := list.Update([]byte(`{ "EnableFolderCreation": true }`)); err != nil {
			t.Error(err)
		}
		ff, err := list.Items().AddValidate(map[string]string{
			"Title":         folderName,
			"FileLeafRef":   folderName,
			"ContentType":   "Folder",
			"ContentTypeId": "0x0120",
		}, nil)
		if err != nil {
			t.Error(err)
		}
		if _, err := list.Items().GetByID(ff.ID()).Update([]byte(`{ "FileLeafRef": "` + folderName + `" }`)); err != nil {
			t.Error(err)
		}

		options := &ValidateAddOptions{NewDocumentUpdate: true, CheckInComment: "test"}
		options.DecodedPath = "Lists/" + newListTitle + "/" + folderName
		data := map[string]string{"Title": "New item in folder"}
		if _, err := list.Items().AddValidate(data, options); err != nil {
			t.Error(err)
		}
	})

	if err := list.Delete(); err != nil {
		t.Error(err)
	}

}
