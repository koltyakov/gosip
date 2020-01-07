package api

import (
	"encoding/json"
	"fmt"

	"github.com/koltyakov/gosip"
)

// CustomActions represent SharePoint CustomActions API queryable collection struct
// Always use NewCustomActions constructor instead of &CustomActions{}
type CustomActions struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers *ODataMods
}

// CustomActionInfo - user custom actions API response payload structure
type CustomActionInfo struct {
	ClientSideComponentID         string           `json:"ClientSideComponentId"`
	ClientSideComponentProperties string           `json:"ClientSideComponentProperties"`
	CommandUIExtension            string           `json:"CommandUIExtension"`
	Description                   string           `json:"Description"`
	Group                         string           `json:"Group"`
	HostProperties                string           `json:"HostProperties"`
	ID                            string           `json:"Id"`
	ImageURL                      string           `json:"ImageUrl"`
	Location                      string           `json:"Location"`
	Name                          string           `json:"Name"`
	RegistrationID                string           `json:"RegistrationId"`
	RegistrationType              int              `json:"RegistrationType"`
	Scope                         int              `json:"Scope"`
	ScriptBlock                   string           `json:"ScriptBlock"`
	ScriptSrc                     string           `json:"ScriptSrc"`
	Sequence                      int              `json:"Sequence"`
	Title                         string           `json:"Title"`
	URL                           string           `json:"Url"`
	VersionOfUserCustomAction     string           `json:"VersionOfUserCustomAction"`
	Rights                        *BasePermissions `json:"Rights"`
}

// NewCustomActions - CustomActions struct constructor function
func NewCustomActions(client *gosip.SPClient, endpoint string, config *RequestConfig) *CustomActions {
	return &CustomActions{
		client:    client,
		endpoint:  endpoint,
		config:    config,
		modifiers: NewODataMods(),
	}
}

// ToURL gets endpoint with modificators raw URL
func (customActions *CustomActions) ToURL() string {
	return toURL(customActions.endpoint, customActions.modifiers)
}

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (customActions *CustomActions) Conf(config *RequestConfig) *CustomActions {
	customActions.config = config
	return customActions
}

// Select adds $select OData modifier
func (customActions *CustomActions) Select(oDataSelect string) *CustomActions {
	customActions.modifiers.AddSelect(oDataSelect)
	return customActions
}

// Filter adds $filter OData modifier
func (customActions *CustomActions) Filter(oDataFilter string) *CustomActions {
	customActions.modifiers.AddFilter(oDataFilter)
	return customActions
}

// Top adds $top OData modifier
func (customActions *CustomActions) Top(oDataTop int) *CustomActions {
	customActions.modifiers.AddTop(oDataTop)
	return customActions
}

// OrderBy adds $orderby OData modifier
func (customActions *CustomActions) OrderBy(oDataOrderBy string, ascending bool) *CustomActions {
	customActions.modifiers.AddOrderBy(oDataOrderBy, ascending)
	return customActions
}

// Get gets event customActions collection
func (customActions *CustomActions) Get() ([]*CustomActionInfo, error) {
	sp := NewHTTPClient(customActions.client)
	data, err := sp.Get(customActions.ToURL(), getConfHeaders(customActions.config))
	if err != nil {
		return nil, err
	}
	data, _ = parseODataCollectionPlain(data)
	res := []*CustomActionInfo{}
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// Add register new user custom action
func (customActions *CustomActions) Add(payload []byte) (*CustomActionInfo, error) {
	body := patchMetadataType(payload, "SP.UserCustomAction")
	sp := NewHTTPClient(customActions.client)
	data, err := sp.Post(customActions.endpoint, body, getConfHeaders(customActions.config))
	if err != nil {
		return nil, err
	}
	data = NormalizeODataItem(data)
	res := &CustomActionInfo{}
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res, nil
}

/* Custom action */

// CustomAction represent SharePoint CustomAction API object
// Always use NewCustomAction constructor instead of &CustomAction{}
type CustomAction struct {
	client   *gosip.SPClient
	config   *RequestConfig
	endpoint string
}

// NewCustomAction - CustomActions struct constructor function
func NewCustomAction(client *gosip.SPClient, endpoint string, config *RequestConfig) *CustomAction {
	return &CustomAction{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// GetByID gets action API onject by ID (GUID)
func (customActions *CustomActions) GetByID(actionID string) *CustomAction {
	return NewCustomAction(
		customActions.client,
		fmt.Sprintf("%s('%s')", customActions.endpoint, actionID),
		customActions.config,
	)
}

// Get gets this action metadata
func (customAction *CustomAction) Get() (*CustomActionInfo, error) {
	sp := NewHTTPClient(customAction.client)
	data, err := sp.Get(customAction.endpoint, getConfHeaders(customAction.config))
	if err != nil {
		return nil, err
	}
	data = NormalizeODataItem(data)
	res := &CustomActionInfo{}
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// Update updates this custom action
func (customAction *CustomAction) Update(payload []byte) (*CustomActionInfo, error) {
	body := patchMetadataType(payload, "SP.UserCustomAction")
	sp := NewHTTPClient(customAction.client)
	data, err := sp.Post(customAction.endpoint, body, getConfHeaders(customAction.config))
	if err != nil {
		return nil, err
	}
	data = NormalizeODataItem(data)
	res := &CustomActionInfo{}
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// Delete deletes this custom action
func (customAction *CustomAction) Delete() error {
	sp := NewHTTPClient(customAction.client)
	_, err := sp.Delete(customAction.endpoint, getConfHeaders(customAction.config))
	return err
}
