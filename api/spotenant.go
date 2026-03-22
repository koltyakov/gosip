package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/csom"
)

// SPOTenant represents SharePoint Online Tenant Administration API
// Always use NewSPOTenant constructor instead of &SPOTenant{}
type SPOTenant struct {
	client   *HTTPClient
	config   *RequestConfig
	endpoint string
}

// NewSPOTenant creates a new SPOTenant API instance
func NewSPOTenant(client *gosip.SPClient, siteURL string, config *RequestConfig) *SPOTenant {
	return &SPOTenant{
		client:   NewHTTPClient(client),
		endpoint: siteURL,
		config:   config,
	}
}

// Conf sets request config for the current chain
func (t *SPOTenant) Conf(config *RequestConfig) *SPOTenant {
	t.config = config
	return t
}

// IncludePersonalSite controls whether OneDrive personal sites appear
// in tenant site property queries.

// IncludePersonalSiteDefault uses the server default (excludes personal sites)
// IncludePersonalSiteInclude includes personal sites
// IncludePersonalSiteOnly returns only personal sites
type IncludePersonalSite int

const (
	IncludePersonalSiteDefault IncludePersonalSite = 0
	IncludePersonalSiteInclude IncludePersonalSite = 1
	IncludePersonalSiteOnly IncludePersonalSite = 2
)

// SitePropertiesFilter defines the filter parameters for GetSitePropertiesFromSharePointByFilters
type SitePropertiesFilter struct {

	// CSOM filter expression (e.g., "Url -like 'tenant-my.sharepoint.com/personal/'")
	Filter string

	// IncludeDetail requests the full set of site properties. When false,
	// some properties may return server defaults rather than actual values
	IncludeDetail bool

	// Site template filter
	// (e.g., use "SPSPERS" for OneDrive personal sites)
	Template string

	// IncludePersonalSite controls whether OneDrive personal sites appear in results.
	// Valid values: IncludePersonalSiteDefault (0), IncludePersonalSiteInclude (1), IncludePersonalSiteOnly (2).
	IncludePersonalSite IncludePersonalSite
}

// SitePropertiesPage represents a page of site property results
type SitePropertiesPage struct {
	Sites       []SitePropertiesResp
	HasNextPage func() bool
	GetNextPage func() (*SitePropertiesPage, error)
}

// SitePropertiesResp represents a single sites raw CSOM JSON response
type SitePropertiesResp []byte

// Data unmarshals the response into a typed SiteProperties struct
func (r SitePropertiesResp) Data() *SiteProperties {
	var sp SiteProperties
	if err := json.Unmarshal(r, &sp); err != nil {
		return &SiteProperties{}
	}
	return &sp
}

// Normalized returns the raw JSON bytes for custom unmarshalling
func (r SitePropertiesResp) Normalized() []byte {
	return []byte(r)
}

