package api

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/google/uuid"
)

// AddChunkedOptions ...
type AddChunkedOptions struct {
	Owerwrite bool
	Progress  func(data *FileUploadProgressData)
	ChunkSize int
}

// FileUploadProgressData ...
type FileUploadProgressData struct {
	UploadID    string
	Stage       string
	ChunkSize   int
	BlockNumber int
	FileOffset  int
}

// AddChunked ...
func (files *Files) AddChunked(name string, stream io.Reader, options *AddChunkedOptions) (FileResp, error) {
	web := NewSP(files.client).Web().Conf(files.config)
	var file *File
	uploadID := uuid.New().String()

	// Default props
	if options == nil {
		options = &AddChunkedOptions{
			Owerwrite: true,
			ChunkSize: 10485760,
		}
	}
	if options.Progress == nil {
		options.Progress = func(data *FileUploadProgressData) {}
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
			return files.Add(name, chunk, options.Owerwrite)
		}

		// Finishing uploading chunked file
		if size < options.ChunkSize && progress.BlockNumber > 0 {
			progress.Stage = "finishing"
			options.Progress(progress)
			return file.finishUpload(uploadID, progress.FileOffset, chunk)
		}

		// Initial chunked upload
		if progress.BlockNumber == 0 {
			options.Progress(progress)
			fileResp, err := files.Add(name, nil, options.Owerwrite)
			if err != nil {
				return nil, err
			}
			file = web.GetFile(fileResp.Data().ServerRelativeURL)
			offset, err := file.startUpload(uploadID, chunk)
			if err != nil {
				return nil, err
			}
			progress.FileOffset = offset
		} else { // or continue chunk upload
			progress.Stage = "continue"
			options.Progress(progress)
			offset, err := file.continueUpload(uploadID, progress.FileOffset, chunk)
			if err != nil {
				return nil, err
			}
			progress.FileOffset = offset
		}

		progress.BlockNumber++
	}

	progress.Stage = "finishing"
	options.Progress(progress)
	return file.finishUpload(uploadID, progress.FileOffset, nil)
}

func (file *File) startUpload(uploadID string, chunk []byte) (int, error) {
	sp := NewHTTPClient(file.client)
	endpoint := fmt.Sprintf("%s/StartUpload(uploadId=guid'%s')", file.endpoint, uploadID)
	data, err := sp.Post(endpoint, chunk, getConfHeaders(file.config))
	if err != nil {
		return 0, err
	}
	data = parseODataItem(data)
	if res, err := strconv.Atoi(fmt.Sprintf("%s", data)); err == nil {
		return res, nil
	}
	res := &struct {
		StartUpload int `json:"StartUpload,string"`
	}{}
	if err := json.Unmarshal(data, &res); err != nil {
		return 0, err
	}
	return res.StartUpload, nil
}

func (file *File) continueUpload(uploadID string, fileOffset int, chunk []byte) (int, error) {
	sp := NewHTTPClient(file.client)
	endpoint := fmt.Sprintf("%s/ContinueUpload(uploadId=guid'%s',fileOffset=%d)", file.endpoint, uploadID, fileOffset)
	data, err := sp.Post(endpoint, chunk, getConfHeaders(file.config))
	if err != nil {
		return 0, err
	}
	data = parseODataItem(data)
	if res, err := strconv.Atoi(fmt.Sprintf("%s", data)); err == nil {
		return res, nil
	}
	res := &struct {
		ContinueUpload int `json:"ContinueUpload,string"`
	}{}
	if err := json.Unmarshal(data, &res); err != nil {
		return 0, err
	}
	return res.ContinueUpload, nil
}

func (file *File) finishUpload(uploadID string, fileOffset int, chunk []byte) (FileResp, error) {
	sp := NewHTTPClient(file.client)
	endpoint := fmt.Sprintf("%s/FinishUpload(uploadId=guid'%s',fileOffset=%d)", file.endpoint, uploadID, fileOffset)
	return sp.Post(endpoint, chunk, getConfHeaders(file.config))
}
