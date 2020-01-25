package api

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/koltyakov/gosip"
)

//go:generate ggen -ent FieldLinks -conf -mods Select,Filter,Top

// FieldLinks represent SharePoint content type FieldLinks API queryable collection struct
// Always use NewFieldLinks constructor instead of &FieldLinks{}
type FieldLinks struct {
	client        *gosip.SPClient
	config        *RequestConfig
	endpoint      string
	modifiers     *ODataMods
	contentTypeID string
}

// FieldLink represent SharePoint content type FieldLink API
// Always use NewFieldLink constructor instead of &FieldLink{}
type FieldLink struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers *ODataMods
}

// FieldLinkInfo field link info
type FieldLinkInfo struct {
	ID                string `json:"Id"`
	Name              string `json:"Name"`
	FieldInternalName string `json:"FieldInternalName"`
	Hidden            bool   `json:"Hidden"`
	Required          bool   `json:"Required"`
}

// FieldLinksResp - fieldLinks response type with helper processor methods
type FieldLinksResp []byte

// FieldLinkResp - fieldLinks response type with helper processor methods
type FieldLinkResp []byte

// NewFieldLinks - FieldLinks struct constructor function
func NewFieldLinks(client *gosip.SPClient, endpoint string, config *RequestConfig) *FieldLinks {
	return &FieldLinks{
		client:    client,
		endpoint:  endpoint,
		config:    config,
		modifiers: NewODataMods(),
	}
}

// NewFieldLins - FieldLink struct constructor function
func NewFieldLink(client *gosip.SPClient, endpoint string, config *RequestConfig) *FieldLink {
	return &FieldLink{
		client:    client,
		endpoint:  endpoint,
		config:    config,
		modifiers: NewODataMods(),
	}
}

// ToURL gets endpoint with modificators raw URL
func (fieldLinks *FieldLinks) ToURL() string {
	return toURL(fieldLinks.endpoint, fieldLinks.modifiers)
}

// Get gets fieds response collection
func (fieldLinks *FieldLinks) Get() (FieldLinksResp, error) {
	sp := NewHTTPClient(fieldLinks.client)
	return sp.Get(fieldLinks.ToURL(), getConfHeaders(fieldLinks.config))
}

// GetByID gets a field link by its ID (GUID)
func (fieldLinks *FieldLinks) GetByID(fieldLinkID string) *FieldLink {
	return NewFieldLink(
		fieldLinks.client,
		fmt.Sprintf("%s('%s')", fieldLinks.endpoint, fieldLinkID),
		fieldLinks.config,
	)
}

// Delete deletes a field link by its ID (GUID)
func (fieldLink *FieldLink) Delete() error {
	sp := NewHTTPClient(fieldLink.client)
	_, err := sp.Delete(fieldLink.endpoint, getConfHeaders(fieldLink.config))
	return err
}

// // Update updates a field link
// func (fieldLink *FieldLink) Update(body []byte) (FieldLinkResp, error) {
// 	body = patchMetadataType(body, "SP.FieldLink")
// 	sp := NewHTTPClient(fieldLink.client)
// 	return sp.Update(fieldLink.endpoint, body, getConfHeaders(fieldLink.config))
// }

// GetFields gets fieds response collection
func (fieldLinks *FieldLinks) GetFields() (FieldsResp, error) {
	endpoint := getPriorEndpoint(fieldLinks.endpoint, "/FieldLinks")
	fields := NewFields(
		fieldLinks.client,
		endpoint,
		fieldLinks.config,
		"contentType",
	)
	fields.modifiers = fieldLinks.modifiers
	return fields.Get()
}

