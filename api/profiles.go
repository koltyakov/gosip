package api

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/koltyakov/gosip"
)

// Profiles ...
type Profiles struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers map[string]string
}

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
func (profiles *Profiles) GetMyProperties() ([]byte, error) {
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
func (profiles *Profiles) GetPropertiesFor(loginName string) ([]byte, error) {
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
func (profiles *Profiles) GetOwnerUserProfile() ([]byte, error) {
	sp := NewHTTPClient(profiles.client)
	apiURL, _ := url.Parse(
		strings.Split(profiles.endpoint, "/_api")[0] +
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
func (profiles *Profiles) UserProfile() ([]byte, error) {
	sp := NewHTTPClient(profiles.client)
	apiURL, _ := url.Parse(
		strings.Split(profiles.endpoint, "/_api")[0] +
			"/_api/sp.userprofiles.profileloader.getprofileloader/getuserprofile",
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
