package rest

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/koltyakov/gosip"
)

// UploadFile sample
func UploadFile(client *gosip.SPClient, folderRelativeURL string, fileName string, contentData []byte) (string, error) {

	endpoint := fmt.Sprintf(
		"%s/_api/web/getFolderByServerRelativeUrl('%s')/files/add(overwrite=true,url='%s')",
		client.AuthCnfg.GetSiteURL(),
		url.QueryEscape(folderRelativeURL),
		url.QueryEscape(fileName),
	)
	req, err := http.NewRequest(
		"POST",
		endpoint,
		bytes.NewBuffer(contentData),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/json;odata=verbose")

	resp, err := client.Execute(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", data), nil
}