// Add adds field link
func (fieldLinks *FieldLinks) Add(name string) (string, error) {
	// // REST API doesn't work in that context as supposed to (https://social.msdn.microsoft.com/Forums/office/en-US/52dc9d24-2eb3-4540-a26a-02b12fe1155b/rest-add-fieldlink-to-content-type?forum=sharepointdevelopment)
	// body := []byte(TrimMultiline(fmt.Sprintf(`{
	// 	"__metadata": { "type": "SP.FieldLink" },
	// 	"FieldInternalName": "%s",
	// 	"Hidden": %t,
	// 	"Required": %t
	// }`, name, hidden, required)))
	// sp := NewHTTPClient(fieldLinks.client)
	// resp, err := sp.Post(fieldLinks.endpoint, body, getConfHeaders(fieldLinks.config))
	// if err != nil {
	// 	return nil, err
	// }
	// linkInfo := &FieldLinkInfo{}
	// if err := json.Unmarshal(resp, &linkInfo); err != nil {
	// 	return nil, err
	// }
	// return linkInfo, nil

	if fieldLinks.contentTypeID == "" {
		ct := NewContentType(
			fieldLinks.client,
			getPriorEndpoint(fieldLinks.endpoint, "/FieldLinks"),
			fieldLinks.config,
		)
		resp, err := ct.Select("StringId").Get()
		if err != nil {
			return "", err
		}
		fieldLinks.contentTypeID = resp.Data().ID
		if fieldLinks.contentTypeID == "" {
			return "", fmt.Errorf("can't get content type ID")
		}
	}

	body := []byte(TrimMultiline(`
		<Request xmlns="http://schemas.microsoft.com/sharepoint/clientquery/2009" SchemaVersion="15.0.0.0" LibraryVersion="16.0.0.0" ApplicationName="Javascript Library">
			<Actions>
				<ObjectIdentityQuery Id="16" ObjectPathId="14" />
				<Method Name="Update" Id="17" ObjectPathId="10">
					<Parameters>
						<Parameter Type="Boolean">false</Parameter>
					</Parameters>
				</Method>
			</Actions>
			<ObjectPaths>
				<StaticProperty Id="0" TypeId="{3747adcd-a3c3-41b9-bfab-4a64dd2f1e0a}" Name="Current" />
				<Property Id="2" ParentId="0" Name="Web" />
				<Property Id="4" ParentId="2" Name="Fields" />
				<Method Id="6" ParentId="4" Name="GetByInternalNameOrTitle">
					<Parameters>
						<Parameter Type="String">` + name + `</Parameter>
					</Parameters>
				</Method>
				<Property Id="8" ParentId="2" Name="ContentTypes" />
				<Method Id="10" ParentId="8" Name="GetById">
					<Parameters>
						<Parameter Type="String">` + fieldLinks.contentTypeID + `</Parameter>
					</Parameters>
				</Method>
				<Property Id="12" ParentId="10" Name="FieldLinks" />
				<Method Id="14" ParentId="12" Name="Add">
					<Parameters>
						<Parameter TypeId="{63fb2c92-8f65-4bbb-a658-b6cd294403f4}">
							<Property Name="Field" ObjectPathId="6" />
						</Parameter>
					</Parameters>
				</Method>
			</ObjectPaths>
		</Request>
	`))

	sp := NewHTTPClient(fieldLinks.client)
	resp, err := sp.ProcessQuery(fieldLinks.client.AuthCnfg.GetSiteURL(), body)
	if err != nil {
		return "", err
	}
	rgx := regexp.MustCompile(`:fl:(.*?)"`)
	rs := rgx.FindStringSubmatch(fmt.Sprintf("%s", resp))
	fieldLinkID := rs[1]

	return fieldLinkID, nil
}

/* Response helpers */

// Data : to get typed data
func (fieldLinksResp *FieldLinksResp) Data() []*FieldLinkInfo {
	collection, _ := normalizeODataCollection(*fieldLinksResp)
	resFieldLinks := []*FieldLinkInfo{}
	for _, f := range collection {
		linkInfo := &FieldLinkInfo{}
		json.Unmarshal(f, &linkInfo)
		resFieldLinks = append(resFieldLinks, linkInfo)
	}
	return resFieldLinks
}

// Normalized returns normalized body
func (fieldLinksResp *FieldLinksResp) Normalized() []byte {
	normalized, _ := NormalizeODataCollection(*fieldLinksResp)
	return normalized
}

// Data : to get typed data
func (fieldLinkResp *FieldLinkResp) Data() *FieldLinkInfo {
	data := NormalizeODataItem(*fieldLinkResp)
	res := &FieldLinkInfo{}
	json.Unmarshal(data, &res)
	return res
}

// Normalized returns normalized body
func (fieldLinkResp *FieldLinkResp) Normalized() []byte {
	return NormalizeODataItem(*fieldLinkResp)
}
