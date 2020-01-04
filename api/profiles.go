package api

import (
	"encoding/json"
	"net/url"

	"github.com/koltyakov/gosip"
)

// Profiles  represent SharePoint User Profiles API queryable collection struct
// Always use NewProfiles constructor instead of &Profiles{}
type Profiles struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers *ODataMods
}

// ProfileInfo - user profile API response payload structure
type ProfileInfo struct {
	AccountName                     string `json:"AccountName"`
	DisplayName                     string `json:"DisplayName"`
	FollowPersonalSiteURL           string `json:"FollowPersonalSiteUrl"`
	IsDefaultDocumentLibraryBlocked bool   `json:"IsDefaultDocumentLibraryBlocked"`
	IsPeopleListPublic              bool   `json:"IsPeopleListPublic"`
	IsPrivacySettingOn              bool   `json:"IsPrivacySettingOn"`
	IsSelf                          bool   `json:"IsSelf"`
	JobTitle                        string `json:"JobTitle"`
	MySiteFirstRunExperience        int    `json:"MySiteFirstRunExperience"`
	MySiteHostURL                   string `json:"MySiteHostUrl"`
	O15FirstRunExperience           int    `json:"O15FirstRunExperience"`
	PersonalSiteCapabilities        int    `json:"PersonalSiteCapabilities"`
	PersonalSiteFirstCreationError  string `json:"PersonalSiteFirstCreationError"`
	PersonalSiteFirstCreationTime   string `json:"PersonalSiteFirstCreationTime"`
	PersonalSiteInstantiationState  int    `json:"PersonalSiteInstantiationState"`
	PersonalSiteLastCreationTime    string `json:"PersonalSiteLastCreationTime"`
	PersonalSiteNumberOfRetries     int    `json:"PersonalSiteNumberOfRetries"`
	PictureImportEnabled            bool   `json:"PictureImportEnabled"`
	PictureURL                      string `json:"PictureUrl"`
	PublicURL                       string `json:"PublicUrl"`
	SipAddress                      string `json:"SipAddress"`
	URLToCreatePersonalSite         string `json:"UrlToCreatePersonalSite"`
}

// ProfilePropsInto ...
type ProfilePropsInto struct {
	AccountName           string           `json:"AccountName"`
	DirectReports         []string         `json:"DirectReports"`
	DisplayName           string           `json:"DisplayName"`
	Email                 string           `json:"Email"`
	ExtendedManagers      []string         `json:"ExtendedManagers"`
	ExtendedReports       []string         `json:"ExtendedReports"`
	Peers                 []string         `json:"Peers"`
	IsFollowed            bool             `json:"IsFollowed"`
	PersonalSiteHostURL   string           `json:"PersonalSiteHostUrl"`
	PersonalURL           string           `json:"PersonalUrl"`
	PictureURL            string           `json:"PictureUrl"`
	Title                 string           `json:"Title"`
	UserURL               string           `json:"UserUrl"`
	UserProfileProperties []*TypedKeyValue `json:"UserProfileProperties"`
	// LatestPost       string `json:"LatestPost"`
}

// ProfileResp - user profile response type with helper processor methods
type ProfileResp []byte

// ProfilePropsResp - user profile props response type with helper processor methods
type ProfilePropsResp []byte

// NewProfiles - Profiles struct constructor function
func NewProfiles(client *gosip.SPClient, endpoint string, config *RequestConfig) *Profiles {
	return &Profiles{
		client:    client,
		endpoint:  endpoint,
		config:    config,
		modifiers: NewODataMods(),
	}
}

// ToURL gets endpoint with modificators raw URL
func (profiles *Profiles) ToURL() string {
	return toURL(profiles.endpoint, profiles.modifiers)
}

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (profiles *Profiles) Conf(config *RequestConfig) *Profiles {
	profiles.config = config
	return profiles
}

// Select adds $select OData modifier
func (profiles *Profiles) Select(oDataSelect string) *Profiles {
	profiles.modifiers.AddSelect(oDataSelect)
	return profiles
}

// Expand adds $expand OData modifier
func (profiles *Profiles) Expand(oDataExpand string) *Profiles {
	profiles.modifiers.AddExpand(oDataExpand)
	return profiles
}

// GetMyProperties gets current context user profile properties
func (profiles *Profiles) GetMyProperties() (ProfilePropsResp, error) {
	sp := NewHTTPClient(profiles.client)
	apiURL, _ := url.Parse(profiles.endpoint + "/GetMyProperties")
	query := apiURL.Query()
	for k, v := range profiles.modifiers.Get() {
		query.Set(k, trimMultiline(v))
	}
	apiURL.RawQuery = query.Encode()
	return sp.Post(apiURL.String(), nil, getConfHeaders(profiles.config))
}

// GetPropertiesFor gets properties of a specified user profile (by user login name)
func (profiles *Profiles) GetPropertiesFor(loginName string) (ProfilePropsResp, error) {
	sp := NewHTTPClient(profiles.client)
	apiURL, _ := url.Parse(
		profiles.endpoint +
			"/GetPropertiesFor('" + url.QueryEscape(loginName) + "')",
	)
	query := apiURL.Query()
	for k, v := range profiles.modifiers.Get() {
		query.Set(k, trimMultiline(v))
	}
	apiURL.RawQuery = query.Encode()
	return sp.Get(apiURL.String(), getConfHeaders(profiles.config))
}

