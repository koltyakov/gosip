package api

import (
	"encoding/json"
	"net/url"

	"github.com/koltyakov/gosip"
)

// Profiles ...
type Profiles struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers map[string]string
}

// ProfileInfo ...
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
	AccountName           string   `json:"AccountName"`
	DirectReports         []string `json:"DirectReports"`
	DisplayName           string   `json:"DisplayName"`
	Email                 string   `json:"Email"`
	ExtendedManagers      []string `json:"ExtendedManagers"`
	ExtendedReports       []string `json:"ExtendedReports"`
	Peers                 []string `json:"Peers"`
	IsFollowed            bool     `json:"IsFollowed"`
	PersonalSiteHostURL   string   `json:"PersonalSiteHostUrl"`
	PersonalURL           string   `json:"PersonalUrl"`
	PictureURL            string   `json:"PictureUrl"`
	Title                 string   `json:"Title"`
	UserURL               string   `json:"UserUrl"`
	UserProfileProperties []struct {
		Key       string `json:"Key"`
		Value     string `json:"Value"`
		ValueType string `json:"ValueType"`
	} `json:"UserProfileProperties"`
	// LatestPost       string `json:"LatestPost"`
}

// ProfileResp ...
type ProfileResp []byte

// ProfilePropsResp ...
type ProfilePropsResp []byte

// NewProfiles ...
func NewProfiles(client *gosip.SPClient, endpoint string, config *RequestConfig) *Profiles {
	return &Profiles{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// ToURL ...
func (profiles *Profiles) ToURL() string {
	apiURL, _ := url.Parse(profiles.endpoint)
	query := apiURL.Query() // url.Values{}
	for k, v := range profiles.modifiers {
		query.Set(k, trimMultiline(v))
	}
	apiURL.RawQuery = query.Encode()
	return apiURL.String()
}

// Conf ...
func (profiles *Profiles) Conf(config *RequestConfig) *Profiles {
	profiles.config = config
	return profiles
}

// Select ...
func (profiles *Profiles) Select(oDataSelect string) *Profiles {
	if profiles.modifiers == nil {
		profiles.modifiers = make(map[string]string)
	}
	profiles.modifiers["$select"] = oDataSelect
	return profiles
}

// Expand ...
func (profiles *Profiles) Expand(oDataExpand string) *Profiles {
	if profiles.modifiers == nil {
		profiles.modifiers = make(map[string]string)
	}
	profiles.modifiers["$expand"] = oDataExpand
	return profiles
}

// Get ...
func (profiles *Profiles) Get() ([]byte, error) {
	sp := NewHTTPClient(profiles.client)
	return sp.Get(profiles.ToURL(), getConfHeaders(profiles.config))
}

// GetMyProperties ...
func (profiles *Profiles) GetMyProperties() (ProfilePropsResp, error) {
	sp := NewHTTPClient(profiles.client)
	apiURL, _ := url.Parse(profiles.endpoint + "/GetMyProperties")
	query := apiURL.Query()
	for k, v := range profiles.modifiers {
		query.Set(k, trimMultiline(v))
	}
	apiURL.RawQuery = query.Encode()
	return sp.Post(apiURL.String(), nil, getConfHeaders(profiles.config))
}

// GetPropertiesFor ...
func (profiles *Profiles) GetPropertiesFor(loginName string) (ProfilePropsResp, error) {
	sp := NewHTTPClient(profiles.client)
	apiURL, _ := url.Parse(
		profiles.endpoint +
			"/GetPropertiesFor('" + url.QueryEscape(loginName) + "')",
	)
	query := apiURL.Query()
	for k, v := range profiles.modifiers {
		query.Set(k, trimMultiline(v))
	}
	apiURL.RawQuery = query.Encode()
	return sp.Get(apiURL.String(), getConfHeaders(profiles.config))
}

// GetUserProfilePropertyFor ...
func (profiles *Profiles) GetUserProfilePropertyFor(loginName string, property string) ([]byte, error) {
	sp := NewHTTPClient(profiles.client)
	endpoint := profiles.endpoint +
		"/GetUserProfilePropertyFor(" +
		"accountname='" + url.QueryEscape(loginName) + "'," +
		"propertyname='" + url.QueryEscape(property) + "')"
	return sp.Get(endpoint, getConfHeaders(profiles.config))
}

// GetOwnerUserProfile ...
func (profiles *Profiles) GetOwnerUserProfile() (ProfileResp, error) {
	sp := NewHTTPClient(profiles.client)
	apiURL, _ := url.Parse(
		getPriorEndpoint(profiles.endpoint, "/_api") +
			"/_api/sp.userprofiles.profileloader.getowneruserprofile",
	)
	query := apiURL.Query()
	for k, v := range profiles.modifiers {
		query.Set(k, trimMultiline(v))
	}
	apiURL.RawQuery = query.Encode()
	return sp.Post(apiURL.String(), nil, getConfHeaders(profiles.config))
}

// UserProfile ...
func (profiles *Profiles) UserProfile() (ProfileResp, error) {
	sp := NewHTTPClient(profiles.client)
	apiURL, _ := url.Parse(
		getPriorEndpoint(profiles.endpoint, "/_api") +
			"/_api/sp.userprofiles.profileloader.getprofileloader/GetUserProfile",
	)
	query := apiURL.Query()
	for k, v := range profiles.modifiers {
		query.Set(k, trimMultiline(v))
	}
	apiURL.RawQuery = query.Encode()
	return sp.Post(apiURL.String(), nil, getConfHeaders(profiles.config))
}

// SetSingleValueProfileProperty ...
func (profiles *Profiles) SetSingleValueProfileProperty(loginName string, property string, value string) ([]byte, error) {
	sp := NewHTTPClient(profiles.client)
	endpoint := profiles.endpoint + "/SetSingleValueProfileProperty"
	prop := map[string]string{}
	prop["accountName"] = loginName
	prop["propertyName"] = property
	prop["propertyValue"] = value
	body, _ := json.Marshal(prop)
	return sp.Post(endpoint, body, getConfHeaders(profiles.config))
}

// SetMultiValuedProfileProperty ...
func (profiles *Profiles) SetMultiValuedProfileProperty(loginName string, property string, values []string) ([]byte, error) {
	sp := NewHTTPClient(profiles.client)
	endpoint := profiles.endpoint + "/SetMultiValuedProfileProperty"
	prop := map[string]interface{}{}
	prop["accountName"] = loginName
	prop["propertyName"] = property
	prop["propertyValues"] = values
	body, _ := json.Marshal(prop)
	return sp.Post(endpoint, body, getConfHeaders(profiles.config))
}

// HideSuggestion ...
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
