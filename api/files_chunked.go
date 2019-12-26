package api

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"strconv"
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
	if options == nil {
		options = &AddChunkedOptions{
			Owerwrite: true,
			ChunkSize: 10485760,
		}
	}

	sp := NewHTTPClient(files.client)
	endpoint := fmt.Sprintf("%s/Add(overwrite=%t,url='%s')", files.endpoint, options.Owerwrite, name)
	if _, err := sp.Post(endpoint, nil, getConfHeaders(files.config)); err != nil {
		return nil, err
	}

	if options.Progress == nil {
		options.Progress = func(data *FileUploadProgressData) {}
	}
	if options.ChunkSize == 0 {
		options.ChunkSize = 10485760
	}

	uploadID := uuid.New().String()

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
		var offset int
		if progress.BlockNumber == 0 {
			options.Progress(progress)
			offset, err = files.startUpload(name, uploadID, chunk)
		} else {
			progress.Stage = "continue"
			options.Progress(progress)
			offset, err = files.continueUpload(name, uploadID, progress.FileOffset, chunk)
		}
		if err != nil {
			return nil, err
		}
		progress.FileOffset = offset
		progress.BlockNumber++
	}

	progress.Stage = "finishing"
	options.Progress(progress)
	return files.finishUpload(name, uploadID, progress.FileOffset, nil)
}

func (files *Files) startUpload(name string, uploadID string, chunk []byte) (int, error) {
	sp := NewHTTPClient(files.client)
	endpoint := fmt.Sprintf("%s/Files('%s')/StartUpload(uploadId=guid'%s')", files.endpoint, name, uploadID)
	data, err := sp.Post(endpoint, chunk, getConfHeaders(files.config))
	if err != nil {
		return 0, err
	}
	if res, err := strconv.Atoi(fmt.Sprintf("%s", data)); err == nil {
		return res, nil
	}
	res := &struct {
		StartUpload int `json:"StartUpload"`
	}{}
	if err := json.Unmarshal(data, &res); err != nil {
		return 0, err
	}
	return res.StartUpload, nil
}

func (files *Files) continueUpload(name string, uploadID string, fileOffset int, chunk []byte) (int, error) {
	sp := NewHTTPClient(files.client)
	endpoint := fmt.Sprintf("%s/Files('%s')/ContinueUpload(uploadId=guid'%s',fileOffset=%d)", files.endpoint, name, uploadID, fileOffset)
	data, err := sp.Post(endpoint, chunk, getConfHeaders(files.config))
	if err != nil {
		return 0, err
	}
	if res, err := strconv.Atoi(fmt.Sprintf("%s", data)); err == nil {
		return res, nil
	}
	res := &struct {
		ContinueUpload int `json:"ContinueUpload"`
	}{}
	if err := json.Unmarshal(data, &res); err != nil {
		return 0, err
	}
	return res.ContinueUpload, nil
}

func (files *Files) finishUpload(name string, uploadID string, fileOffset int, chunk []byte) (FileResp, error) {
	sp := NewHTTPClient(files.client)
	endpoint := fmt.Sprintf("%s/Files('%s')/FinishUpload(uploadId=guid'%s',fileOffset=%d)", files.endpoint, name, uploadID, fileOffset)
	return sp.Post(endpoint, chunk, getConfHeaders(files.config))
}
