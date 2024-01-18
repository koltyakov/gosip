package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/google/uuid"
)

// AddChunkedOptions provides optional settings for AddChunked method
type AddChunkedOptions struct {
	Overwrite bool                                    // should overwrite existing file
	Progress  func(data *FileUploadProgressData) bool // on progress callback, execute custom logic on each chunk, if the Progress is used it should return "true" to continue upload otherwise upload is canceled
	ChunkSize int                                     // chunk size in bytes
}

// FileUploadProgressData describes Progress callback options
type FileUploadProgressData struct {
	UploadID    string
	Stage       string
	ChunkSize   int
	BlockNumber int
	FileOffset  int
}

// AddChunked uploads a file in chunks (streaming), is a good fit for large files. Supported starting from SharePoint 2016.
func (files *Files) AddChunked(ctx context.Context, name string, stream io.Reader, options *AddChunkedOptions) (FileResp, error) {
	web := NewSP(files.client).Web().Conf(files.config)
	var file *File
	uploadID := uuid.New().String()

	cancelUpload := func(ctx context.Context, file *File, uploadID string) error {
		if err := file.cancelUpload(ctx, uploadID); err != nil {
			return err
		}
		return fmt.Errorf("file upload was canceled")
	}

	// Default props
	if options == nil {
		options = &AddChunkedOptions{
			Overwrite: true,
			ChunkSize: 10485760,
		}
	}
	if options.Progress == nil {
		options.Progress = func(data *FileUploadProgressData) bool {
			return true
		}
	}
	if options.ChunkSize == 0 {
		options.ChunkSize = 10485760
	}

	progress := &FileUploadProgressData{
		UploadID:    uploadID,
		Stage:       "starting",
		ChunkSize:   options.ChunkSize,
		BlockNumber: 0,
		FileOffset:  0,
	}

	slot := make([]byte, options.ChunkSize)
	for {
		size, err := stream.Read(slot)
		if err == io.EOF {
			break
		}
		chunk := slot[:size]

		// Upload in a call if file size is less than chunk size
		if size < options.ChunkSize && progress.BlockNumber == 0 {
			return files.Add(ctx, name, chunk, options.Overwrite)
		}

		// Finishing uploading chunked file
		if size < options.ChunkSize && progress.BlockNumber > 0 {
			progress.Stage = "finishing"
			if !options.Progress(progress) {
				return nil, cancelUpload(ctx, file, uploadID)
			}
			if file == nil {
				return nil, fmt.Errorf("can't get file object")
			}
			return file.finishUpload(ctx, uploadID, progress.FileOffset, chunk)
		}

		// Initial chunked upload
		if progress.BlockNumber == 0 {
			progress.Stage = "starting"
			if !options.Progress(progress) {
				return nil, fmt.Errorf("file upload was canceled") // cancelUpload(file, uploadID)
			}
			fileResp, err := files.Add(ctx, name, nil, options.Overwrite)
			if err != nil {
				return nil, err
			}
			file = web.GetFile(fileResp.Data().ServerRelativeURL)
			offset, err := file.startUpload(ctx, uploadID, chunk)
			if err != nil {
				return nil, err
			}
			progress.FileOffset = offset
		} else { // or continue chunk upload
			progress.Stage = "continue"
			if !options.Progress(progress) {
				return nil, cancelUpload(ctx, file, uploadID)
			}
			if file == nil {
				return nil, fmt.Errorf("can't get file object")
			}
			offset, err := file.continueUpload(ctx, uploadID, progress.FileOffset, chunk)
			if err != nil {
				return nil, err
			}
			progress.FileOffset = offset
		}

		progress.BlockNumber++
	}

	progress.Stage = "finishing"
	if !options.Progress(progress) {
		return nil, cancelUpload(ctx, file, uploadID)
	}
	if file == nil {
		return nil, fmt.Errorf("can't get file object")
	}
	return file.finishUpload(ctx, uploadID, progress.FileOffset, nil)
}

// startUpload starts uploading a document using chunk API
func (file *File) startUpload(ctx context.Context, uploadID string, chunk []byte) (int, error) {
	client := NewHTTPClient(file.client)
	endpoint := fmt.Sprintf("%s/StartUpload(uploadId=guid'%s')", file.endpoint, uploadID)
	data, err := client.Post(ctx, endpoint, bytes.NewBuffer(chunk), file.config)
	if err != nil {
		return 0, err
	}
	data = NormalizeODataItem(data)
	if res, err := strconv.Atoi(string(data)); err == nil {
		return res, nil
	}
	res := &struct {
		StartUpload int `json:"StartUpload,string"` // With OData=verbose
		Value       int `json:"value,string"`       // With OData=minimanmetadata or nometadata
	}{}
	if err := json.Unmarshal(data, &res); err != nil {
		return 0, err
	}
	offset := res.StartUpload
	if offset == 0 {
		offset = res.Value
	}
	return offset, nil
}

// continueUpload continues uploading a document using chunk API
func (file *File) continueUpload(ctx context.Context, uploadID string, fileOffset int, chunk []byte) (int, error) {
	client := NewHTTPClient(file.client)
	endpoint := fmt.Sprintf("%s/ContinueUpload(uploadId=guid'%s',fileOffset=%d)", file.endpoint, uploadID, fileOffset)
	data, err := client.Post(ctx, endpoint, bytes.NewBuffer(chunk), file.config)
	if err != nil {
		return 0, err
	}
	data = NormalizeODataItem(data)
	if res, err := strconv.Atoi(string(data)); err == nil {
		return res, nil
	}
	res := &struct {
		ContinueUpload int `json:"ContinueUpload,string"` // With OData=verbose
		Value          int `json:"value,string"`          // With OData=minimanmetadata or nometadata
	}{}
	if err := json.Unmarshal(data, &res); err != nil {
		return 0, err
	}
	offset := res.ContinueUpload
	if offset == 0 {
		offset = res.Value
	}
	return offset, nil
}

// cancelUpload cancels document upload using chunk API
func (file *File) cancelUpload(ctx context.Context, uploadID string) error {
	client := NewHTTPClient(file.client)
	endpoint := fmt.Sprintf("%s/CancelUpload(uploadId=guid'%s')", file.endpoint, uploadID)
	_, err := client.Post(ctx, endpoint, nil, file.config)
	return err
}

// finishUpload finishes uploading a document using chunk API
func (file *File) finishUpload(ctx context.Context, uploadID string, fileOffset int, chunk []byte) (FileResp, error) {
	client := NewHTTPClient(file.client)
	endpoint := fmt.Sprintf("%s/FinishUpload(uploadId=guid'%s',fileOffset=%d)", file.endpoint, uploadID, fileOffset)
	return client.Post(ctx, endpoint, bytes.NewBuffer(chunk), file.config)
}
