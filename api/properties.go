package api

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/koltyakov/gosip"
)

//go:generate ggen -ent Properties -conf -mods Select,Expand

// Properties represent SharePoint Property Bags API queryable collection struct
// Always use NewProperties constructor instead of &Properties{}
type Properties struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers *ODataMods
	entity    string
}

// PropsResp - property bags response type with helper processor methods
type PropsResp []byte

// NewProperties - Properties struct constructor function
func NewProperties(client *gosip.SPClient, endpoint string, config *RequestConfig, entity string) *Properties {
	return &Properties{
		client:    client,
		endpoint:  endpoint,
		config:    config,
		modifiers: NewODataMods(),
		entity:    entity,
	}
}

// ToURL gets endpoint with modificators raw URL
func (properties *Properties) ToURL() string {
	return toURL(properties.endpoint, properties.modifiers)
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
	scoped := NewProperties(properties.client, properties.endpoint, properties.config, properties.entity)
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
	if properties.entity == "web" {
		return properties.setWebProps(props)
	}
	if properties.entity == "folder" {
		return properties.setFolderProps(props)
	}
	if properties.entity == "file" {
		return properties.setFileProps(props)
	}
	return fmt.Errorf("unknown parent entity %s", properties.entity)
}

// setWebProps sets multiple properties defined in string map object (CSOM helper)
func (properties *Properties) setWebProps(props map[string]string) error {
	// ToDo: exclude extra call to site and web metadata
	web := NewWeb(properties.client, getPriorEndpoint(properties.endpoint, "/AllProperties"), properties.config)

	site := NewSP(properties.client).Site()
	siteR, err := site.Select("Id").Get()
	if err != nil {
		return err
	}

	webR, err := web.Select("Id").Get()
	if err != nil {
		return err
	}

	identity := fmt.Sprintf("740c6a0b-85e2-48a0-a494-e0f1759d4aa7:site:%s:web:%s", siteR.Data().ID, webR.Data().ID)

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
				<Property Id="4" ParentId="2" Name="AllProperties" />
			</ObjectPaths>
		</Request>
	`))
	_, err = sp.ProcessQuery(properties.client.AuthCnfg.GetSiteURL(), body)

	printNoScriptWarning(properties.endpoint, err)
	return err
}

// setFolderProps sets multiple properties defined in string map object (CSOM helper)
func (properties *Properties) setFolderProps(props map[string]string) error {
	// ToDo: exclude extra call to site, web and folder metadata
	identity := ""

	web := NewWeb(properties.client, getIncludeEndpoint(properties.endpoint, "/Web"), properties.config)
	folder := NewFolder(properties.client, getPriorEndpoint(properties.endpoint, "/Properties"), properties.config)

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

	folderR, err := folder.Select("UniqueId").Get()
	if err != nil {
		return err
	}
	identity = fmt.Sprintf("7394289f-308a-9000-9495-3d03f105ec57|%s:folder:%s", identity, folderR.Data().UniqueID)

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
				<Property Id="4" ParentId="2" Name="Properties" />
			</ObjectPaths>
		</Request>
	`))
	_, err = sp.ProcessQuery(properties.client.AuthCnfg.GetSiteURL(), body)

	printNoScriptWarning(properties.endpoint, err)
	return err
}

// setFileProps sets multiple properties defined in string map object (CSOM helper)
func (properties *Properties) setFileProps(props map[string]string) error {
	file := NewFile(properties.client, getPriorEndpoint(properties.endpoint, "/Properties"), properties.config)

	fileR, err := file.Select("UniqueId").Get()
	if err != nil {
		return err
	}

	methods := ""
	csomIndex := 24
	for key, val := range props {
		methods += TrimMultiline(`
			<Method Name="SetFieldValue" Id="` + strconv.Itoa(csomIndex) + `" ObjectPathId="23">
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
				<Method Name="Update" Id="` + strconv.Itoa(csomIndex) + `" ObjectPathId="21" />
			</Actions>
			<ObjectPaths>
				<StaticProperty Id="7" TypeId="{3747adcd-a3c3-41b9-bfab-4a64dd2f1e0a}" Name="Current" />
				<Property Id="9" ParentId="7" Name="Web" />
				<Method Id="21" ParentId="9" Name="GetFileById">
					<Parameters>
						<Parameter Type="String">` + fileR.Data().UniqueID + `</Parameter>
					</Parameters>
				</Method>
				<Property Id="23" ParentId="21" Name="Properties" />
			</ObjectPaths>
		</Request>
	`))
	_, err = sp.ProcessQuery(properties.client.AuthCnfg.GetSiteURL(), body)

	printNoScriptWarning(properties.endpoint, err)
	return err
}

func printNoScriptWarning(endpoint string, err error) {
	if err != nil && strings.Contains(err.Error(), "System.UnauthorizedAccessException") {
		siteURL := getPriorEndpoint(endpoint, "/_api")
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

// Normalized returns normalized body
func (propsResp *PropsResp) Normalized() []byte {
	return NormalizeODataItem(*propsResp)
}
