package api

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/koltyakov/gosip"
)

// Web represents SharePoint Web API queryable object struct
// Always use NewWeb constructor instead of &Web{}
type Web struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers *ODataMods
}

// WebInfo - web API response payload structure
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

// WebResp - web response type with helper processor methods
type WebResp []byte

// NewWeb - Web struct constructor function
func NewWeb(client *gosip.SPClient, endpoint string, config *RequestConfig) *Web {
	return &Web{
		client:    client,
		endpoint:  endpoint,
		config:    config,
		modifiers: NewODataMods(),
	}
}

// ToURL gets endpoint with modificators raw URL
func (web *Web) ToURL() string {
	return toURL(web.endpoint, web.modifiers)
}

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (web *Web) Conf(config *RequestConfig) *Web {
	web.config = config
	return web
}

// Select adds $select OData modifier
func (web *Web) Select(oDataSelect string) *Web {
	web.modifiers.AddSelect(oDataSelect)
	return web
}

// Expand adds $expand OData modifier
func (web *Web) Expand(oDataExpand string) *Web {
	web.modifiers.AddExpand(oDataExpand)
	return web
}

// Get gets this Web info
func (web *Web) Get() (WebResp, error) {
	sp := NewHTTPClient(web.client)
	return sp.Get(web.ToURL(), getConfHeaders(web.config))
}

// Delete deletes this Web
func (web *Web) Delete() ([]byte, error) {
	sp := NewHTTPClient(web.client)
	return sp.Delete(web.endpoint, getConfHeaders(web.config))
}

// Update updates Web's metadata with properties provided in `body` parameter
// where `body` is byte array representation of JSON string payload relevalt to SP.Web object
func (web *Web) Update(body []byte) ([]byte, error) {
	body = patchMetadataType(body, "SP.Web")
	sp := NewHTTPClient(web.client)
	return sp.Update(web.endpoint, body, getConfHeaders(web.config))
}

// Lists gets Lists API instance object
func (web *Web) Lists() *Lists {
	return NewLists(
		web.client,
		fmt.Sprintf("%s/Lists", web.endpoint),
		web.config,
	)
}

// GetChangeToken gets current change token for this Web
func (web *Web) GetChangeToken() (string, error) {
	scoped := NewWeb(web.client, web.endpoint, web.config)
	data, err := scoped.Select("CurrentChangeToken").Get()
	if err != nil {
		return "", err
	}
	return data.Data().CurrentChangeToken.StringValue, nil
}

// GetChanges gets changes on this Web due to the configuration provided as `changeQuery` parameter
func (web *Web) GetChanges(changeQuery *ChangeQuery) ([]*ChangeInfo, error) {
	return NewChanges(
		web.client,
		fmt.Sprintf("%s/GetChanges", web.endpoint),
		web.config,
	).GetChanges(changeQuery)
}

// Webs gets Webs API instance queryable collection for this Web (Subwebs)
func (web *Web) Webs() *Webs {
	return NewWebs(
		web.client,
		fmt.Sprintf("%s/Webs", web.endpoint),
		web.config,
	)
}

// Props gets WebProps API instance queryable collection for this Web
func (web *Web) Props() *WebProps {
	return NewWebProps(
		web.client,
		fmt.Sprintf("%s/AllProperties", web.endpoint),
		web.config,
	)
}

// ContentTypes gets ContentTypes API instance queryable collection for this Web
func (web *Web) ContentTypes() *ContentTypes {
	return NewContentTypes(
		web.client,
		fmt.Sprintf("%s/ContentTypes", web.endpoint),
		web.config,
	)
}

// Fields gets Fields API instance queryable collection for this Web (Site Columns)
func (web *Web) Fields() *Fields {
	return NewFields(
		web.client,
		fmt.Sprintf("%s/Fields", web.endpoint),
		web.config,
	)
}

// GetList gets Lists API instance queryable collection for this Web (Web's Lists)
func (web *Web) GetList(listURI string) *List {
	return NewList(
		web.client,
		fmt.Sprintf("%s/GetList('%s')", web.endpoint, checkGetRelativeURL(listURI, web.endpoint)),
		web.config,
	)
}

// EnsureUser ensures a user by a `loginName` parameter and returns UserInfo
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

// SiteGroups gets Groups API instance queryable collection for this Web (Site Groups)
func (web *Web) SiteGroups() *Groups {
	return NewGroups(
		web.client,
		fmt.Sprintf("%s/SiteGroups", web.endpoint),
		web.config,
	)
}

