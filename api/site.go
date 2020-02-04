package api

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/koltyakov/gosip"
)

//go:generate ggen -ent Site -conf -mods Select,Expand -helpers Data,Normalized

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

// FromURL gets Site object using its API URL
func (site *Site) FromURL(url string) *Site {
	url = strings.Split(url, "?")[0]
	return NewSite(site.client, url, site.config)
}

// Get gets this Site data object
func (site *Site) Get() (SiteResp, error) {
	sp := NewHTTPClient(site.client)
	return sp.Get(site.ToURL(), getConfHeaders(site.config))
}

// Update updates Site's metadata with properties provided in `body` parameter
// where `body` is byte array representation of JSON string payload relevalt to SP.Site object
func (site *Site) Update(body []byte) (SiteResp, error) {
	body = patchMetadataType(body, "SP.Site")
	sp := NewHTTPClient(site.client)
	return sp.Update(site.endpoint, bytes.NewBuffer(body), getConfHeaders(site.config))
}

// Delete deletes current site (can't be restored from a recycle bin)
func (site *Site) Delete() error {
	sp := NewHTTPClient(site.client)
	_, err := sp.Delete(site.endpoint, getConfHeaders(site.config))
	return err
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

// Changes gets changes API scoped object
func (site *Site) Changes() *Changes {
	return NewChanges(
		site.client,
		site.endpoint,
		site.config,
	)
}

// EventReceivers gets EventReceivers API scoped object
func (site *Site) EventReceivers() *EventReceivers {
	return NewEventReceivers(
		site.client,
		fmt.Sprintf("%s/EventReceivers", site.endpoint),
		site.config,
	)
}

// CustomActions gets CustomActions API scoped object
func (site *Site) CustomActions() *CustomActions {
	return NewCustomActions(
		site.client,
		fmt.Sprintf("%s/UserCustomActions", site.endpoint),
		site.config,
	)
}

// Owner gets site's owner user
func (site *Site) Owner() *User {
	return NewUser(
		site.client,
		fmt.Sprintf("%s/Owner", site.endpoint),
		site.config,
	)
}
