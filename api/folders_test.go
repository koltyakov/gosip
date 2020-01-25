package api

import (
	"bytes"
	"testing"

	"github.com/google/uuid"
)

func TestFolders(t *testing.T) {
	checkClient(t)

	web := NewSP(spClient).Web()
	newFolderName := uuid.New().String()
	rootFolderURI := getRelativeURL(spClient.AuthCnfg.GetSiteURL()) + "/Shared%20Documents"

	t.Run("Modifiers", func(t *testing.T) {
		f := web.GetFolder(rootFolderURI).Folders()
		mods := f.Select("*").Expand("*").Filter("*").Top(1).OrderBy("*", true).modifiers
		if mods == nil || len(mods.mods) != 5 {
			t.Error("can't add modifiers")
		}
	})

	t.Run("Add", func(t *testing.T) {
		if _, err := web.GetFolder(rootFolderURI).Folders().Add(newFolderName); err != nil {
			t.Error(err)
		}
	})

	t.Run("Update", func(t *testing.T) {
		fm := []byte(`{"Name":"Test"}`)
		if _, err := web.GetFolder(rootFolderURI + "/" + newFolderName).Update(fm); err != nil {
			t.Error(err)
		}
	})

	t.Run("Get", func(t *testing.T) {
		data, err := web.GetFolder(rootFolderURI).Folders().Select("Id").Conf(headers.verbose).Get()
		if err != nil {
			t.Error(err)
		}
		if len(data.Data()) == 0 {
			t.Error("can't get folders")
		}
		if bytes.Compare(data, data.Normalized()) == -1 {
			t.Error("response normalization error")
		}
	})

	t.Run("GetByName", func(t *testing.T) {
		if _, err := web.GetFolder(rootFolderURI).Folders().GetByName(newFolderName).Get(); err != nil {
			t.Error(err)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		if err := web.GetFolder(rootFolderURI + "/" + newFolderName).Delete(); err != nil {
			t.Error(err)
		}
	})

	// ToDo:
	// Recycle

}