// GetUserProfilePropertyFor gets specific properte of a specified user profile (by user login name)
func (profiles *Profiles) GetUserProfilePropertyFor(loginName string, property string) (string, error) {
	sp := NewHTTPClient(profiles.client)
	endpoint := profiles.endpoint +
		"/GetUserProfilePropertyFor(" +
		"accountname='" + url.QueryEscape(loginName) + "'," +
		"propertyname='" + url.QueryEscape(property) + "')"
	data, err := sp.Get(endpoint, getConfHeaders(profiles.config))
	if err != nil {
		return "", err
	}
	data = parseODataItem(data)
	res := &struct {
		Value                     string `json:"value"`
		GetUserProfilePropertyFor string `json:"GetUserProfilePropertyFor"`
	}{}
	if err := json.Unmarshal(data, &res); err != nil {
		return "", err
	}
	propertyValue := res.GetUserProfilePropertyFor
	if res.Value != "" {
		propertyValue = res.Value
	}
	return propertyValue, nil
}

// GetOwnerUserProfile gets owner's user profile
func (profiles *Profiles) GetOwnerUserProfile() (ProfileResp, error) {
	sp := NewHTTPClient(profiles.client)
	apiURL, _ := url.Parse(
		getPriorEndpoint(profiles.endpoint, "/_api") +
			"/_api/sp.userprofiles.profileloader.getowneruserprofile",
	)
	query := apiURL.Query()
	for k, v := range profiles.modifiers.Get() {
		query.Set(k, trimMultiline(v))
	}
	apiURL.RawQuery = query.Encode()
	return sp.Post(apiURL.String(), nil, getConfHeaders(profiles.config))
}

// UserProfile gets current context user profile object
func (profiles *Profiles) UserProfile() (ProfileResp, error) {
	sp := NewHTTPClient(profiles.client)
	apiURL, _ := url.Parse(
		getPriorEndpoint(profiles.endpoint, "/_api") +
			"/_api/sp.userprofiles.profileloader.getprofileloader/GetUserProfile",
	)
	query := apiURL.Query()
	for k, v := range profiles.modifiers.Get() {
		query.Set(k, trimMultiline(v))
	}
	apiURL.RawQuery = query.Encode()
	return sp.Post(apiURL.String(), nil, getConfHeaders(profiles.config))
}

// SetSingleValueProfileProperty sets a single value property for the profile by its email
func (profiles *Profiles) SetSingleValueProfileProperty(loginName string, property string, value string) error {
	sp := NewHTTPClient(profiles.client)
	endpoint := profiles.endpoint + "/SetSingleValueProfileProperty"
	prop := map[string]string{}
	prop["accountName"] = loginName
	prop["propertyName"] = property
	prop["propertyValue"] = value
	body, _ := json.Marshal(prop)
	_, err := sp.Post(endpoint, body, getConfHeaders(profiles.config))
	return err
}

// SetMultiValuedProfileProperty sets a multi value property for the profile by its email
func (profiles *Profiles) SetMultiValuedProfileProperty(loginName string, property string, values []string) error {
	sp := NewHTTPClient(profiles.client)
	endpoint := profiles.endpoint + "/SetMultiValuedProfileProperty"
	prop := map[string]interface{}{}
	prop["accountName"] = loginName
	prop["propertyName"] = property
	prop["propertyValues"] = values
	body, _ := json.Marshal(prop)
	_, err := sp.Post(endpoint, body, getConfHeaders(profiles.config))
	return err
}

// HideSuggestion removes the specified user from the user's list of suggested people to follow
func (profiles *Profiles) HideSuggestion(loginName string) ([]byte, error) {
	sp := NewHTTPClient(profiles.client)
	endpoint := profiles.endpoint + "/HideSuggestion('" + url.QueryEscape(loginName) + "')"
	return sp.Post(endpoint, nil, getConfHeaders(profiles.config))
}

/* Response helpers */

// Data : to get typed data
func (profileResp *ProfileResp) Data() *ProfileInfo {
	data := parseODataItem(*profileResp)
	data = normalizeMultiLookups(data)
	res := &ProfileInfo{}
	json.Unmarshal(data, &res)
	return res
}

// Unmarshal : to unmarshal to custom object
func (profileResp *ProfileResp) Unmarshal(obj interface{}) error {
	data := parseODataItem(*profileResp)
	return json.Unmarshal(data, obj)
}

// Data : to get typed data
func (profilePropsResp *ProfilePropsResp) Data() *ProfilePropsInto {
	data := parseODataItem(*profilePropsResp)
	data = normalizeMultiLookups(data)
	res := &ProfilePropsInto{}
	json.Unmarshal(data, &res)
	return res
}

// Unmarshal : to unmarshal to custom object
func (profilePropsResp *ProfilePropsResp) Unmarshal(obj interface{}) error {
	data := parseODataItem(*profilePropsResp)
	return json.Unmarshal(data, obj)
}
