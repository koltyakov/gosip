package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/koltyakov/gosip"
)

// List ...
type List struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers map[string]string
}

// ListInfo ...
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

// ListResp ...
type ListResp []byte

// NewList ...
func NewList(client *gosip.SPClient, endpoint string, config *RequestConfig) *List {
	return &List{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// ToURL ...
func (list *List) ToURL() string {
	apiURL, _ := url.Parse(list.endpoint)
	query := apiURL.Query() // url.Values{}
	for k, v := range list.modifiers {
		query.Set(k, trimMultiline(v))
	}
	apiURL.RawQuery = query.Encode()
	return apiURL.String()
}

// Conf ...
func (list *List) Conf(config *RequestConfig) *List {
	list.config = config
	return list
}

// Select ...
func (list *List) Select(oDataSelect string) *List {
	if list.modifiers == nil {
		list.modifiers = make(map[string]string)
	}
	list.modifiers["$select"] = oDataSelect
	return list
}

// Expand ...
func (list *List) Expand(oDataExpand string) *List {
	if list.modifiers == nil {
		list.modifiers = make(map[string]string)
	}
	list.modifiers["$expand"] = oDataExpand
	return list
}

// Get ...
func (list *List) Get() (ListResp, error) {
	sp := NewHTTPClient(list.client)
	return sp.Get(list.ToURL(), getConfHeaders(list.config))
}

// Delete ...
func (list *List) Delete() ([]byte, error) {
	sp := NewHTTPClient(list.client)
	return sp.Delete(list.endpoint, getConfHeaders(list.config))
}

// Update ...
func (list *List) Update(body []byte) ([]byte, error) {
	sp := NewHTTPClient(list.client)
	return sp.Update(list.endpoint, body, getConfHeaders(list.config))
}

// Items ...
func (list *List) Items() *Items {
	return NewItems(
		list.client,
		fmt.Sprintf("%s/Items", list.endpoint),
		list.config,
	)
}

// ContentTypes ...
func (list *List) ContentTypes() *ContentTypes {
	return NewContentTypes(
		list.client,
		fmt.Sprintf("%s/ContentTypes", list.endpoint),
		list.config,
	)
}

// GetChangeToken ...
func (list *List) GetChangeToken() (string, error) {
	scoped := *list
	data, err := scoped.Select("CurrentChangeToken").Get()
	if err != nil {
		return "", err
	}
	return data.Data().CurrentChangeToken.StringValue, nil
}

// GetChanges ...
func (list *List) GetChanges(changeQuery *ChangeQuery) ([]*ChangeInfo, error) {
	return NewChanges(
		list.client,
		fmt.Sprintf("%s/GetChanges", list.endpoint),
		list.config,
	).GetChanges(changeQuery)
}

// Fields ...
func (list *List) Fields() *Fields {
	return NewFields(
		list.client,
		fmt.Sprintf("%s/Fields", list.endpoint),
		list.config,
	)
}

// Views ...
func (list *List) Views() *Views {
	return NewViews(
		list.client,
		fmt.Sprintf("%s/Views", list.endpoint),
		list.config,
	)
}

// ParentWeb ...
func (list *List) ParentWeb() *Web {
	return NewWeb(
		list.client,
		fmt.Sprintf("%s/ParentWeb", list.endpoint),
		list.config,
	)
}

// GetEntityType ...
func (list *List) GetEntityType() (string, error) {
	scoped := *list
	data, err := scoped.Select("ListItemEntityTypeFullName").Get()
	if err != nil {
		return "", err
	}
	return data.Data().ListItemEntityTypeFullName, nil
}

// Roles ...
func (list *List) Roles() *Roles {
	return NewRoles(list.client, list.endpoint, list.config)
}

// ToDo:
// Fields
// Content Type
// Views

/* Response helpers */

// Data : to get typed data
func (listResp *ListResp) Data() *ListInfo {
	data := parseODataItem(*listResp)
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
	data := parseODataItem(*listResp)
	return json.Unmarshal(data, obj)
}
