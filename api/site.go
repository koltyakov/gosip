package api

import (
	"encoding/json"
	"fmt"

	"github.com/koltyakov/gosip"
)

// Site represents SharePoint Site API queryable object struct
// Always use NewSite constructor instead of &Site{}
type Site struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers *ODataMods
}

// SiteInfo - site API response payload structure
type SiteInfo struct {
	AllowCreateDeclarativeWorkflow         bool         `json:"AllowCreateDeclarativeWorkflow"`
	AllowDesigner                          bool         `json:"AllowDesigner"`
	AllowMasterPageEditing                 bool         `json:"AllowMasterPageEditing"`
	AllowRevertFromTemplate                bool         `json:"AllowRevertFromTemplate"`
	AllowSaveDeclarativeWorkflowAsTemplate bool         `json:"AllowSaveDeclarativeWorkflowAsTemplate"`
	AllowSavePublishDeclarativeWorkflow    bool         `json:"AllowSavePublishDeclarativeWorkflow"`
	AllowSelfServiceUpgrade                bool         `json:"AllowSelfServiceUpgrade"`
	AllowSelfServiceUpgradeEvaluation      bool         `json:"AllowSelfServiceUpgradeEvaluation"`
	AuditLogTrimmingRetention              int          `json:"AuditLogTrimmingRetention"`
	CompatibilityLevel                     int          `json:"CompatibilityLevel"`
	CurrentChangeToken                     *StringValue `json:"CurrentChangeToken"`
	DisableAppViews                        bool         `json:"DisableAppViews"`
	DisableCompanyWideSharingLinks         bool         `json:"DisableCompanyWideSharingLinks"`
	DisableFlows                           bool         `json:"DisableFlows"`
	ExternalSharingTipsEnabled             bool         `json:"ExternalSharingTipsEnabled"`
	GeoLocation                            string       `json:"GeoLocation"`
	GroupID                                string       `json:"GroupId"`
	HubSiteID                              string       `json:"HubSiteId"`
	ID                                     string       `json:"Id"`
	IsHubSite                              bool         `json:"IsHubSite"`
	MaxItemsPerThrottledOperation          int          `json:"MaxItemsPerThrottledOperation"`
	NeedsB2BUpgrade                        bool         `json:"NeedsB2BUpgrade"`
	PrimaryURI                             string       `json:"PrimaryUri"`
	ReadOnly                               bool         `json:"ReadOnly"`
	RequiredDesignerVersion                string       `json:"RequiredDesignerVersion"`
	ResourcePath                           *DecodedURL  `json:"ResourcePath"`
	SandboxedCodeActivationCapability      int          `json:"SandboxedCodeActivationCapability"`
	SensitivityLabel                       string       `json:"SensitivityLabel"`
	SensitivityLabelID                     string       `json:"SensitivityLabelId"`
	ServerRelativeURL                      string       `json:"ServerRelativeUrl"`
	ShareByEmailEnabled                    bool         `json:"ShareByEmailEnabled"`
	ShareByLinkEnabled                     bool         `json:"ShareByLinkEnabled"`
	ShowURLStructure                       bool         `json:"ShowUrlStructure"`
	TrimAuditLog                           bool         `json:"TrimAuditLog"`
	UIVersionConfigurationEnabled          bool         `json:"UIVersionConfigurationEnabled"`
	UpgradeReminderDate                    string       `json:"UpgradeReminderDate"` // time.Time
	UpgradeScheduled                       bool         `json:"UpgradeScheduled"`
	UpgradeScheduledDate                   string       `json:"UpgradeScheduledDate"` // time.Time
	Upgrading                              bool         `json:"Upgrading"`
	URL                                    string       `json:"Url"`
}

// SiteResp - site response type with helper processor methods
type SiteResp []byte

// NewSite - Site struct constructor function
func NewSite(client *gosip.SPClient, endpoint string, config *RequestConfig) *Site {
	return &Site{
		client:    client,
		endpoint:  endpoint,
		config:    config,
		modifiers: NewODataMods(),
	}
}

// ToURL gets endpoint with modificators raw URL
func (site *Site) ToURL() string {
	return toURL(site.endpoint, site.modifiers)
}

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (site *Site) Conf(config *RequestConfig) *Site {
	site.config = config
	return site
}

// Select adds $select OData modifier
func (site *Site) Select(oDataSelect string) *Site {
	site.modifiers.AddSelect(oDataSelect)
	return site
}

// Expand adds $expand OData modifier
func (site *Site) Expand(oDataExpand string) *Site {
	site.modifiers.AddExpand(oDataExpand)
	return site
}

// Get gets this Site data object
func (site *Site) Get() (SiteResp, error) {
	sp := NewHTTPClient(site.client)
	return sp.Get(site.ToURL(), getConfHeaders(site.config))
}

// Delete deletes current site (can't be restored from a recycle bin)
func (site *Site) Delete() ([]byte, error) {
	sp := NewHTTPClient(site.client)
	return sp.Delete(site.endpoint, getConfHeaders(site.config))
}

// RootWeb gets Site's root web queryable API object
func (site *Site) RootWeb() *Web {
	endpoint := fmt.Sprintf("%s/RootWeb", site.endpoint)
	return NewWeb(site.client, endpoint, site.config)
}

// OpenWebByID gets a Web data object by its ID (GUID)
func (site *Site) OpenWebByID(webID string) (WebResp, error) {
	endpoint := fmt.Sprintf("%s/OpenWebById('%s')", site.endpoint, webID)
	sp := NewHTTPClient(site.client)
	return sp.Post(endpoint, nil, getConfHeaders(site.config))
}

// Features gets Features API instance queryable collection for this Site
func (site *Site) Features() *Features {
	endpoint := fmt.Sprintf("%s/Features", site.endpoint)
	return NewFeatures(site.client, endpoint, site.config)
}

// RecycleBin gets RecycleBin API instance object for this Site
func (site *Site) RecycleBin() *RecycleBin {
	endpoint := fmt.Sprintf("%s/RecycleBin", site.endpoint)
	return NewRecycleBin(site.client, endpoint, site.config)
}

// GetChangeToken gets current change token for this Site
func (site *Site) GetChangeToken() (string, error) {
	scoped := NewSite(site.client, site.endpoint, site.config)
	data, err := scoped.Select("CurrentChangeToken").Get()
	if err != nil {
		return "", err
	}
	return data.Data().CurrentChangeToken.StringValue, nil
}

// GetChanges gets changes on this Site due to the configuration provided as `changeQuery` parameter
func (site *Site) GetChanges(changeQuery *ChangeQuery) ([]*ChangeInfo, error) {
	return NewChanges(
		site.client,
		fmt.Sprintf("%s/GetChanges", site.endpoint),
		site.config,
	).GetChanges(changeQuery)
}

// ToDo:
// Custom actions

/* Response helpers */

// Data : to get typed data
func (siteResp *SiteResp) Data() *SiteInfo {
	data := parseODataItem(*siteResp)
	res := &SiteInfo{}
	json.Unmarshal(data, &res)
	return res
}

// Unmarshal : to unmarshal to custom object
func (siteResp *SiteResp) Unmarshal(obj interface{}) error {
	data := parseODataItem(*siteResp)
	return json.Unmarshal(data, obj)
}
