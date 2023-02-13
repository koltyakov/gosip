package api

import (
	"bytes"
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
		fileName := "TinyFile.txt"
		stream := strings.NewReader("Less than a chunk content")
		if _, err := web.GetFolder(newFolderURI).Files().AddChunked(fileName, stream, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("AddChunked", func(t *testing.T) {
		fileName := "ChunkedFile.txt"
		content := "Greater than a chunk content..."
		stream := strings.NewReader(content)
		options := &AddChunkedOptions{
			Overwrite: true,
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
		if !bytes.Equal([]byte(content), data) {
			t.Error("wrong file content after chunked upload")
		}
	})

	t.Run("AddChunkedNotEmtyOffset", func(t *testing.T) {
		fileName := "ChunkedFile1.txt"
		content := "Greater than a chunk content..."
		stream := strings.NewReader(content)
		for _, reqConfig := range []*RequestConfig{HeadersPresets.Minimalmetadata, HeadersPresets.Nometadata} {
			var offset int
			options := &AddChunkedOptions{
				Overwrite: true,
				ChunkSize: 5,
				Progress: func(progress *FileUploadProgressData) bool {
					if progress.BlockNumber == 0 {
						return true
					}
					offset = progress.FileOffset
					return false
				},
			}
			_, _ = web.GetFolder(newFolderURI).Files().
				Conf(reqConfig).
				AddChunked(fileName, stream, options)

			if offset == 0 {
				t.Error("wrong offset value")
			}
		}
	})

	t.Run("AddChunkedNilFinishPackage", func(t *testing.T) {
		fileName := "ChunkedFile.txt"
		content := "1234512345" // with combination of ChunkSize finishUpload package is nil
		stream := strings.NewReader(content)
		options := &AddChunkedOptions{
			Overwrite: true,
			ChunkSize: 5,
		}
		if _, err := web.GetFolder(newFolderURI).Files().AddChunked(fileName, stream, options); err != nil {
			t.Error(err)
		}
	})

	t.Run("AddChunkedZeroSize", func(t *testing.T) {
		fileName := "ChunkedFile.txt"
		content := "1234512345"
		stream := strings.NewReader(content)
		options := &AddChunkedOptions{
			Overwrite: true,
			ChunkSize: 0,
		}
		if _, err := web.GetFolder(newFolderURI).Files().AddChunked(fileName, stream, options); err != nil {
			t.Error(err)
		}
	})

	t.Run("AddChunkedCancel", func(t *testing.T) {
		fileName := "ChunkedFile.txt"
		content := "Greater than a chunk content..."
		stream := strings.NewReader(content)
		options := &AddChunkedOptions{
			Overwrite: true,
			ChunkSize: 5,
			Progress: func(data *FileUploadProgressData) bool {
				return data.BlockNumber <= 0 // cancel upload after first chunk
			},
		}
		_, err := web.GetFolder(newFolderURI).Files().AddChunked(fileName, stream, options)
		if err == nil {
			t.Error("cancel upload was not handled")
		}
		if err != nil && err.Error() != "file upload was canceled" {
			t.Error(err)
		}
	})

	t.Run("AddChunkedImmediateCancel", func(t *testing.T) {
		fileName := "ChunkedFile.txt"
		content := "Greater than a chunk content..."
		stream := strings.NewReader(content)
		options := &AddChunkedOptions{
			Overwrite: true,
			ChunkSize: 5,
			Progress: func(data *FileUploadProgressData) bool {
				return false // cancel upload immediately
			},
		}
		_, err := web.GetFolder(newFolderURI).Files().AddChunked(fileName, stream, options)
		if err == nil {
			t.Error("cancel upload was not handled")
		}
		if err != nil && err.Error() != "file upload was canceled" {
			t.Error(err)
		}
	})

	if err := web.GetFolder(newFolderURI).Delete(); err != nil {
		t.Error(err)
	}
}
