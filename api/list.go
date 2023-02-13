package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/koltyakov/gosip"
)

//go:generate ggen -ent List -conf -mods Select,Expand -helpers Normalized

// List represents SharePoint List API queryable object struct
// Always use NewList constructor instead of &List{}
type List struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers *ODataMods
}

// ListInfo - list instance response payload structure
type ListInfo struct {
	ID                               string       `json:"Id"`
	Title                            string       `json:"Title"`
	AllowContentTypes                bool         `json:"AllowContentTypes"`
	BaseTemplate                     int          `json:"BaseTemplate"`
	BaseType                         int          `json:"BaseType"`
	ContentTypesEnabled              bool         `json:"ContentTypesEnabled"`
	CrawlNonDefaultViews             bool         `json:"CrawlNonDefaultViews"`
	Created                          time.Time    `json:"Created"`
	DefaultContentApprovalWorkflowID string       `json:"DefaultContentApprovalWorkflowId"`
	DefaultItemOpenUseListSetting    bool         `json:"DefaultItemOpenUseListSetting"`
	Description                      string       `json:"Description"`
	Direction                        string       `json:"Direction"`
	DisableGridEditing               bool         `json:"DisableGridEditing"`
	DocumentTemplateURL              string       `json:"DocumentTemplateUrl"`
	DraftVersionVisibility           int          `json:"DraftVersionVisibility"`
	EnableAttachments                bool         `json:"EnableAttachments"`
	EnableFolderCreation             bool         `json:"EnableFolderCreation"`
	EnableMinorVersions              bool         `json:"EnableMinorVersions"`
	EnableModeration                 bool         `json:"EnableModeration"`
	EnableRequestSignOff             bool         `json:"EnableRequestSignOff"`
	EnableVersioning                 bool         `json:"EnableVersioning"`
	EntityTypeName                   string       `json:"EntityTypeName"`
	FileSavePostProcessingEnabled    bool         `json:"FileSavePostProcessingEnabled"`
	ForceCheckout                    bool         `json:"ForceCheckout"`
	HasExternalDataSource            bool         `json:"HasExternalDataSource"`
	Hidden                           bool         `json:"Hidden"`
	ImageURL                         string       `json:"ImageUrl"`
	IrmEnabled                       bool         `json:"IrmEnabled"`
	IrmExpire                        bool         `json:"IrmExpire"`
	IrmReject                        bool         `json:"IrmReject"`
	IsApplicationList                bool         `json:"IsApplicationList"`
	IsCatalog                        bool         `json:"IsCatalog"`
	IsPrivate                        bool         `json:"IsPrivate"`
	ItemCount                        int          `json:"ItemCount"`
	LastItemDeletedDate              time.Time    `json:"LastItemDeletedDate"`
	LastItemModifiedDate             time.Time    `json:"LastItemModifiedDate"`
	LastItemUserModifiedDate         time.Time    `json:"LastItemUserModifiedDate"`
	ListExperienceOptions            int          `json:"ListExperienceOptions"`
	ListItemEntityTypeFullName       string       `json:"ListItemEntityTypeFullName"`
	MajorVersionLimit                int          `json:"MajorVersionLimit"`
	MajorWithMinorVersionsLimit      int          `json:"MajorWithMinorVersionsLimit"`
	MultipleDataList                 bool         `json:"MultipleDataList"`
	NoCrawl                          bool         `json:"NoCrawl"`
	ParentWebURL                     string       `json:"ParentWebUrl"`
	ParserDisabled                   bool         `json:"ParserDisabled"`
	ServerTemplateCanCreateFolders   bool         `json:"ServerTemplateCanCreateFolders"`
	TemplateFeatureID                string       `json:"TemplateFeatureId"`
	ParentWebPath                    *DecodedURL  `json:"ParentWebPath"`
	ImagePath                        *DecodedURL  `json:"ImagePath"`
	CurrentChangeToken               *StringValue `json:"CurrentChangeToken"`
}

// ListResp - list response type with helper processor methods
type ListResp []byte

// RenderListDataInfo ...
type RenderListDataInfo struct {
	Row                    []map[string]interface{} `json:"Row"`
	FirstRow               int                      `json:"FirstRow"`
	FolderPermissions      string                   `json:"FolderPermissions"`
	LastRow                int                      `json:"LastRow"`
	RowLimit               int                      `json:"RowLimit"`
	FilterLink             string                   `json:"FilterLink"`
	ForceNoHierarchy       string                   `json:"ForceNoHierarchy"`
	HierarchyHasIndention  string                   `json:"HierarchyHasIndention"`
	CurrentFolderSpItemURL string                   `json:"CurrentFolderSpItemUrl"`
}

// RenderListDataResp - renderListData method response type with helper processor methods
type RenderListDataResp []byte

// NewList - List struct constructor function
func NewList(client *gosip.SPClient, endpoint string, config *RequestConfig) *List {
	return &List{
		client:    client,
		endpoint:  endpoint,
		config:    config,
		modifiers: NewODataMods(),
	}
}

// ToURL gets endpoint with modificators raw URL
func (list *List) ToURL() string {
	return toURL(list.endpoint, list.modifiers)
}

// Get gets list's data object
func (list *List) Get() (ListResp, error) {
	client := NewHTTPClient(list.client)
	return client.Get(list.ToURL(), list.config)
}

// Delete deletes a list (can't be restored from a recycle bin)
func (list *List) Delete() error {
	client := NewHTTPClient(list.client)
	_, err := client.Delete(list.endpoint, list.config)
	return err
}

