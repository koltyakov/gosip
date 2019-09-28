package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/koltyakov/gosip"
	strategy "github.com/koltyakov/gosip/auth/saml"
	"github.com/koltyakov/gosip/cmd/samples/files/rest"
)

var (
	fileRelativeURL = flag.String("fileRelativeUrl", "", "File relative URL")
	uploadFile      = flag.String("uploadFile", "", "Upload file path (local)")
	downloadTo      = flag.String("downloadTo", "", "Download file location")
)

func main() {
	flag.Parse()

	if *fileRelativeURL == "" {
		fmt.Println("fileRelativeURL can't be blank")
		return
	}

	configPath := "./config/private.saml.json"
	auth := &strategy.AuthCnfg{}

	err := auth.ReadConfig(configPath)
	if err != nil {
		fmt.Printf("unable to get config: %v\n", err)
		return
	}

	client := &gosip.SPClient{
		AuthCnfg: auth,
	}

	// Upload sample

	if *uploadFile != "" {

		file, err := os.Open(*uploadFile)
		if err != nil {
			fmt.Printf("unable to read a file: %v\n", err)
			return
		}
		defer file.Close()

		contentData, _ := ioutil.ReadAll(file)

		res, err := rest.UploadFile(
			client,
			filepath.Dir(*fileRelativeURL),
			filepath.Base(*fileRelativeURL),
			contentData,
		)
		if err != nil {
			fmt.Printf("unable to upload a file: %v\n", err)
			return
		}

		fmt.Printf("upload response: %s\n", res)

	}

	// Download sample

	if *uploadFile == "" {

		if *downloadTo == "" {
			*downloadTo = fmt.Sprintf("download/%s", filepath.Base(*fileRelativeURL))
		}

		data, err := rest.DownloadFile(client, *fileRelativeURL)
		if err != nil {
			fmt.Printf("unable to download a file: %v\n", err)
			return
		}

		_ = os.MkdirAll(filepath.Dir(*downloadTo), os.ModePerm)
		file, err := os.Create(*downloadTo)
		if err != nil {
			fmt.Printf("unable to create a file: %v\n", err)
			return
		}
		defer file.Close()

		_, err = file.Write(data)
		if err != nil {
			fmt.Printf("unable to write to file: %v\n", err)
			return
		}

		file.Sync()

	}

}
