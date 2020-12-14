package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/koltyakov/gosip"
)

//go:generate ggen -ent Changes -conf -coll -mods Top

// Changes represent SharePoint Changes API queryable collection struct
// Always use NewChanges constructor instead of &Changes{}
type Changes struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers *ODataMods
}

// ChangesResp - changes response type with helper processor methods
type ChangesResp struct {
	Data        func() []*ChangeInfo         // gets change info collection
	GetNextPage func() (*ChangesResp, error) // gets next page changes collection
	data        []*ChangeInfo
}

// ChangeInfo - changes API response payload structure
type ChangeInfo struct {
	ChangeToken       *StringValue `json:"ChangeToken"`
	ChangeType        int          `json:"ChangeType"` // see more https://docs.microsoft.com/en-us/previous-versions/office/sharepoint-csom/ee543793(v%3Doffice.15)
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
	ChangeTokenStart      string // Specifies the start date and start time for changes that are returned through the query
	ChangeTokenEnd        string // Specifies the end date and end time for changes that are returned through the query
	Add                   bool   // Specifies whether add changes are included in the query
	Alert                 bool   // Specifies whether changes to alerts are included in the query
	ContentType           bool   // Specifies whether changes to content types are included in the query
	DeleteObject          bool   // Specifies whether deleted objects are included in the query
	Field                 bool   // Specifies whether changes to fields are included in the query
	File                  bool   // Specifies whether changes to files are included in the query
	Folder                bool   // Specifies whether changes to folders are included in the query
	Group                 bool   // Specifies whether changes to groups are included in the query
	GroupMembershipAdd    bool   // Specifies whether adding users to groups is included in the query
	GroupMembershipDelete bool   // Specifies whether deleting users from the groups is included in the query
	Item                  bool   // Specifies whether general changes to list items are included in the query
	List                  bool   // Specifies whether changes to lists are included in the query
	Move                  bool   // Specifies whether move changes are included in the query
	Navigation            bool   // Specifies whether changes to the navigation structure of a site collection are included in the query
	Rename                bool   // Specifies whether renaming changes are included in the query
	Restore               bool   // Specifies whether restoring items from the recycle bin or from backups is included in the query
	RoleAssignmentAdd     bool   // Specifies whether adding role assignments is included in the query
	RoleAssignmentDelete  bool   // Specifies whether adding role assignments is included in the query
	RoleDefinitionAdd     bool   // Specifies whether adding role assignments is included in the query
	RoleDefinitionDelete  bool   // Specifies whether adding role assignments is included in the query
	RoleDefinitionUpdate  bool   // Specifies whether adding role assignments is included in the query
	SecurityPolicy        bool   // Specifies whether modifications to security policies are included in the query
	Site                  bool   // Specifies whether changes to site collections are included in the query
	SystemUpdate          bool   // Specifies whether updates made using the item SystemUpdate method are included in the query
	Update                bool   // Specifies whether update changes are included in the query
	User                  bool   // Specifies whether changes to users are included in the query
	View                  bool   // Specifies whether changes to views are included in the query
	Web                   bool   // Specifies whether changes to Web sites are included in the query
}

// NewChanges - Changes struct constructor function
func NewChanges(client *gosip.SPClient, endpoint string, config *RequestConfig) *Changes {
	return &Changes{
		client:    client,
		endpoint:  endpoint,
		config:    config,
		modifiers: NewODataMods(),
	}
}

// GetCurrentToken gets current change token for this parent entity
func (changes *Changes) GetCurrentToken() (string, error) {
	endpoint := fmt.Sprintf("%s?$select=CurrentChangeToken", changes.endpoint)
	client := NewHTTPClient(changes.client)
	data, err := client.Get(endpoint, changes.config)
	if err != nil {
		return "", err
	}
	data = NormalizeODataItem(data)
	res := &struct {
		CurrentChangeToken StringValue `json:"CurrentChangeToken"`
	}{}
	if err := json.Unmarshal(data, &res); err != nil {
		return "", err
	}
	return res.CurrentChangeToken.StringValue, nil
}

