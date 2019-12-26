package api

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestFilesChunked(t *testing.T) {
	checkClient(t)

	web := NewSP(spClient).Web()
	newFolderName := uuid.New().String()
	rootFolderURI := getRelativeURL(spClient.AuthCnfg.GetSiteURL()) + "/Shared%20Documents"
	newFolderURI := rootFolderURI + "/" + newFolderName
	if _, err := web.GetFolder(rootFolderURI).Folders().Add(newFolderName); err != nil {
		t.Error(err)
	}

	t.Run("AddChunked01", func(t *testing.T) {
		fileName := fmt.Sprintf("TinyFile.txt")
		stream := strings.NewReader(fmt.Sprintf("File %s data", fileName))
		if _, err := web.GetFolder(newFolderURI).Files().AddChunked(fileName, stream, nil); err != nil {
			t.Error(err)
		}
	})

	// if _, err := web.GetFolder(newFolderURI).Delete(); err != nil {
	// 	t.Error(err)
	// }
}
