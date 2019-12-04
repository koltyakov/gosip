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

// WebInfo ...
type WebInfo struct {
	ID                            string       `json:"Id"`
	Title                         string       `json:"Title"`
	AllowRssFeeds                 bool         `json:"AllowRssFeeds"`
	AlternateCSSURL               string       `json:"AlternateCssUrl"`
	AppInstanceID                 string       `json:"AppInstanceId"`
	ClassicWelcomePage            string       `json:"ClassicWelcomePage"`
	Configuration                 int          `json:"Configuration"`
	Created                       string       `json:"Created"` // time.Time
	CustomMasterURL               string       `json:"CustomMasterUrl"`
	Description                   string       `json:"Description"`
	DesignPackageID               string       `json:"DesignPackageId"`
	EnableMinimalDownload         bool         `json:"EnableMinimalDownload"`
	FooterEmphasis                int          `json:"FooterEmphasis"`
	FooterEnabled                 bool         `json:"FooterEnabled"`
	FooterLayout                  int          `json:"FooterLayout"`
	HeaderEmphasis                int          `json:"HeaderEmphasis"`
	HeaderLayout                  int          `json:"HeaderLayout"`
	HorizontalQuickLaunch         bool         `json:"HorizontalQuickLaunch"`
	IsHomepageModernized          bool         `json:"IsHomepageModernized"`
	IsMultilingual                bool         `json:"IsMultilingual"`
	IsRevertHomepageLinkHidden    bool         `json:"IsRevertHomepageLinkHidden"`
	Language                      int          `json:"Language"`
	LastItemModifiedDate          string       `json:"LastItemModifiedDate"`     // time.Time
	LastItemUserModifiedDate      string       `json:"LastItemUserModifiedDate"` // time.Time
	MasterURL                     string       `json:"MasterUrl"`
	MegaMenuEnabled               bool         `json:"MegaMenuEnabled"`
	NavAudienceTargetingEnabled   bool         `json:"NavAudienceTargetingEnabled"`
	NoCrawl                       bool         `json:"NoCrawl"`
	ObjectCacheEnabled            bool         `json:"ObjectCacheEnabled"`
	OverwriteTranslationsOnChange bool         `json:"OverwriteTranslationsOnChange"`
	QuickLaunchEnabled            bool         `json:"QuickLaunchEnabled"`
	RecycleBinEnabled             bool         `json:"RecycleBinEnabled"`
	SearchScope                   int          `json:"SearchScope"`
	ServerRelativeURL             string       `json:"ServerRelativeUrl"`
	SiteLogoURL                   string       `json:"SiteLogoUrl"`
	SyndicationEnabled            bool         `json:"SyndicationEnabled"`
	TenantAdminMembersCanShare    int          `json:"TenantAdminMembersCanShare"`
	TreeViewEnabled               bool         `json:"TreeViewEnabled"`
	UIVersion                     int          `json:"UIVersion"`
	UIVersionConfigurationEnabled bool         `json:"UIVersionConfigurationEnabled"`
	URL                           string       `json:"Url"`
	WebTemplate                   string       `json:"WebTemplate"`
	WelcomePage                   string       `json:"WelcomePage"`
	CurrentChangeToken            *StringValue `json:"CurrentChangeToken"`
	ResourcePath                  *DecodedURL  `json:"ResourcePath"`
}

// WebResp ...
type WebResp []byte