// GetChanges gets changes in scope of the parent container using provided change query
func (changes *Changes) GetChanges(changeQuery *ChangeQuery) (*ChangesResp, error) {
	endpoint := toURL(fmt.Sprintf("%s/GetChanges", changes.endpoint), changes.modifiers)
	client := NewHTTPClient(changes.client)
	metadata := map[string]interface{}{}
	if changeQuery != nil {
		optsRaw, _ := json.Marshal(changeQuery)
		_ = json.Unmarshal(optsRaw, &metadata)
	}
	metadata["__metadata"] = map[string]string{"type": "SP.ChangeQuery"}
	if changeQuery != nil && changeQuery.ChangeTokenStart != "" {
		metadata["ChangeTokenStart"] = map[string]string{"StringValue": changeQuery.ChangeTokenStart}
	}
	if changeQuery != nil && changeQuery.ChangeTokenEnd != "" {
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
	data, err := client.Post(endpoint, bytes.NewBuffer(body), changes.config)
	if err != nil {
		return nil, err
	}

	result := &ChangesResp{}
	result.Data = func() []*ChangeInfo {
		result.data = nil
		collection, _ := normalizeODataCollection(data)
		var changesInfo []*ChangeInfo
		for _, changeItem := range collection {
			c := &ChangeInfo{}
			if err := json.Unmarshal(changeItem, &c); err == nil {
				changesInfo = append(changesInfo, c)
			}
		}
		result.data = changesInfo
		return changesInfo
	}
	result.GetNextPage = func() (*ChangesResp, error) {
		if result.data == nil {
			result.Data()
		}
		if len(result.data) == 0 {
			return nil, fmt.Errorf("can't get next page of an empty collection")
		}
		changeQuery.ChangeTokenStart = result.data[len(result.data)-1].ChangeToken.StringValue
		return changes.GetChanges(changeQuery)
	}

	return result, nil
}

// GetChangeType gets verbose change type
// https://docs.microsoft.com/en-us/previous-versions/office/sharepoint-csom/ee543793(v%3Doffice.15)
func (changes *Changes) GetChangeType(changeType int) string {
	changeTypes := map[int]string{
		0:  "NoChange",         // Enumeration whose values indicate that no change has taken place. The value = 0.
		1:  "Add",              // Enumeration whose values specify that an object has been added within the scope of a list, site, site collection, or content database. The value = 1.
		2:  "Update",           // Enumeration whose values specify that an object has been modified within the scope of a list, site, site collection, or content database. The value = 2.
		3:  "DeleteObject",     // Enumeration whose values specify that an object has been deleted within the scope of a list, site, site collection, or content database. The value = 3.
		4:  "Rename",           // Enumeration whose values specify that the leaf in a URL has been renamed. The value = 4.
		5:  "MoveAway",         // Enumeration whose values specify that a non-leaf section within a URL has been renamed. The object was moved away from the location within the site specified by the change. The value = 5.
		6:  "MoveInto",         // Enumeration whose values specify that a non-leaf section within a URL has been renamed. The object was moved into the location within the site specified by the change. The value = 6.
		7:  "Restore",          // Enumeration whose values specify that an object has restored from a backup or from the recycle bin. The value = 7.
		8:  "RoleAdd",          // Enumeration whose values specify that a role definition has been added. The value = 8.
		9:  "RoleDelete",       // Enumeration whose values specify that a role definition has been deleted. The value = 9.
		10: "RoleUpdate",       // Enumeration whose values specify that a role definition has been updated. The value = 10.
		11: "AssignmentAdd",    // Enumeration whose values specify that a user has been given permissions to a list. The value = 11. The list must have unique permissions enabled.
		12: "AssignmentDelete", // Enumeration whose values specify that a user has lost permissions to a list. The value = 12. The list must have unique permissions enabled.
		13: "MemberAdd",        // Enumeration whose values specify that a user has been added to a group. The value = 13.
		14: "MemberDelete",     // Enumeration whose values specify that a user has been removed from a group. The value = 14.
		15: "SystemUpdate",     // Enumeration whose values specify that a change has been made to an item using the SystemUpdate method. The value = 15.
		16: "Navigation",       // Enumeration whose values specify that a change in the navigation structure of a site collection has been made. The value = 16.
		17: "ScopeAdd",         // Enumeration whose values specify that a change in permissions scope has been made to break inheritance from an object’s parent. The value = 17.
		18: "ScopeDelete",      // Enumeration whose values specify that a change in permissions scope has been made to revert back to inheriting permissions from an object’s parent. The value = 18.
	}
	changeName := changeTypes[changeType]
	if len(changeName) == 0 {
		changeName = "Unknown"
	}
	return changeName
}