// Recycle moves this list to the recycle bin
func (list *List) Recycle() error {
	endpoint := fmt.Sprintf("%s/Recycle", list.endpoint)
	client := NewHTTPClient(list.client)
	_, err := client.Post(endpoint, nil, list.config)
	return err
}

// Update updates List's metadata with properties provided in `body` parameter
// where `body` is byte array representation of JSON string payload relevant to SP.List object
func (list *List) Update(body []byte) (ListResp, error) {
	body = patchMetadataType(body, "SP.List")
	client := NewHTTPClient(list.client)
	return client.Update(list.endpoint, bytes.NewBuffer(body), list.config)
}

// Items gets Items API instance queryable collection
func (list *List) Items() *Items {
	return NewItems(
		list.client,
		fmt.Sprintf("%s/Items", list.endpoint),
		list.config,
	)
}

// ContentTypes gets list's Content Types API instance queryable collection
func (list *List) ContentTypes() *ContentTypes {
	return NewContentTypes(
		list.client,
		fmt.Sprintf("%s/ContentTypes", list.endpoint),
		list.config,
	)
}

// Subscriptions gets list's subscriptions API instance queryable collection
func (list *List) Subscriptions() *Subscriptions {
	return NewSubscriptions(
		list.client,
		fmt.Sprintf("%s/Subscriptions", list.endpoint),
		list.config,
	)
}

// Changes gets changes API scoped object
func (list *List) Changes() *Changes {
	return NewChanges(
		list.client,
		list.endpoint,
		list.config,
	)
}

// Fields gets list's Fields API instance queryable collection
func (list *List) Fields() *Fields {
	return NewFields(
		list.client,
		fmt.Sprintf("%s/Fields", list.endpoint),
		list.config,
		"list",
	)
}

// Views gets list's Views API instance queryable collection
func (list *List) Views() *Views {
	return NewViews(
		list.client,
		fmt.Sprintf("%s/Views", list.endpoint),
		list.config,
	)
}

// ParentWeb gets list's Parent Web API instance object
func (list *List) ParentWeb() *Web {
	return NewWeb(
		list.client,
		fmt.Sprintf("%s/ParentWeb", list.endpoint),
		list.config,
	)
}

// RootFolder gets list's root folder object
func (list *List) RootFolder() *Folder {
	return NewFolder(
		list.client,
		fmt.Sprintf("%s/RootFolder", list.endpoint),
		list.config,
	)
}

// GetEntityType gets list's ListItemEntityTypeFullName
func (list *List) GetEntityType() (string, error) {
	scoped := NewList(list.client, list.endpoint, list.config)
	data, err := scoped.Select("ListItemEntityTypeFullName").Get()
	if err != nil {
		return "", err
	}
	return data.Data().ListItemEntityTypeFullName, nil
}

// ReserveListItemID reserves item's ID in this list
func (list *List) ReserveListItemID() (int, error) {
	client := NewHTTPClient(list.client)
	endpoint := fmt.Sprintf("%s/ReserveListItemId", list.endpoint)
	data, err := client.Post(endpoint, nil, list.config)
	if err != nil {
		return 0, err
	}
	data = NormalizeODataItem(data)
	if res, err := strconv.Atoi(string(data)); err == nil {
		return res, nil
	}
	res := &struct {
		ReserveListItemID int `json:"ReserveListItemId"`
	}{}
	if err := json.Unmarshal(data, &res); err != nil {
		return 0, err
	}
	return res.ReserveListItemID, nil
}

// RenderListData renders lists content using CAML
func (list *List) RenderListData(viewXML string) (RenderListDataResp, error) {
	client := NewHTTPClient(list.client)
	apiURL, _ := url.Parse(fmt.Sprintf("%s/RenderListData(@viewXml)", list.endpoint))
	query := apiURL.Query()
	query.Set("@viewXml", `'`+TrimMultiline(viewXML)+`'`)
	apiURL.RawQuery = query.Encode()
	data, err := client.Post(apiURL.String(), nil, list.config)
	if err != nil {
		return nil, err
	}
	data = NormalizeODataItem(data)
	res := &struct {
		Value          string `json:"value"`
		RenderListData string `json:"RenderListData"`
	}{}
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	if res.RenderListData != "" {
		data = []byte(res.RenderListData)
	}
	if res.Value != "" {
		data = []byte(res.Value)
	}
	return data, nil
}

// ToDo:
// RenderListDataAsStream ...

// Roles gets list's Roles API instance queryable collection
func (list *List) Roles() *Roles {
	return NewRoles(list.client, list.endpoint, list.config)
}

// ContextInfo gets context info for a web of current list
func (list *List) ContextInfo() (*ContextInfo, error) {
	return NewContext(list.client, list.ToURL(), list.config).Get()
}

/* Response helpers */

// Data : to get typed data
func (listResp *ListResp) Data() *ListInfo {
	data := NormalizeODataItem(*listResp)
	data = fixDatesInResponse(data, []string{
		"Created",
		"LastItemDeletedDate",
		"LastItemModifiedDate",
		"LastItemUserModifiedDate",
	})
	res := &ListInfo{}
	_ = json.Unmarshal(data, &res)
	return res
}

// // Normalized returns normalized body
// func (listResp *ListResp) Normalized() []byte {
// 	return NormalizeODataItem(*listResp)
// }

// Data : to get typed data
func (listData *RenderListDataResp) Data() *RenderListDataInfo {
	res := &RenderListDataInfo{}
	_ = json.Unmarshal(*listData, &res)
	return res
}
