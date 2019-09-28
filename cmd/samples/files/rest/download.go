package rest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/koltyakov/gosip"
)

// DownloadFile sample
func DownloadFile(client *gosip.SPClient, fileRelativeURL string) ([]byte, error) {

	endpoint := fmt.Sprintf(
		"%s/_api/Web/GetFileByServerRelativeUrl(@FileServerRelativeUrl)/$value?@FileServerRelativeUrl='%s'",
		client.AuthCnfg.GetSiteURL(),
		url.QueryEscape(fileRelativeURL),
	)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.TransferEncoding = []string{"null"}

	resp, err := client.Execute(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