// SiteProperties holds site metadata returned by CSOM tenant queries.
// Fields map to Microsoft.Online.SharePoint.TenantAdministration.SiteProperties.
type SiteProperties struct {
	AllowDownloadingNonWebViewableFiles                        bool        `json:"AllowDownloadingNonWebViewableFiles"`
	AllowEditing                                               bool        `json:"AllowEditing"`
	AllowFileArchive                                           bool        `json:"AllowFileArchive"`
	AllowSelfServiceUpgrade                                    bool        `json:"AllowSelfServiceUpgrade"`
	AllowWebPropertyBagUpdateWhenDenyAddAndCustomizePagesIsEnabled bool     `json:"AllowWebPropertyBagUpdateWhenDenyAddAndCustomizePagesIsEnabled"`
	AnonymousLinkExpirationInDays                              int         `json:"AnonymousLinkExpirationInDays"`
	ApplyToExistingDocumentLibraries                           bool        `json:"ApplyToExistingDocumentLibraries"`
	ApplyToNewDocumentLibraries                                bool        `json:"ApplyToNewDocumentLibraries"`
	ArchiveStatus                                              string      `json:"ArchiveStatus"`
	ArchivedBy                                                 string      `json:"ArchivedBy"`
	ArchivedFileDiskUsed                                       int64       `json:"ArchivedFileDiskUsed"`
	ArchivedTime                                               string      `json:"ArchivedTime"`
	AuthContextStrength                                        *string     `json:"AuthContextStrength"`
	AuthenticationContextLimitedAccess                         bool        `json:"AuthenticationContextLimitedAccess"`
	AuthenticationContextName                                  *string     `json:"AuthenticationContextName"`
	AverageResourceUsage                                       int         `json:"AverageResourceUsage"`
	BlockDownloadLinksFileType                                 int         `json:"BlockDownloadLinksFileType"`
	BlockDownloadMicrosoft365GroupIds                           []string    `json:"BlockDownloadMicrosoft365GroupIds"`
	BlockDownloadPolicy                                        bool        `json:"BlockDownloadPolicy"`
	BlockDownloadPolicyFileTypeIds                             []int       `json:"BlockDownloadPolicyFileTypeIds"`
	BlockGuestsAsSiteAdmin                                     int         `json:"BlockGuestsAsSiteAdmin"`
	BonusDiskQuota                                             int64       `json:"BonusDiskQuota"`
	ClearGroupId                                               bool        `json:"ClearGroupId"`
	ClearRestrictedAccessControl                               bool        `json:"ClearRestrictedAccessControl"`
	CommentsOnSitePagesDisabled                                bool        `json:"CommentsOnSitePagesDisabled"`
	CompatibilityLevel                                         int         `json:"CompatibilityLevel"`
	ConditionalAccessPolicy                                    int         `json:"ConditionalAccessPolicy"`
	CreatedTime                                                string      `json:"CreatedTime"`
	CurrentResourceUsage                                       int         `json:"CurrentResourceUsage"`
	DefaultLinkPermission                                      int         `json:"DefaultLinkPermission"`
	DefaultLinkToExistingAccess                                bool        `json:"DefaultLinkToExistingAccess"`
	DefaultLinkToExistingAccessReset                           bool        `json:"DefaultLinkToExistingAccessReset"`
	DefaultShareLinkRole                                       int         `json:"DefaultShareLinkRole"`
	DefaultShareLinkScope                                      int         `json:"DefaultShareLinkScope"`
	DefaultSharingLinkType                                     int         `json:"DefaultSharingLinkType"`
	DenyAddAndCustomizePages                                   int         `json:"DenyAddAndCustomizePages"`
	Description                                                *string     `json:"Description"`
	DisableAppViews                                            int         `json:"DisableAppViews"`
	DisableClassicPageBaselineSecurityMode                     bool        `json:"DisableClassicPageBaselineSecurityMode"`
	DisableCompanyWideSharingLinks                             int         `json:"DisableCompanyWideSharingLinks"`
	DisableFlows                                               int         `json:"DisableFlows"`
	DisableSiteBranding                                        bool        `json:"DisableSiteBranding"`
	EnableAutoExpirationVersionTrim                            bool        `json:"EnableAutoExpirationVersionTrim"`
	ExcludeBlockDownloadPolicySiteOwners                       bool        `json:"ExcludeBlockDownloadPolicySiteOwners"`
	ExcludeBlockDownloadSharePointGroups                       []string    `json:"ExcludeBlockDownloadSharePointGroups"`
	ExcludedBlockDownloadGroupIds                              []string    `json:"ExcludedBlockDownloadGroupIds"`
	ExpireVersionsAfterDays                                    int         `json:"ExpireVersionsAfterDays"`
	ExternalUserExpirationInDays                               int         `json:"ExternalUserExpirationInDays"`
	FileTypesForVersionExpiration                              *string     `json:"FileTypesForVersionExpiration"`
	GroupId                                                    string      `json:"GroupId"`
	GroupOwnerLoginName                                        *string     `json:"GroupOwnerLoginName"`
	HasHolds                                                   bool        `json:"HasHolds"`
	HidePeoplePreviewingFiles                                  bool        `json:"HidePeoplePreviewingFiles"`
	HidePeopleWhoHaveListsOpen                                 bool        `json:"HidePeopleWhoHaveListsOpen"`
	HubSiteId                                                  string      `json:"HubSiteId"`
	IBMode                                                     *string     `json:"IBMode"`
	IBSegments                                                 []string    `json:"IBSegments"`
	IBSegmentsToAdd                                            []string    `json:"IBSegmentsToAdd"`
	IBSegmentsToRemove                                         []string    `json:"IBSegmentsToRemove"`
	InheritVersionPolicyFromTenant                             bool        `json:"InheritVersionPolicyFromTenant"`
	IsAuthoritative                                            bool        `json:"IsAuthoritative"`
	IsGroupOwnerSiteAdmin                                      bool        `json:"IsGroupOwnerSiteAdmin"`
	IsHubSite                                                  bool        `json:"IsHubSite"`
	IsTeamsChannelConnected                                    bool        `json:"IsTeamsChannelConnected"`
	IsTeamsConnected                                           bool        `json:"IsTeamsConnected"`
	LastContentModifiedDate                                    string      `json:"LastContentModifiedDate"`
	Lcid                                                       int         `json:"Lcid"`
	LimitedAccessFileType                                      int         `json:"LimitedAccessFileType"`
	ListsShowHeaderAndNavigation                               bool        `json:"ListsShowHeaderAndNavigation"`
	LockIssue                                                  *string     `json:"LockIssue"`
	LockReason                                                 int         `json:"LockReason"`
	LockState                                                  string      `json:"LockState"`
	LoopDefaultSharingLinkRole                                 int         `json:"LoopDefaultSharingLinkRole"`
	LoopDefaultSharingLinkScope                                int         `json:"LoopDefaultSharingLinkScope"`
	MajorVersionLimit                                          int         `json:"MajorVersionLimit"`
	MajorWithMinorVersionsLimit                                int         `json:"MajorWithMinorVersionsLimit"`
	MediaTranscription                                         int         `json:"MediaTranscription"`
	OrganizationLinkMaxExpirationInDays                        int         `json:"OrganizationLinkMaxExpirationInDays"`
	OrganizationLinkRecommendedExpirationInDays                int         `json:"OrganizationLinkRecommendedExpirationInDays"`
	OverrideBlockUserInfoVisibility                            int         `json:"OverrideBlockUserInfoVisibility"`
	OverrideSharingCapability                                  bool        `json:"OverrideSharingCapability"`
	OverrideTenantAnonymousLinkExpirationPolicy                bool        `json:"OverrideTenantAnonymousLinkExpirationPolicy"`
	OverrideTenantExternalUserExpirationPolicy                 bool        `json:"OverrideTenantExternalUserExpirationPolicy"`
	OverrideTenantOrganizationLinkExpirationPolicy             bool        `json:"OverrideTenantOrganizationLinkExpirationPolicy"`
	Owner                                                      string      `json:"Owner"`
	OwnerEmail                                                 *string     `json:"OwnerEmail"`
	OwnerLoginName                                             *string     `json:"OwnerLoginName"`
	OwnerName                                                  *string     `json:"OwnerName"`
	PWAEnabled                                                 int         `json:"PWAEnabled"`
	ReadOnlyAccessPolicy                                       bool        `json:"ReadOnlyAccessPolicy"`
	ReadOnlyForBlockDownloadPolicy                             bool        `json:"ReadOnlyForBlockDownloadPolicy"`
	ReadOnlyForUnmanagedDevices                                bool        `json:"ReadOnlyForUnmanagedDevices"`
	RelatedGroupId                                             string      `json:"RelatedGroupId"`
	RemoveVersionExpirationFileTypeOverride                    *string     `json:"RemoveVersionExpirationFileTypeOverride"`
	RequestFilesLinkEnabled                                    bool        `json:"RequestFilesLinkEnabled"`
	RequestFilesLinkExpirationInDays                           int         `json:"RequestFilesLinkExpirationInDays"`
	RestrictContentOrgWideSearch                               bool        `json:"RestrictContentOrgWideSearch"`
	RestrictedAccessControl                                    bool        `json:"RestrictedAccessControl"`
	RestrictedAccessControlGroups                              []string    `json:"RestrictedAccessControlGroups"`
	RestrictedAccessControlGroupsToAdd                         []string    `json:"RestrictedAccessControlGroupsToAdd"`
	RestrictedAccessControlGroupsToRemove                      []string    `json:"RestrictedAccessControlGroupsToRemove"`
	RestrictedContentDiscoveryForCopilotAndAgents              bool        `json:"RestrictedContentDiscoveryforCopilotAndAgents"`
	RestrictedToRegion                                         int         `json:"RestrictedToRegion"`
	SandboxedCodeActivationCapability                          int         `json:"SandboxedCodeActivationCapability"`
	SensitivityLabel                                           string      `json:"SensitivityLabel"`
	SensitivityLabel2                                          *string     `json:"SensitivityLabel2"`
	SetOwnerWithoutUpdatingSecondaryAdmin                      bool        `json:"SetOwnerWithoutUpdatingSecondaryAdmin"`
	SharingAllowedDomainList                                   *string     `json:"SharingAllowedDomainList"`
	SharingBlockedDomainList                                   *string     `json:"SharingBlockedDomainList"`
	SharingCapability                                          int         `json:"SharingCapability"`
	SharingDomainRestrictionMode                               int         `json:"SharingDomainRestrictionMode"`
	SharingLockDownCanBeCleared                                bool        `json:"SharingLockDownCanBeCleared"`
	SharingLockDownEnabled                                     bool        `json:"SharingLockDownEnabled"`
	ShowPeoplePickerSuggestionsForGuestUsers                   bool        `json:"ShowPeoplePickerSuggestionsForGuestUsers"`
	SiteDefinedSharingCapability                               int         `json:"SiteDefinedSharingCapability"`
	SiteID                                                     string      `json:"SiteId"`
	SocialBarOnSitePagesDisabled                               bool        `json:"SocialBarOnSitePagesDisabled"`
	Status                                                     string      `json:"Status"`
	StorageMaximumLevel                                        int64       `json:"StorageMaximumLevel"`
	StorageQuotaType                                           *string     `json:"StorageQuotaType"`
	StorageUsage                                               int64       `json:"StorageUsage"`
	StorageWarningLevel                                        int64       `json:"StorageWarningLevel"`
	TeamsChannelType                                           int         `json:"TeamsChannelType"`
	Template                                                   string      `json:"Template"`
	TimeZoneId                                                 int         `json:"TimeZoneId"`
	Title                                                      string      `json:"Title"`
	TitleTranslations                                          []string    `json:"TitleTranslations"`
	URL                                                        string      `json:"Url"`
	UserCodeMaximumLevel                                       int         `json:"UserCodeMaximumLevel"`
	UserCodeWarningLevel                                       int         `json:"UserCodeWarningLevel"`
	VersionCount                                               int         `json:"VersionCount"`
	VersionPolicyFileTypeOverride                              []string    `json:"VersionPolicyFileTypeOverride"`
	VersionSize                                                int64       `json:"VersionSize"`
	WebsCount                                                  int         `json:"WebsCount"`
}

