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
	ChangeTokenStart      string // Gets or sets a value that specifies the start date and start time for changes that are returned through the query
	ChangeTokenEnd        string // Gets or sets a value that specifies the end date and end time for changes that are returned through the query
	Add                   bool   // Gets or sets a value that specifies whether add changes are included in the query
	Alert                 bool   // Gets or sets a value that specifies whether changes to alerts are included in the query
	ContentType           bool   // Gets or sets a value that specifies whether changes to content types are included in the query
	DeleteObject          bool   // Gets or sets a value that specifies whether deleted objects are included in the query
	Field                 bool   // Gets or sets a value that specifies whether changes to fields are included in the query
	File                  bool   // Gets or sets a value that specifies whether changes to files are included in the query
	Folder                bool   // Gets or sets value that specifies whether changes to folders are included in the query
	Group                 bool   // Gets or sets a value that specifies whether changes to groups are included in the query
	GroupMembershipAdd    bool   // Gets or sets a value that specifies whether adding users to groups is included in the query
	GroupMembershipDelete bool   // Gets or sets a value that specifies whether deleting users from the groups is included in the query
	Item                  bool   // Gets or sets a value that specifies whether general changes to list items are included in the query
	List                  bool   // Gets or sets a value that specifies whether changes to lists are included in the query
	Move                  bool   // Gets or sets a value that specifies whether move changes are included in the query
	Navigation            bool   // Gets or sets a value that specifies whether changes to the navigation structure of a site collection are included in the query
	Rename                bool   // Gets or sets a value that specifies whether renaming changes are included in the query
	Restore               bool   // Gets or sets a value that specifies whether restoring items from the recycle bin or from backups is included in the query
	RoleAssignmentAdd     bool   // Gets or sets a value that specifies whether adding role assignments is included in the query
	RoleAssignmentDelete  bool   // Gets or sets a value that specifies whether adding role assignments is included in the query
	RoleDefinitionAdd     bool   // Gets or sets a value that specifies whether adding role assignments is included in the query
	RoleDefinitionDelete  bool   // Gets or sets a value that specifies whether adding role assignments is included in the query
	RoleDefinitionUpdate  bool   // Gets or sets a value that specifies whether adding role assignments is included in the query
	SecurityPolicy        bool   // Gets or sets a value that specifies whether modifications to security policies are included in the query
	Site                  bool   // Gets or sets a value that specifies whether changes to site collections are included in the query
	SystemUpdate          bool   // Gets or sets a value that specifies whether updates made using the item SystemUpdate method are included in the query
	Update                bool   // Gets or sets a value that specifies whether update changes are included in the query
	User                  bool   // Gets or sets a value that specifies whether changes to users are included in the query
	View                  bool   // Gets or sets a value that specifies whether changes to views are included in the query
	Web                   bool   // Gets or sets a value that specifies whether changes to Web sites are included in the query
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

// GetChanges gets changes in scope of the parent container using provided change query
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
