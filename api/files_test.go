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

	t.Run("AddSeries", func(t *testing.T) {
		for i := 1; i <= 5; i++ {
			fileName := fmt.Sprintf("File_%d.txt", i)
			fileData := []byte(fmt.Sprintf("File %d data", i))
			if _, err := web.GetFolder(newFolderURI).Files().Add(fileName, fileData, true); err != nil {
				t.Error(err)
			}
		}
	})

	t.Run("Upload", func(t *testing.T) {
		fileName := "File_6.txt"
		fileReader := bytes.NewBuffer([]byte("File 6 data"))
		fileResp, err := web.GetFolder(newFolderURI).Files().Upload(fileName, fileReader, true)
		if err != nil {
			t.Error(err)
		}
		if err := web.GetFile(fileResp.Data().ServerRelativeURL).Recycle(); err != nil {
			t.Error(err)
		}
	})

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

	t.Run("Conf", func(t *testing.T) {
		files := web.GetFolder(newFolderURI).Files()
		hs := map[string]*RequestConfig{
			"nometadata":      HeadersPresets.Nometadata,
			"minimalmetadata": HeadersPresets.Minimalmetadata,
			"verbose":         HeadersPresets.Verbose,
		}
		for key, preset := range hs {
			f := files.Conf(preset)
			if f.config != preset {
				t.Errorf("can't %v config", key)
			}
		}
	})

	t.Run("Modifiers", func(t *testing.T) {
		files := web.GetFolder(newFolderURI).Files()
		mods := files.Select("*").Expand("*").Filter("*").Top(1).OrderBy("*", true).modifiers
		if mods == nil || len(mods.mods) != 5 {
			t.Error("wrong number of modifiers")
		}
	})

	if err := web.GetFolder(newFolderURI).Delete(); err != nil {
		t.Error(err)
	}
}