// CreatedTimeDate parses CreatedTime from CSOM date format into time.Time.
// Returns the zero time if the format is unrecognized.
func (s *SiteProperties) CreatedTimeDate() time.Time {
	return parseCSOMDate(s.CreatedTime)
}

// LastContentModifiedDateTime parses LastContentModifiedDate from CSOM date format into time.Time.
// Returns the zero time if the format is unrecognized.
func (s *SiteProperties) LastContentModifiedDateTime() time.Time {
	return parseCSOMDate(s.LastContentModifiedDate)
}

// CleanSiteID returns the SiteID with any "/Guid(...)/" wrapper stripped.
func (s *SiteProperties) CleanSiteID() string {
	clean := strings.TrimPrefix(s.SiteID, "/Guid(")
	clean = strings.TrimSuffix(clean, ")/")
	return clean
}

// GetSiteProperties returns a filtered list of site collections and their properties using the
// GetSitePropertiesFromSharePointByFilters CSOM method.
// Results will include OneDrive personal sites when the filter's IncludePersonalSite is set.
// For non-personal site queries, the sp.Web() and sp.Site() fluent methods are the recommended approach.
func (t *SPOTenant) GetSiteProperties(filter *SitePropertiesFilter) (*SitePropertiesPage, error) {
	if filter == nil {
		filter = &SitePropertiesFilter{}
	}
	return t.getSitePropertiesPage(filter, "0")
}

