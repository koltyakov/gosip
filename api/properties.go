package api

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/koltyakov/gosip"
)

// Properties represent SharePoint Property Bags API queryable collection struct
// Always use NewProperties constructor instead of &Properties{}
type Properties struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers *ODataMods
}

// PropsResp - property bags response type with helper processor methods
type PropsResp []byte

// NewProperties - Properties struct constructor function
func NewProperties(client *gosip.SPClient, endpoint string, config *RequestConfig) *Properties {
	return &Properties{
		client:    client,
		endpoint:  endpoint,
		config:    config,
		modifiers: NewODataMods(),
	}
}

// ToURL gets endpoint with modificators raw URL
func (properties *Properties) ToURL() string {
	return toURL(properties.endpoint, properties.modifiers)
}

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (properties *Properties) Conf(config *RequestConfig) *Properties {
	properties.config = config
	return properties
}

// Select adds $select OData modifier
func (properties *Properties) Select(oDataSelect string) *Properties {
	properties.modifiers.AddSelect(oDataSelect)
	return properties
}

// Expand adds $expand OData modifier
func (properties *Properties) Expand(oDataExpand string) *Properties {
	properties.modifiers.AddExpand(oDataExpand)
	return properties
}

// Get gets properties collection
func (properties *Properties) Get() (PropsResp, error) {
	sp := NewHTTPClient(properties.client)
	return sp.Get(properties.ToURL(), getConfHeaders(properties.config))
}

// GetProps gets specific props values
func (properties *Properties) GetProps(props []string) (map[string]string, error) {
	for indx, prop := range props {
		key := strings.Replace(strings.Replace(prop, "_x005f_", "_", -1), "_", "_x005f_", -1)
		props[indx] = key
	}
	scoped := NewProperties(properties.client, properties.endpoint, properties.config)
	selectProps := ""
	for _, prop := range props {
		if len(selectProps) > 0 {
			selectProps += ","
		}
		selectProps += prop
	}
	res, err := scoped.Select(selectProps).Get()
	if err != nil {
		scoped.modifiers = &ODataMods{}
		res, err := scoped.Get()
		if err != nil {
			return nil, err
		}
		resProps := map[string]string{}
		for key, val := range res.Data() {
			for _, p := range props {
				if p == key {
					resProps[key] = val
				}
				p = strings.Replace(p, "_x005f_", "_", -1)
				if p == key {
					resProps[key] = val
				}
			}
		}
		return resProps, nil
	}
	return res.Data(), nil
}

// Set sets a single property (CSOM helper)
func (properties *Properties) Set(prop string, value string) error {
	return properties.SetProps(map[string]string{prop: value})
}

// SetProps sets multiple properties defined in string map object (CSOM helper)
func (properties *Properties) SetProps(props map[string]string) error {
	var web *Web
	var folder *Folder

	identity := "" // keeps folder of web identity path
	property := "" // takes AllProperties or Properties value based on root object

	// `/AllProperties` endpoint part is from a Web object
	// `/Properties` endpoint part is from a Folder object
	if strings.Contains(strings.ToLower(properties.endpoint), "/allproperties") {
		web = NewWeb(properties.client, getPriorEndpoint(properties.endpoint, "/AllProperties"), properties.config)
		property = "AllProperties"
	} else {
		web = NewWeb(properties.client, getIncludeEndpoint(properties.endpoint, "/Web"), properties.config)
		folder = NewFolder(properties.client, getPriorEndpoint(properties.endpoint, "/Properties"), properties.config)
		property = "Properties"
	}

	site := NewSP(properties.client).Site()
	siteR, err := site.Select("Id").Get()
	if err != nil {
		return err
	}
	identity = fmt.Sprintf("740c6a0b-85e2-48a0-a494-e0f1759d4aa7:site:%s", siteR.Data().ID)

	webR, err := web.Select("Id").Get()
	if err != nil {
		return err
	}
	identity = fmt.Sprintf("%s:web:%s", identity, webR.Data().ID)

	if folder != nil {
		folderR, err := folder.Select("UniqueId").Get()
		if err != nil {
			return err
		}
		identity = fmt.Sprintf("7394289f-308a-9000-9495-3d03f105ec57|%s:folder:%s", identity, folderR.Data().UniqueID)
	}

	methods := ""
	csomIndex := 9
	for key, val := range props {
		methods += TrimMultiline(`
			<Method Name="SetFieldValue" Id="` + strconv.Itoa(csomIndex) + `" ObjectPathId="4">
				<Parameters>
					<Parameter Type="String">` + key + `</Parameter>
					<Parameter Type="String">` + val + `</Parameter>
				</Parameters>
			</Method>
		`)
		csomIndex++
	}
	sp := NewHTTPClient(properties.client)
	body := []byte(TrimMultiline(`
		<Request xmlns="http://schemas.microsoft.com/sharepoint/clientquery/2009" SchemaVersion="15.0.0.0" LibraryVersion="16.0.0.0" ApplicationName="Gosip">
			<Actions>
				` + methods + `
				<Method Name="Update" Id="` + strconv.Itoa(csomIndex) + `" ObjectPathId="2" />
			</Actions>
			<ObjectPaths>
				<Identity Id="2" Name="` + identity + `" />
				<Property Id="4" ParentId="2" Name="` + property + `" />
			</ObjectPaths>
		</Request>
	`))
	_, err = sp.ProcessQuery(properties.client.AuthCnfg.GetSiteURL(), body)

	if err != nil && strings.Contains(err.Error(), "System.UnauthorizedAccessException") {
		siteURL := getPriorEndpoint(properties.endpoint, "/_api")
		if strings.Contains(strings.ToLower(siteURL), ".sharepoint.com") {
			noScriptSiteDisable := fmt.Sprintf("spo site classic set --url %s --noScriptSite false", siteURL)
			err = fmt.Errorf(
				"%s. You probably have \"noScriptSite\" enabled on your site. "+
					"You can enable it using PnP Office 365 CLI by running \"%s\". "+
					"See more: https://pnp.github.io/office365-cli",
				err,
				noScriptSiteDisable,
			)
		}
	}

	return err
}

/* Response helpers */

// Data : to get typed data
func (propsResp *PropsResp) Data() map[string]string {
	data := NormalizeODataItem(*propsResp)
	resAll := map[string]interface{}{}
	json.Unmarshal(data, &resAll)
	res := map[string]string{}
	for key, val := range resAll {
		if reflect.TypeOf(val).String() == "string" {
			key = strings.Replace(key, "_x005f_", "_", -1)
			res[key] = val.(string)
		}
	}
	return res
}