// SiteUsers gets Users API instance queryable collection for this Web (Site Users)
func (web *Web) SiteUsers() *Users {
	return NewUsers(
		web.client,
		fmt.Sprintf("%s/SiteUsers", web.endpoint),
		web.config,
	)
}

// AssociatedGroups gets associated groups scoped constructor
func (web *Web) AssociatedGroups() *AssociatedGroups {
	return NewAssociatedGroups(
		web.client,
		web.endpoint,
		web.config,
	)
}

// CurrentUser gets current User API instance object
func (web *Web) CurrentUser() *User {
	return NewUser(
		web.client,
		fmt.Sprintf("%s/CurrentUser", web.endpoint),
		web.config,
	)
}

// GetFolder gets a folder by its relevant URI, URI can be host relevant (e.g. `/sites/site/lib/folder`)
// or web relavant (e.g. `lib/folder`, with web relevant URI there should be no slash at the begining)
func (web *Web) GetFolder(serverRelativeURL string) *Folder {
	return NewFolder(
		web.client,
		fmt.Sprintf(
			"%s/GetFolderByServerRelativeUrl('%s')",
			web.endpoint,
			checkGetRelativeURL(serverRelativeURL, web.endpoint),
		),
		web.config,
	)
}

// EnsureFolder is a helper to ensure a folder by its relevalt URI, when there was no folder it's created
func (web *Web) EnsureFolder(serverRelativeURL string) ([]byte, error) {
	return ensureFolder(web, serverRelativeURL, serverRelativeURL)
}

// GetFile gets File API instance object by its relevant URI
// File URI can be host relevant (e.g. `/sites/site/lib/folder/file.txt`)
// or web relavant (e.g. `lib/folder/file.txt`, with web relevant URI there should be no slash at the begining)
func (web *Web) GetFile(serverRelativeURL string) *File {
	return NewFile(
		web.client,
		fmt.Sprintf(
			"%s/GetFileByServerRelativeUrl('%s')",
			web.endpoint,
			checkGetRelativeURL(serverRelativeURL, web.endpoint),
		),
		web.config,
	)
}

// Roles gets Roles API instance queryable collection for this Web
func (web *Web) Roles() *Roles {
	return NewRoles(web.client, web.endpoint, web.config)
}

// RoleDefinitions gets RoleDefinitions API instance queryable collection for this Web
func (web *Web) RoleDefinitions() *RoleDefinitions {
	return NewRoleDefinitions(
		web.client,
		fmt.Sprintf("%s/RoleDefinitions", web.endpoint),
		web.config,
	)
}

// Features gets Features API instance queryable collection for this Web
func (web *Web) Features() *Features {
	endpoint := fmt.Sprintf("%s/Features", web.endpoint)
	return NewFeatures(web.client, endpoint, web.config)
}

// RecycleBin gets RecycleBin API instance object for this Web
func (web *Web) RecycleBin() *RecycleBin {
	endpoint := fmt.Sprintf("%s/RecycleBin", web.endpoint)
	return NewRecycleBin(web.client, endpoint, web.config)
}

// ContextInfo gets Context info object for this Web
func (web *Web) ContextInfo() (*ContextInfo, error) {
	return NewContext(web.client, web.ToURL(), web.config).Get()
}

// ToDo:
// Custom actions

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

/* Miscellaneous */

func ensureFolder(web *Web, serverRelativeURL string, currentRelativeURL string) ([]byte, error) {
	data, err := web.GetFolder(currentRelativeURL).Get()
	if err != nil {
		splitted := strings.Split(currentRelativeURL, "/")
		if len(splitted) == 1 {
			return nil, err
		}
		splitted = splitted[0 : len(splitted)-1]
		currentRelativeURL = strings.Join(splitted, "/")
		return ensureFolder(web, serverRelativeURL, currentRelativeURL)
	}

	curFolders := strings.Split(currentRelativeURL, "/")
	expFolders := strings.Split(serverRelativeURL, "/")

	if len(curFolders) == len(expFolders) {
		return data, nil
	}

	createFolders := expFolders[len(curFolders):]
	for _, folder := range createFolders {
		data, err = web.GetFolder(currentRelativeURL).Folders().Add(folder)
		if err != nil {
			return nil, err
		}
		currentRelativeURL += "/" + folder
	}

	return data, nil
}