func (t *SPOTenant) getSitePropertiesPage(filter *SitePropertiesFilter, startIndex string) (*SitePropertiesPage, error) {
	csomPkg, err := buildSitePropertiesQuery(filter, startIndex)
	if err != nil {
		return nil, fmt.Errorf("failed to build CSOM query: %w", err)
	}

	jsomResp, err := t.client.ProcessQuery(t.endpoint, bytes.NewBuffer([]byte(csomPkg)), t.config)
	if err != nil {
		return nil, fmt.Errorf("CSOM query failed: %w", err)
	}

	resp, err := parseSitePropertiesResponse(jsomResp)
	if err != nil {
		return nil, err
	}

	nextStartIndex := resp.nextStartIndex
	page := &SitePropertiesPage{
		Sites: resp.sites,
		HasNextPage: func() bool {
			return strings.TrimSpace(nextStartIndex) != ""
		},
		GetNextPage: func() (*SitePropertiesPage, error) {
			if strings.TrimSpace(nextStartIndex) == "" {
				return nil, fmt.Errorf("no next page available")
			}
			return t.getSitePropertiesPage(filter, nextStartIndex)
		},
	}
	return page, nil
}

// SPO Tenant CSOM TypeIDs
const (
	spoTenantTypeID              = "{268004ae-ef6b-4e9b-8425-127220d84719}"
	spoSitePropertiesFilterTypeID = "{b92aeee2-c92c-4b67-abcc-024e471bc140}"
)

