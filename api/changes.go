package api

import (
	"encoding/json"
	"time"

	"github.com/koltyakov/gosip"
)

// Changes represent SharePoint Changes API queryable collection struct
// Always use NewChanges constructor instead of &Changes{}
type Changes struct {
	client   *gosip.SPClient
	config   *RequestConfig
	endpoint string
}

// ChangeInfo - changes API response payload structure
type ChangeInfo struct {
	ChangeToken       *StringValue `json:"ChangeToken"`
	ChangeType        int          `json:"ChangeType"`
	Editor            string       `json:"Editor"`
	EditorEmailHint   string       `json:"EditorEmailHint"`
	ItemID            int          `json:"ItemId"`
	ListID            string       `json:"ListId"`
	ServerRelativeURL string       `json:"ServerRelativeUrl"`
	SharedByUser      string       `json:"SharedByUser"`
	SharedWithUsers   string       `json:"SharedWithUsers"`
	SiteID            string       `json:"SiteId"`
	Time              time.Time    `json:"Time"`
	UniqueID          string       `json:"UniqueId"`
	WebID             string       `json:"WebId"`
}

// ChangeQuery ...
type ChangeQuery struct {
	ChangeTokenStart      string
	ChangeTokenEnd        string
	Add                   bool
	Alert                 bool
	ContentType           bool
	DeleteObject          bool
	Field                 bool
	File                  bool
	Folder                bool
	Group                 bool
	GroupMembershipAdd    bool
	GroupMembershipDelete bool
	Item                  bool
	List                  bool
	Move                  bool
	Navigation            bool
	Rename                bool
	Restore               bool
	RoleAssignmentAdd     bool
	RoleAssignmentDelete  bool
	RoleDefinitionAdd     bool
	RoleDefinitionDelete  bool
	RoleDefinitionUpdate  bool
	SecurityPolicy        bool
	Site                  bool
	SystemUpdate          bool
	Update                bool
	User                  bool
	View                  bool
	Web                   bool
}

// NewChanges - Changes struct constructor function
func NewChanges(client *gosip.SPClient, endpoint string, config *RequestConfig) *Changes {
	return &Changes{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// ToURL gets endpoint with modificators raw URL gets endpoint with modificators raw URL
func (changes *Changes) ToURL() string {
	return changes.endpoint
}

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (changes *Changes) Conf(config *RequestConfig) *Changes {
	changes.config = config
	return changes
}

// GetChanges ...
func (changes *Changes) GetChanges(changeQuery *ChangeQuery) ([]*ChangeInfo, error) {
	sp := NewHTTPClient(changes.client)
	metadata := map[string]interface{}{}
	if changeQuery != nil {
		optsRaw, _ := json.Marshal(changeQuery)
		json.Unmarshal(optsRaw, &metadata)
	}
	metadata["__metadata"] = map[string]string{"type": "SP.ChangeQuery"}
	if changeQuery.ChangeTokenStart != "" {
		metadata["ChangeTokenStart"] = map[string]string{"StringValue": changeQuery.ChangeTokenStart}
	}
	if changeQuery.ChangeTokenEnd != "" {
		metadata["ChangeTokenEnd"] = map[string]string{"StringValue": changeQuery.ChangeTokenEnd}
	}
	for k, v := range metadata {
		if v == false || v == "" || v == nil {
			delete(metadata, k)
		}
	}
	query := map[string]interface{}{"query": metadata}
	body, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}
	data, err := sp.Post(changes.ToURL(), body, getConfHeaders(changes.config))
	if err != nil {
		return nil, err
	}
	collection, _ := parseODataCollection(data)
	results := []*ChangeInfo{}
	for _, changeItem := range collection {
		c := &ChangeInfo{}
		if err := json.Unmarshal(changeItem, &c); err == nil {
			results = append(results, c)
		}
	}
	return results, nil
}

// ToDo:
// Pagination
