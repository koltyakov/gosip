package api

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/koltyakov/gosip"
)

//go:generate ggen -ent ContentTypes -item ContentType -conf -coll -mods Select,Expand,Filter,Top,OrderBy -helpers Data,Normalized

// ContentTypes represent SharePoint Content Types API queryable collection struct
// Always use NewContentTypes constructor instead of &ContentTypes{}
type ContentTypes struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers *ODataMods
}

// ContentTypeCreationInfo new content type metadata
type ContentTypeCreationInfo struct {
	ID                  string // Content type ID, e.g. 0x010000BE397685D43B428513CD9CC1E75CE4, optional is ParentContentTypeID is provided
	Name                string // Content type display name
	Group               string // Content type group name
	Description         string // Description text
	ParentContentTypeID string // Parent content type ID, e.g. 0x01, optional is ID is provided
}

// ContentTypesResp - content types response type with helper processor methods
type ContentTypesResp []byte

// NewContentTypes - ContentTypes struct constructor function
func NewContentTypes(client *gosip.SPClient, endpoint string, config *RequestConfig) *ContentTypes {
	return &ContentTypes{
		client:    client,
		endpoint:  endpoint,
		config:    config,
		modifiers: NewODataMods(),
	}
}

// ToURL gets endpoint with modificators raw URL
func (contentTypes *ContentTypes) ToURL() string {
	return toURL(contentTypes.endpoint, contentTypes.modifiers)
}

// Get gets content typed queryable collection response
func (contentTypes *ContentTypes) Get() (ContentTypesResp, error) {
	sp := NewHTTPClient(contentTypes.client)
	return sp.Get(contentTypes.ToURL(), getConfHeaders(contentTypes.config))
}

// GetByID gets a content type by its ID (GUID)
func (contentTypes *ContentTypes) GetByID(contentTypeID string) *ContentType {
	return NewContentType(
		contentTypes.client,
		fmt.Sprintf("%s('%s')", contentTypes.endpoint, contentTypeID),
		contentTypes.config,
	)
}

// Add adds Content Types with properties provided in `body` parameter
// where `body` is byte array representation of JSON string payload relevalt to SP.ContentType object
func (contentTypes *ContentTypes) Add(body []byte) (ContentTypeResp, error) {
	// REST API doesn't work in that context as supposed to https://github.com/pnp/pnpjs/issues/457
	body = patchMetadataType(body, "SP.ContentType")
	sp := NewHTTPClient(contentTypes.client)
	return sp.Post(contentTypes.endpoint, body, getConfHeaders(contentTypes.config))
}

// Create adds Content Type using CSOM polyfill as REST's Add method is limited (https://github.com/pnp/pnpjs/issues/457)
func (contentTypes *ContentTypes) Create(contentTypeInfo *ContentTypeCreationInfo) (string, error) {
	sp := NewHTTPClient(contentTypes.client)
	objectPaths := TrimMultiline(`
		<StaticProperty Id="0" TypeId="{3747adcd-a3c3-41b9-bfab-4a64dd2f1e0a}" Name="Current" />
		<Property Id="2" ParentId="0" Name="Web" />
		<Property Id="4" ParentId="2" Name="ContentTypes" />
	`)
	nextIndex := 6
	if contentTypeInfo.ParentContentTypeID != "" {
		objectPaths = TrimMultiline(objectPaths + `
			<Method Id="6" ParentId="4" Name="GetById">
				<Parameters>
					<Parameter Type="String">` + contentTypeInfo.ParentContentTypeID + `</Parameter>
				</Parameters>
			</Method>
		`)
		nextIndex = 8
	}
	ctIDProp := `<Property Name="Id" Type="Null" />`
	if contentTypeInfo.ID != "" {
		ctIDProp = `<Property Name="Id" Type="String">` + contentTypeInfo.ID + `</Property>`
	}
	pctIDProp := `<Property Name="ParentContentType" Type="Null" />`
	if contentTypeInfo.ParentContentTypeID != "" {
		pctIDProp = `<Property Name="ParentContentType" ObjectPathId="6" />`
	}
	objectPaths = TrimMultiline(objectPaths + `
		<Method Id="` + strconv.Itoa(nextIndex) + `" ParentId="4" Name="Add">
			<Parameters>
				<Parameter TypeId="{168f3091-4554-4f14-8866-b20d48e45b54}">
					` + ctIDProp + `
					<Property Name="Name" Type="String">` + contentTypeInfo.Name + `</Property>
					<Property Name="Group" Type="String">` + contentTypeInfo.Group + `</Property>
					<Property Name="Description" Type="String">` + contentTypeInfo.Description + `</Property>
					` + pctIDProp + `
				</Parameter>
			</Parameters>
		</Method>
	`)
	body := []byte(TrimMultiline(`
		<Request xmlns="http://schemas.microsoft.com/sharepoint/clientquery/2009" SchemaVersion="15.0.0.0" LibraryVersion="16.0.0.0" ApplicationName="Javascript Library">
			<Actions>
				<ObjectIdentityQuery Id="` + strconv.Itoa(nextIndex+2) + `" ObjectPathId="` + strconv.Itoa(nextIndex) + `" />
			</Actions>
			<ObjectPaths>
				` + objectPaths + `
			</ObjectPaths>
		</Request>
	`))
	resp, err := sp.ProcessQuery(contentTypes.client.AuthCnfg.GetSiteURL(), body)
	if err != nil {
		return "", nil
	}
	rgx := regexp.MustCompile(`:contenttype:(.*?)"`)
	rs := rgx.FindStringSubmatch(fmt.Sprintf("%s", resp))
	return rs[1], nil
}
