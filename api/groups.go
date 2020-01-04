package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/koltyakov/gosip"
)

// Groups represent SharePoint Site Groups API queryable collection struct
// Always use NewGroups constructor instead of &Groups{}
type Groups struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers *ODataMods
}

// GroupsResp - groups response type with helper processor methods
type GroupsResp []byte

// NewGroups - Groups struct constructor function
func NewGroups(client *gosip.SPClient, endpoint string, config *RequestConfig) *Groups {
	return &Groups{
		client:    client,
		endpoint:  endpoint,
		config:    config,
		modifiers: NewODataMods(),
	}
}

// ToURL gets endpoint with modificators raw URL
func (groups *Groups) ToURL() string {
	return toURL(groups.endpoint, groups.modifiers)
}

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (groups *Groups) Conf(config *RequestConfig) *Groups {
	groups.config = config
	return groups
}

// Select adds $select OData modifier
func (groups *Groups) Select(oDataSelect string) *Groups {
	groups.modifiers.AddSelect(oDataSelect)
	return groups
}

// Expand adds $expand OData modifier
func (groups *Groups) Expand(oDataExpand string) *Groups {
	groups.modifiers.AddExpand(oDataExpand)
	return groups
}

// Filter adds $filter OData modifier
func (groups *Groups) Filter(oDataFilter string) *Groups {
	groups.modifiers.AddFilter(oDataFilter)
	return groups
}

// Top adds $top OData modifier
func (groups *Groups) Top(oDataTop int) *Groups {
	groups.modifiers.AddTop(oDataTop)
	return groups
}

// OrderBy adds $orderby OData modifier
func (groups *Groups) OrderBy(oDataOrderBy string, ascending bool) *Groups {
	groups.modifiers.AddOrderBy(oDataOrderBy, ascending)
	return groups
}

// Get gets site Groups response - a collection of GroupInfo
func (groups *Groups) Get() (GroupsResp, error) {
	sp := NewHTTPClient(groups.client)
	return sp.Get(groups.ToURL(), getConfHeaders(groups.config))
}

// Add creates new group with a specified name. Additional metadata can optionally be provided as string map object.
// `metadata` should correspond to SP.Group type.
func (groups *Groups) Add(title string, metadata map[string]interface{}) (GroupResp, error) {
	if metadata == nil {
		metadata = make(map[string]interface{})
	}
	metadata["__metadata"] = map[string]string{
		"type": "SP.Group",
	}
	metadata["Title"] = title
	body, _ := json.Marshal(metadata)
	sp := NewHTTPClient(groups.client)
	return sp.Post(groups.endpoint, body, getConfHeaders(groups.config))
}

// GetByID gets a group object by its ID
func (groups *Groups) GetByID(groupID int) *Group {
	return NewGroup(
		groups.client,
		fmt.Sprintf("%s/GetById(%d)", groups.endpoint, groupID),
		groups.config,
	)
}

// GetByName gets a group object by its name
func (groups *Groups) GetByName(groupName string) *Group {
	return NewGroup(
		groups.client,
		fmt.Sprintf("%s/GetByName('%s')", groups.endpoint, groupName),
		groups.config,
	)
}

// RemoveByID deletes a group object by its ID
func (groups *Groups) RemoveByID(groupID int) ([]byte, error) {
	endpoint := fmt.Sprintf("%s/RemoveById(%d)", groups.ToURL(), groupID)
	sp := NewHTTPClient(groups.client)
	return sp.Post(endpoint, nil, getConfHeaders(groups.config))
}

// RemoveByLoginName deletes a group object by its login name
func (groups *Groups) RemoveByLoginName(loginName string) ([]byte, error) {
	endpoint := fmt.Sprintf(
		"%s/RemoveByLoginName('%s')",
		groups.endpoint,
		url.QueryEscape(loginName),
	)
	sp := NewHTTPClient(groups.client)
	return sp.Post(endpoint, nil, getConfHeaders(groups.config))
}

/* Response helpers */

// Data : to get typed data
func (groupsResp *GroupsResp) Data() []GroupResp {
	collection, _ := parseODataCollection(*groupsResp)
	cts := []GroupResp{}
	for _, ct := range collection {
		cts = append(cts, GroupResp(ct))
	}
	return cts
}

// Unmarshal : to unmarshal to custom object
func (groupsResp *GroupsResp) Unmarshal(obj interface{}) error {
	// collection := parseODataCollection(*groupsResp)
	// data, _ := json.Marshal(collection)
	data, _ := parseODataCollectionPlain(*groupsResp)
	return json.Unmarshal(data, obj)
}
