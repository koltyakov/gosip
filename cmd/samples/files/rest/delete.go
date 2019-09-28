package rest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/koltyakov/gosip"
)

// DeleteFile sample
func DeleteFile(client *gosip.SPClient, fileRelativeURL string) (string, error) {

	endpoint := fmt.Sprintf(
		"%s/_api/Web/GetFileByServerRelativeUrl(@FileServerRelativeUrl)/$value?@FileServerRelativeUrl='%s'",
		client.AuthCnfg.GetSiteURL(),
		url.QueryEscape(fileRelativeURL),
	)
	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("X-HTTP-Method", "DELETE")
	req.Header.Add("Accept", "application/json; odata=verbose")

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