// buildSitePropertiesQuery constructs the CSOM XML for GetSitePropertiesFromSharePointByFilters
func buildSitePropertiesQuery(filter *SitePropertiesFilter, startIndex string) (string, error) {
	builder := csom.NewBuilder()

	constructorXML := `<Constructor Id="{{.ID}}" TypeId="` + spoTenantTypeID + `" />`
	tenantObject, _ := builder.AddObject(csom.NewObject(constructorXML), nil)

	methodParams := []string{
		fmt.Sprintf(`
			<Parameter TypeId="%s">
				<Property Name="Filter" Type="String">%s</Property>
				<Property Name="IncludeDetail" Type="Boolean">%s</Property>
				<Property Name="IncludePersonalSite" Type="Enum">%d</Property>
				<Property Name="StartIndex" Type="String">%s</Property>
				<Property Name="Template" Type="String">%s</Property>
			</Parameter>`,
			spoSitePropertiesFilterTypeID,
			filter.Filter,
			strings.ToLower(strconv.FormatBool(filter.IncludeDetail)),
			filter.IncludePersonalSite,
			startIndex,
			filter.Template,
		),
	}

	methodObject, _ := builder.AddObject(
		csom.NewObjectMethod("GetSitePropertiesFromSharePointByFilters", methodParams),
		tenantObject,
	)

	builder.AddAction(csom.NewAction(`<ObjectPath Id="{{.ID}}" ObjectPathId="{{.ObjectID}}" />`), tenantObject)
	builder.AddAction(csom.NewAction(`<ObjectPath Id="{{.ID}}" ObjectPathId="{{.ObjectID}}" />`), methodObject)

	queryXML := `<Query Id="{{.ID}}" ObjectPathId="{{.ObjectID}}">` +
		`<Query SelectAllProperties="true"><Properties /></Query>` +
		`<ChildItemQuery SelectAllProperties="true"><Properties /></ChildItemQuery>` +
		`</Query>`
	builder.AddAction(csom.NewAction(queryXML), methodObject)

	return builder.Compile()
}

// sitePropertiesRawResponse holds the parsed CSOM response fields
// before they are wrapped into a SitePropertiesPage.
type sitePropertiesRawResponse struct {
	nextStartIndex string
	sites          []SitePropertiesResp
}

// parseSitePropertiesResponse parses the raw CSOM JSON array response
func parseSitePropertiesResponse(data []byte) (*sitePropertiesRawResponse, error) {
	var rawResponse []json.RawMessage
	if err := json.Unmarshal(data, &rawResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal CSOM response: %w", err)
	}

	if len(rawResponse) == 0 {
		return nil, fmt.Errorf("empty CSOM response")
	}

	// The site data is in the last element of the response array
	var envelope struct {
		NextStartIndex string            `json:"NextStartIndexFromSharePoint"`
		ChildItems     []json.RawMessage `json:"_Child_Items_"`
	}
	if err := json.Unmarshal(rawResponse[len(rawResponse)-1], &envelope); err != nil {
		return nil, fmt.Errorf("failed to unmarshal site properties envelope: %w", err)
	}

	sites := make([]SitePropertiesResp, len(envelope.ChildItems))
	for i, raw := range envelope.ChildItems {
		sites[i] = SitePropertiesResp(raw)
	}

	return &sitePropertiesRawResponse{
		nextStartIndex: envelope.NextStartIndex,
		sites:          sites,
	}, nil
}

// csomDateRegex matches the CSOM date format: /Date(year,month,day,hour,min,sec,ms)/
var csomDateRegex = regexp.MustCompile(`/Date\((\d+),(\d+),(\d+),(\d+),(\d+),(\d+),(\d+)\)/`)

// parseCSOMDate parses a SharePoint CSOM date string "/Date(2024,0,15,10,30,0,0)/"
// into a time.Time in UTC. The month field is zero based (0=January).
// Returns the zero time if the input is malformed.
func parseCSOMDate(input string) time.Time {
	matches := csomDateRegex.FindStringSubmatch(input)
	if len(matches) != 8 {
		return time.Time{}
	}

	values := make([]int, 7)
	for i := 0; i < 7; i++ {
		n, err := strconv.Atoi(matches[i+1])
		if err != nil {
			return time.Time{}
		}
		values[i] = n
	}

	return time.Date(
		values[0], time.Month(values[1]+1), values[2],
		values[3], values[4], values[5], values[6]*1_000_000,
		time.UTC,
	)
}
