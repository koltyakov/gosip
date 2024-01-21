package api

import (
	"bytes"
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestFolder(t *testing.T) {
	checkClient(t)

	web := NewSP(spClient).Web()
	newFolderName := uuid.New().String()
	rootFolderURI := getRelativeURL(spClient.AuthCnfg.GetSiteURL()) + "/Shared%20Documents"

	t.Run("Add", func(t *testing.T) {
		if _, err := web.GetFolder(rootFolderURI).Folders().Add(context.Background(), newFolderName); err != nil {
			t.Error(err)
		}
	})

	t.Run("Get", func(t *testing.T) {
		data, err := web.GetFolder(rootFolderURI + "/" + newFolderName).Get(context.Background())
		if err != nil {
			t.Error(err)
		}
		if bytes.Compare(data, data.Normalized()) == -1 {
			t.Error("response normalization error")
		}
	})

	t.Run("GetFolderByID", func(t *testing.T) {
		data, err := web.GetFolder(rootFolderURI + "/" + newFolderName).Get(context.Background())
		if err != nil {
			t.Error(err)
		}
		d, err := web.GetFolderByID(data.Data().UniqueID).Get(context.Background())
		if err != nil {
			t.Error(err)
		}
		if d.Data().ServerRelativeURL != data.Data().ServerRelativeURL {
			t.Error("can't get folder by ID")
		}
	})

	t.Run("GetFolderByPath", func(t *testing.T) {
		data, err := web.GetFolder(rootFolderURI + "/" + newFolderName).Get(context.Background())
		if err != nil {
			t.Error(err)
		}
		d, err := web.GetFolderByPath(data.Data().ServerRelativeURL).Get(context.Background())
		if err != nil {
			t.Error(err)
		}
		if d.Data().ServerRelativeURL != data.Data().ServerRelativeURL {
			t.Error("can't get folder by ID")
		}
	})

	t.Run("ContextInfo", func(t *testing.T) {
		if _, err := web.GetFolder(rootFolderURI + "/" + newFolderName).ContextInfo(context.Background()); err != nil {
			t.Error(err)
		}
	})

	t.Run("ParentFolder", func(t *testing.T) {
		pf, err := web.GetFolder(rootFolderURI + "/" + newFolderName).ParentFolder().Get(context.Background())
		if err != nil {
			t.Error(err)
		}

		if pf.Data().Name == "" {
			t.Error("wrong parent folder name")
		}
	})

	t.Run("Props", func(t *testing.T) {
		folder := web.GetFolder(rootFolderURI + "/" + newFolderName)
		props, err := folder.Props().Get(context.Background())
		if err != nil {
			t.Error(err)
		}
		if len(props.Data()) == 0 {
			t.Error("can't get property bags")
		}
		if err := folder.Props().Set(context.Background(), "MyProp", "MyValue"); err != nil {
			t.Error("can't set property bags")
		}
	})

	t.Run("GetItem", func(t *testing.T) {
		item, err := web.GetFolder(rootFolderURI + "/" + newFolderName).GetItem(context.Background())
		if err != nil {
			t.Error(err)
		}
		data, err := item.Get(context.Background())
		if err != nil {
			t.Error(err)
		}
		if data.Data().ID == 0 {
			t.Error("can't get folder's item")
		}
	})

	t.Run("Delete", func(t *testing.T) {
		if err := web.GetFolder(rootFolderURI + "/" + newFolderName).Delete(context.Background()); err != nil {
			t.Error(err)
		}
	})

	t.Run("Recycle", func(t *testing.T) {
		guid := uuid.New().String()
		fr, err := web.GetFolder(rootFolderURI).Folders().Add(context.Background(), guid)
		if err != nil {
			t.Error(err)
		}
		if err := web.GetFolder(fr.Data().ServerRelativeURL).Recycle(context.Background()); err != nil {
			t.Error(err)
		}
	})

}
