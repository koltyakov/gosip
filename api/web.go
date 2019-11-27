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
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers map[string]string
}

// Conf ...
func (web *Web) Conf(config *RequestConfig) *Web {
	web.config = config
	return web
}

// Select ...
func (web *Web) Select(oDataSelect string) *Web {
	if web.modifiers == nil {
		web.modifiers = make(map[string]string)
	}
	web.modifiers["$select"] = oDataSelect
	return web
}

// Expand ...
func (web *Web) Expand(oDataExpand string) *Web {
	if web.modifiers == nil {
		web.modifiers = make(map[string]string)
	}
	web.modifiers["$expand"] = oDataExpand
	return web
}

// Get ...
func (web *Web) Get() ([]byte, error) {
	apiURL, _ := url.Parse(web.endpoint)
	query := url.Values{}
	for k, v := range web.modifiers {
		query.Add(k, trimMultiline(v))
	}
	apiURL.RawQuery = query.Encode()
	sp := NewHTTPClient(web.client)
	return sp.Get(apiURL.String(), getConfHeaders(web.config))
}

// Delete ...
func (web *Web) Delete() ([]byte, error) {
	sp := NewHTTPClient(web.client)
	return sp.Delete(web.endpoint, getConfHeaders(web.config))
}

// Update ...
func (web *Web) Update(body []byte) ([]byte, error) {
	sp := NewHTTPClient(web.client)
	return sp.Update(web.endpoint, body, getConfHeaders(web.config))
}

// Lists ...
func (web *Web) Lists() *Lists {
	return &Lists{
		client:   web.client,
		config:   web.config,
		endpoint: fmt.Sprintf("%s/lists", web.endpoint),
	}
}

// Webs ...
func (web *Web) Webs() *Webs {
	return &Webs{
		client:   web.client,
		config:   web.config,
		endpoint: fmt.Sprintf("%s/webs", web.endpoint),
	}
}

// GetList ...
func (web *Web) GetList(listURI string) *List {
	// Prepend web relative URL to "Lists/ListPath" URIs
	if string([]rune(listURI)[0]) != "/" {
		absoluteURL := strings.Split(web.endpoint, "/_api")[0]
		relativeURL := getRelativeURL(absoluteURL)
		listURI = fmt.Sprintf("%s/%s", relativeURL, listURI)
	}
	return &List{
		client: web.client,
		config: web.config,
		endpoint: fmt.Sprintf(
			"%s/getList('%s')",
			web.endpoint,
			listURI,
		),
	}
}

// EnsureUser ...
func (web *Web) EnsureUser(loginName string) (*UserInfo, error) {
	sp := NewHTTPClient(web.client)
	endpoint := fmt.Sprintf("%s/ensureUser", web.endpoint)

	headers := getConfHeaders(web.config)
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
		config: web.config,
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
		config: web.config,
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
		config: web.config,
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
		config: web.config,
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
		config:   web.config,
		endpoint: web.endpoint,
	}
}

// RoleDefinitions ...
func (web *Web) RoleDefinitions() *RoleDefinitions {
	return &RoleDefinitions{
		client: web.client,
		config: web.config,
		endpoint: fmt.Sprintf(
			"%s/RoleDefinitions",
			web.endpoint,
		),
	}
}