// NewWeb ...
func NewWeb(client *gosip.SPClient, endpoint string, config *RequestConfig) *Web {
	return &Web{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// ToURL ...
func (web *Web) ToURL() string {
	apiURL, _ := url.Parse(web.endpoint)
	query := url.Values{}
	for k, v := range web.modifiers {
		query.Add(k, trimMultiline(v))
	}
	apiURL.RawQuery = query.Encode()
	return apiURL.String()
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
func (web *Web) Get() (WebResp, error) {
	sp := NewHTTPClient(web.client)
	return sp.Get(web.ToURL(), getConfHeaders(web.config))
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
	return NewLists(
		web.client,
		fmt.Sprintf("%s/Lists", web.endpoint),
		web.config,
	)
}

// GetChangeToken ...
func (web *Web) GetChangeToken() (string, error) {
	scoped := *web
	data, err := scoped.Select("CurrentChangeToken").Get()
	if err != nil {
		return "", err
	}
	return data.Data().CurrentChangeToken.StringValue, nil
}

// GetChanges ...
func (web *Web) GetChanges(changeQuery *ChangeQuery) ([]*ChangeInfo, error) {
	return NewChanges(
		web.client,
		fmt.Sprintf("%s/GetChanges", web.endpoint),
		web.config,
	).GetChanges(changeQuery)
}

// Webs ...
func (web *Web) Webs() *Webs {
	return NewWebs(
		web.client,
		fmt.Sprintf("%s/Webs", web.endpoint),
		web.config,
	)
}

// Props ...
func (web *Web) Props() *WebProps {
	return NewWebProps(
		web.client,
		fmt.Sprintf("%s/AllProperties", web.endpoint),
		web.config,
	)
}

// ContentTypes ...
func (web *Web) ContentTypes() *ContentTypes {
	return NewContentTypes(
		web.client,
		fmt.Sprintf("%s/ContentTypes", web.endpoint),
		web.config,
	)
}

// Fields ...
func (web *Web) Fields() *Fields {
	return NewFields(
		web.client,
		fmt.Sprintf("%s/Fields", web.endpoint),
		web.config,
	)
}

// GetList ...
func (web *Web) GetList(listURI string) *List {
	// Prepend web relative URL to "Lists/ListPath" URIs
	if string([]rune(listURI)[0]) != "/" {
		absoluteURL := strings.Split(web.endpoint, "/_api")[0]
		relativeURL := getRelativeURL(absoluteURL)
		listURI = fmt.Sprintf("%s/%s", relativeURL, listURI)
	}
	return NewList(
		web.client,
		fmt.Sprintf("%s/GetList('%s')", web.endpoint, listURI),
		web.config,
	)
}

// EnsureUser ...
func (web *Web) EnsureUser(loginName string) (*UserInfo, error) {
	sp := NewHTTPClient(web.client)
	endpoint := fmt.Sprintf("%s/EnsureUser", web.endpoint)

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
	return NewGroups(
		web.client,
		fmt.Sprintf("%s/SiteGroups", web.endpoint),
		web.config,
	)
}

// SiteUsers ...
func (web *Web) SiteUsers() *Users {
	return NewUsers(
		web.client,
		fmt.Sprintf("%s/SiteUsers", web.endpoint),
		web.config,
	)
}

// CurrentUser ...
func (web *Web) CurrentUser() *User {
	return NewUser(
		web.client,
		fmt.Sprintf("%s/CurrentUser", web.endpoint),
		web.config,
	)
}

// GetFolder ...
func (web *Web) GetFolder(serverRelativeURL string) *Folder {
	return NewFolder(
		web.client,
		fmt.Sprintf(
			"%s/GetFolderByServerRelativeUrl('%s')",
			web.endpoint,
			serverRelativeURL,
		),
		web.config,
	)
}

// GetFile ...
func (web *Web) GetFile(serverRelativeURL string) *File {
	return NewFile(
		web.client,
		fmt.Sprintf(
			"%s/GetFileByServerRelativeUrl('%s')",
			web.endpoint,
			serverRelativeURL,
		),
		web.config,
	)
}

// Roles ...
func (web *Web) Roles() *Roles {
	return NewRoles(web.client, web.endpoint, web.config)
}

// RoleDefinitions ...
func (web *Web) RoleDefinitions() *RoleDefinitions {
	return NewRoleDefinitions(
		web.client,
		fmt.Sprintf("%s/RoleDefinitions", web.endpoint),
		web.config,
	)
}

/* Response helpers */

// Data : to get typed data
func (webResp *WebResp) Data() *WebInfo {
	data := parseODataItem(*webResp)
	res := &WebInfo{}
	json.Unmarshal(data, &res)
	return res
}

// Unmarshal : to unmarshal to custom object
func (webResp *WebResp) Unmarshal(obj interface{}) error {
	data := parseODataItem(*webResp)
	return json.Unmarshal(data, obj)
}
