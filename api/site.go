package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/koltyakov/gosip"
)

// Site ...
type Site struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers map[string]string
}

// SiteInfo ...
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

// SiteResp ...
type SiteResp []byte

// NewSite ...
func NewSite(client *gosip.SPClient, endpoint string, config *RequestConfig) *Site {
	return &Site{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// ToURL ...
func (site *Site) ToURL() string {
	apiURL, _ := url.Parse(site.endpoint)
	query := apiURL.Query() // url.Values{}
	for k, v := range site.modifiers {
		query.Set(k, trimMultiline(v))
	}
	apiURL.RawQuery = query.Encode()
	return apiURL.String()
}

// Conf ...
func (site *Site) Conf(config *RequestConfig) *Site {
	site.config = config
	return site
}

// Select ...
func (site *Site) Select(oDataSelect string) *Site {
	if site.modifiers == nil {
		site.modifiers = make(map[string]string)
	}
	site.modifiers["$select"] = oDataSelect
	return site
}

// Expand ...
func (site *Site) Expand(oDataExpand string) *Site {
	if site.modifiers == nil {
		site.modifiers = make(map[string]string)
	}
	site.modifiers["$expand"] = oDataExpand
	return site
}

// Get ...
func (site *Site) Get() (SiteResp, error) {
	sp := NewHTTPClient(site.client)
	return sp.Get(site.ToURL(), getConfHeaders(site.config))
}

// Delete ...
func (site *Site) Delete() ([]byte, error) {
	sp := NewHTTPClient(site.client)
	return sp.Delete(site.endpoint, getConfHeaders(site.config))
}

// RootWeb ...
func (site *Site) RootWeb() *Web {
	endpoint := fmt.Sprintf("%s/RootWeb", site.endpoint)
	return NewWeb(site.client, endpoint, site.config)
}

// OpenWebByID ...
func (site *Site) OpenWebByID(webID string) (WebResp, error) {
	endpoint := fmt.Sprintf("%s/OpenWebById('%s')", site.endpoint, webID)
	sp := NewHTTPClient(site.client)
	return sp.Post(endpoint, nil, getConfHeaders(site.config))
}

// GetChangeToken ...
func (site *Site) GetChangeToken() (string, error) {
	scoped := *site
	data, err := scoped.Select("CurrentChangeToken").Get()
	if err != nil {
		return "", err
	}
	return data.Data().CurrentChangeToken.StringValue, nil
}

// GetChanges ...
func (site *Site) GetChanges(changeQuery *ChangeQuery) ([]*ChangeInfo, error) {
	return NewChanges(
		site.client,
		fmt.Sprintf("%s/GetChanges", site.endpoint),
		site.config,
	).GetChanges(changeQuery)
}

// ToDo:
// Features
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
