package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/koltyakov/gosip"
)

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

// FromURL gets List object using its API URL
func (list *List) FromURL(url string) *List {
	url = strings.Split(url, "?")[0]
	return NewList(list.client, url, list.config)
}

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (list *List) Conf(config *RequestConfig) *List {
	list.config = config
	return list
}

// Select adds $select OData modifier
func (list *List) Select(oDataSelect string) *List {
	list.modifiers.AddSelect(oDataSelect)
	return list
}

// Expand adds $expand OData modifier
func (list *List) Expand(oDataExpand string) *List {
	list.modifiers.AddExpand(oDataExpand)
	return list
}

// Get gets list's data object
func (list *List) Get() (ListResp, error) {
	sp := NewHTTPClient(list.client)
	return sp.Get(list.ToURL(), getConfHeaders(list.config))
}

// Delete deletes a list (can't be restored from a recycle bin)
func (list *List) Delete() error {
	sp := NewHTTPClient(list.client)
	_, err := sp.Delete(list.endpoint, getConfHeaders(list.config))
	return err
}

// Recycle moves this list to the recycle bin
func (list *List) Recycle() error {
	endpoint := fmt.Sprintf("%s/Recycle", list.endpoint)
	sp := NewHTTPClient(list.client)
	_, err := sp.Post(endpoint, nil, getConfHeaders(list.config))
	return err
}

// Update updates Lists's metadata with properties provided in `body` parameter
// where `body` is byte array representation of JSON string payload relevalt to SP.List object
func (list *List) Update(body []byte) (ListResp, error) {
	body = patchMetadataType(body, "SP.List")
	sp := NewHTTPClient(list.client)
	return sp.Update(list.endpoint, body, getConfHeaders(list.config))
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

// RootFolder gets lists's root folder object
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
	sp := NewHTTPClient(list.client)
	endpoint := fmt.Sprintf("%s/ReserveListItemId", list.endpoint)
	data, err := sp.Post(endpoint, nil, getConfHeaders(list.config))
	if err != nil {
		return 0, err
	}
	data = NormalizeODataItem(data)
	if res, err := strconv.Atoi(fmt.Sprintf("%s", data)); err == nil {
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
	sp := NewHTTPClient(list.client)
	apiURL, _ := url.Parse(fmt.Sprintf("%s/RenderListData(@viewXml)", list.endpoint))
	query := apiURL.Query()
	query.Set("@viewXml", `'`+TrimMultiline(viewXML)+`'`)
	apiURL.RawQuery = query.Encode()
	data, err := sp.Post(apiURL.String(), nil, getConfHeaders(list.config))
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
	return []byte(data), nil
}

// RenderListDataAsStream ...
// ToDo

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
	json.Unmarshal(data, &res)
	return res
}

// Unmarshal : to unmarshal to custom object
func (listResp *ListResp) Unmarshal(obj interface{}) error {
	data := NormalizeODataItem(*listResp)
	return json.Unmarshal(data, obj)
}

// Data : to get typed data
func (listData *RenderListDataResp) Data() *RenderListDataInfo {
	// data := parseODataItem(*listData)
	res := &RenderListDataInfo{}
	json.Unmarshal(*listData, &res)
	return res
}
