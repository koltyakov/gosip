package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/koltyakov/gosip"
)

//go:generate ggen -ent Web -conf -mods Select,Expand -helpers Data,Normalized

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

// FromURL gets Web object using its API URL
func (web *Web) FromURL(url string) *Web {
	url = strings.Split(url, "?")[0]
	return NewWeb(web.client, url, web.config)
}

// Get gets this Web info
func (web *Web) Get() (WebResp, error) {
	client := NewHTTPClient(web.client)
	return client.Get(web.ToURL(), web.config)
}

// Delete deletes this Web
func (web *Web) Delete() error {
	client := NewHTTPClient(web.client)
	_, err := client.Delete(web.endpoint, web.config)
	return err
}

// Update updates Web's metadata with properties provided in `body` parameter
// where `body` is byte array representation of JSON string payload relevant to SP.Web object
func (web *Web) Update(body []byte) (WebResp, error) {
	body = patchMetadataType(body, "SP.Web")
	client := NewHTTPClient(web.client)
	return client.Update(web.endpoint, bytes.NewBuffer(body), web.config)
}

// Lists gets Lists API instance object
func (web *Web) Lists() *Lists {
	return NewLists(
		web.client,
		fmt.Sprintf("%s/Lists", web.endpoint),
		web.config,
	)
}

// Changes gets changes API scoped object
func (web *Web) Changes() *Changes {
	return NewChanges(
		web.client,
		web.endpoint,
		web.config,
	)
}

// Webs gets Webs API instance queryable collection for this Web (sub webs)
func (web *Web) Webs() *Webs {
	return NewWebs(
		web.client,
		fmt.Sprintf("%s/Webs", web.endpoint),
		web.config,
	)
}

// AllProps gets Properties API instance queryable collection for this Web
func (web *Web) AllProps() *Properties {
	return NewProperties(
		web.client,
		fmt.Sprintf("%s/AllProperties", web.endpoint),
		web.config,
		"web",
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

// AvailableContentTypes gets ContentTypes API instance queryable collection for this Web
func (web *Web) AvailableContentTypes() *ContentTypes {
	return NewContentTypes(
		web.client,
		fmt.Sprintf("%s/AvailableContentTypes", web.endpoint),
		web.config,
	)
}

// Fields gets Fields API instance queryable collection for this Web (Site Columns)
func (web *Web) Fields() *Fields {
	return NewFields(
		web.client,
		fmt.Sprintf("%s/Fields", web.endpoint),
		web.config,
		"web",
	)
}

// RootFolder gets web's root folder object
func (web *Web) RootFolder() *Folder {
	return NewFolder(
		web.client,
		fmt.Sprintf("%s/RootFolder", web.endpoint),
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
	client := NewHTTPClient(web.client)
	endpoint := fmt.Sprintf("%s/EnsureUser", web.endpoint)

	headers := getConfHeaders(web.config)
	headers["Accept"] = "application/json;odata=verbose"

	body := fmt.Sprintf(`{"logonName": "%s"}`, loginName)

	data, err := client.Post(endpoint, bytes.NewBuffer([]byte(body)), patchConfigHeaders(web.config, headers))
	if err != nil {
		return nil, err
	}

	res := &struct {
		User *UserInfo `json:"d"`
	}{}

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, fmt.Errorf("unable to parse the response: %w", err)
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

// UserInfoList gets site UIL (User Information List) API instance
func (web *Web) UserInfoList() *List {
	return NewList(
		web.client,
		fmt.Sprintf("%s/SiteUserInfoList", web.endpoint),
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
// or web relevant (e.g. `lib/folder`, with web relevant URI there should be no slash at the beginning)
// A wrapper of `GetFolderByServerRelativeUrl`
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

// GetFolderByPath gets a folder by its relevant URI, URI can be host relevant (e.g. `/sites/site/lib/folder`)
// or web relevant (e.g. `lib/folder`, with web relevant URI there should be no slash at the beginning)
// A wrapper of `GetFolderByServerRelativePath`
// Supported only in modern SharePoint, differs from GetFile with its capability of dealing with special chars in path
func (web *Web) GetFolderByPath(serverRelativeURL string) *Folder {
	return NewFolder(
		web.client,
		fmt.Sprintf(
			"%s/GetFolderByServerRelativePath(decodedUrl='%s')",
			web.endpoint,
			checkGetRelativeURL(serverRelativeURL, web.endpoint),
		),
		web.config,
	)
}

// GetFolderByID gets Folder API instance object by its unique ID
// Supported only in modern SharePoint
func (web *Web) GetFolderByID(uniqueID string) *Folder {
	return NewFolder(
		web.client,
		fmt.Sprintf("%s/GetFolderByID('%s')", web.endpoint, uniqueID),
		web.config,
	)
}

// EnsureFolder is a helper to ensure a folder by its relevant URI, when there was no folder it's created
func (web *Web) EnsureFolder(serverRelativeURL string) ([]byte, error) {
	return ensureFolder(web, serverRelativeURL, serverRelativeURL)
}

// GetFile gets File API instance object by its relevant URI
// File URI can be host relevant (e.g. `/sites/site/lib/folder/file.txt`)
// or web relevant (e.g. `lib/folder/file.txt`, with web relevant URI there should be no slash at the beginning)
// A wrapper of `GetFileByServerRelativeUrl`
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

// GetFileByPath gets File API instance object by its relevant URI
// File URI can be host relevant (e.g. `/sites/site/lib/folder/file.txt`)
// or web relevant (e.g. `lib/folder/file.txt`, with web relevant URI there should be no slash at the beginning)
// A wrapper of `GetFileByServerRelativePath`
// Supported only in modern SharePoint, differs from GetFile with its capability of dealing with special chars in path
func (web *Web) GetFileByPath(serverRelativeURL string) *File {
	return NewFile(
		web.client,
		fmt.Sprintf(
			"%s/GetFileByServerRelativePath(decodedUrl='%s')",
			web.endpoint,
			checkGetRelativeURL(serverRelativeURL, web.endpoint),
		),
		web.config,
	)
}

// GetFileByID gets File API instance object by its unique ID
// Supported only in modern SharePoint
func (web *Web) GetFileByID(uniqueID string) *File {
	return NewFile(
		web.client,
		fmt.Sprintf("GetFileByID('%s')", uniqueID),
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

// EventReceivers gets EventReceivers API scoped object
func (web *Web) EventReceivers() *EventReceivers {
	return NewEventReceivers(
		web.client,
		fmt.Sprintf("%s/EventReceivers", web.endpoint),
		web.config,
	)
}

// CustomActions gets CustomActions API scoped object
func (web *Web) CustomActions() *CustomActions {
	return NewCustomActions(
		web.client,
		fmt.Sprintf("%s/UserCustomActions", web.endpoint),
		web.config,
	)
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
