package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/koltyakov/gosip"
)

// Web ...
type Web struct {
	client   *gosip.SPClient
	conf     *Conf
	endpoint string
	oSelect  string
	oExpand  string
}

// Conf ...
func (web *Web) Conf(conf *Conf) *Web {
	web.conf = conf
	return web
}

// Select ...
func (web *Web) Select(oDataSelect string) *Web {
	web.oSelect = oDataSelect
	return web
}

// Expand ...
func (web *Web) Expand(oDataExpand string) *Web {
	web.oExpand = oDataExpand
	return web
}

// Get ...
func (web *Web) Get() ([]byte, error) {
	apiURL, _ := url.Parse(web.endpoint)
	query := url.Values{}
	if web.oSelect != "" {
		query.Add("$select", TrimMultiline(web.oSelect))
	}
	if web.oExpand != "" {
		query.Add("$expand", TrimMultiline(web.oExpand))
	}
	apiURL.RawQuery = query.Encode()
	sp := &HTTPClient{SPClient: web.client}
	return sp.Get(apiURL.String(), GetConfHeaders(web.conf))
}

// Delete ...
func (web *Web) Delete() ([]byte, error) {
	sp := &HTTPClient{SPClient: web.client}
	return sp.Delete(web.endpoint, GetConfHeaders(web.conf))
}

// Update ...
func (web *Web) Update(body []byte) ([]byte, error) {
	sp := &HTTPClient{SPClient: web.client}
	return sp.Update(web.endpoint, body, GetConfHeaders(web.conf))
}

// Lists ...
func (web *Web) Lists() *Lists {
	return &Lists{
		client:   web.client,
		conf:     web.conf,
		endpoint: fmt.Sprintf("%s/lists", web.endpoint),
	}
}

// GetList ...
func (web *Web) GetList(listURI string) *List {
	// Prepend web relative URL to "Lists/ListPath" URIs
	if string([]rune(listURI)[0]) != "/" {
		absoluteURL := strings.Split(web.endpoint, "/_api")[0]
		relativeURL := GetRelativeURL(absoluteURL)
		listURI = fmt.Sprintf("%s/%s", relativeURL, listURI)
	}
	return &List{
		client: web.client,
		conf:   web.conf,
		endpoint: fmt.Sprintf(
			"%s/getList('%s')",
			web.endpoint,
			listURI,
		),
	}
}

// EnsureUser ...
func (web *Web) EnsureUser(loginName string) (*UserInfo, error) {
	sp := &HTTPClient{SPClient: web.client}
	endpoint := fmt.Sprintf("%s/ensureUser", web.endpoint)

	headers := GetConfHeaders(web.conf)
	headers["Accept"] = "application/json;odata=verbose"

	body := fmt.Sprintf(`{"logonName": "%s"}`, loginName)

	data, err := sp.Post(endpoint, []byte(body), headers)
	if err != nil {
		return nil, err
	}

	res := &struct {
		User *UserInfo `json:"d"`
	}{}

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, fmt.Errorf("unable to parse the response: %v", err)
	}

	return res.User, nil
}

// SiteGroups ...
func (web *Web) SiteGroups() *Groups {
	return &Groups{
		client: web.client,
		conf:   web.conf,
		endpoint: fmt.Sprintf(
			"%s/SiteGroups",
			web.endpoint,
		),
	}
}

// SiteUsers ...
func (web *Web) SiteUsers() *Users {
	return &Users{
		client: web.client,
		conf:   web.conf,
		endpoint: fmt.Sprintf(
			"%s/SiteUsers",
			web.endpoint,
		),
	}
}

// GetFolder ...
func (web *Web) GetFolder(serverRelativeURL string) *Folder {
	return &Folder{
		client: web.client,
		conf:   web.conf,
		endpoint: fmt.Sprintf(
			"%s/GetFolderByServerRelativeUrl('%s')",
			web.endpoint,
			serverRelativeURL,
		),
	}
}

// GetFile ...
func (web *Web) GetFile(serverRelativeURL string) *File {
	return &File{
		client: web.client,
		conf:   web.conf,
		endpoint: fmt.Sprintf(
			"%s/GetFileByServerRelativeUrl('%s')",
			web.endpoint,
			serverRelativeURL,
		),
	}
}

// Roles ...
func (web *Web) Roles() *Roles {
	return &Roles{
		client:   web.client,
		conf:     web.conf,
		endpoint: web.endpoint,
	}
}

// RoleDefinitions ...
func (web *Web) RoleDefinitions() *RoleDefinitions {
	return &RoleDefinitions{
		client: web.client,
		conf:   web.conf,
		endpoint: fmt.Sprintf(
			"%s/RoleDefinitions",
			web.endpoint,
		),
	}
}
