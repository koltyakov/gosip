package api

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestFilesChunked(t *testing.T) {
	checkClient(t)

	if envCode == "2013" {
		t.Skip("is not supported with SP 2013")
	}

	web := NewSP(spClient).Web()
	newFolderName := uuid.New().String()
	rootFolderURI := getRelativeURL(spClient.AuthCnfg.GetSiteURL()) + "/Shared%20Documents"
	newFolderURI := rootFolderURI + "/" + newFolderName
	if _, err := web.GetFolder(rootFolderURI).Folders().Add(newFolderName); err != nil {
		t.Error(err)
	}

	t.Run("AddChunkedMicro", func(t *testing.T) {
		fileName := fmt.Sprintf("TinyFile.txt")
		stream := strings.NewReader("Less than a chunk content")
		if _, err := web.GetFolder(newFolderURI).Files().AddChunked(fileName, stream, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("AddChunked", func(t *testing.T) {
		fileName := fmt.Sprintf("ChunkedFile.txt")
		content := "Greater than a chunk content..."
		stream := strings.NewReader(content)
		options := &AddChunkedOptions{
			Owerwrite: true,
			ChunkSize: 5,
		}
		fileResp, err := web.GetFolder(newFolderURI).Files().AddChunked(fileName, stream, options)
		if err != nil {
			t.Error(err)
		}
		data, err := web.GetFile(fileResp.Data().ServerRelativeURL).Download()
		if err != nil {
			t.Error(err)
		}
		if bytes.Compare([]byte(content), data) != 0 {
			t.Error("wrong file content after chunked upload")
		}
	})

	t.Run("AddChunkedCancel", func(t *testing.T) {
		fileName := fmt.Sprintf("ChunkedFile.txt")
		content := "Greater than a chunk content..."
		stream := strings.NewReader(content)
		options := &AddChunkedOptions{
			Owerwrite: true,
			ChunkSize: 5,
			Progress: func(data *FileUploadProgressData) bool {
				if data.BlockNumber > 0 {
					return false // cancel upload
				}
				return true
			},
		}
		_, err := web.GetFolder(newFolderURI).Files().AddChunked(fileName, stream, options)
		if err == nil {
			t.Error("cancel upload was not handled")
		}
		if err.Error() != "file upload was canceled" {
			t.Error(err)
		}
	})

	if err := web.GetFolder(newFolderURI).Delete(); err != nil {
		t.Error(err)
	}
}
