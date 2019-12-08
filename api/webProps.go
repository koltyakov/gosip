package api

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/koltyakov/gosip"
)

// WebProps ...
type WebProps struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers map[string]string
}

// NewWebProps ...
func NewWebProps(client *gosip.SPClient, endpoint string, config *RequestConfig) *WebProps {
	return &WebProps{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// ToURL ...
func (webProps *WebProps) ToURL() string {
	apiURL, _ := url.Parse(webProps.endpoint)
	query := apiURL.Query() // url.Values{}
	for k, v := range webProps.modifiers {
		query.Set(k, trimMultiline(v))
	}
	apiURL.RawQuery = query.Encode()
	return apiURL.String()
}

// Conf ...
func (webProps *WebProps) Conf(config *RequestConfig) *WebProps {
	webProps.config = config
	return webProps
}

// Select ...
func (webProps *WebProps) Select(oDataSelect string) *WebProps {
	if webProps.modifiers == nil {
		webProps.modifiers = make(map[string]string)
	}
	webProps.modifiers["$select"] = oDataSelect
	return webProps
}

// Expand ...
func (webProps *WebProps) Expand(oDataExpand string) *WebProps {
	if webProps.modifiers == nil {
		webProps.modifiers = make(map[string]string)
	}
	webProps.modifiers["$expand"] = oDataExpand
	return webProps
}

// Get ...
func (webProps *WebProps) Get() ([]byte, error) {
	sp := NewHTTPClient(webProps.client)
	headers := map[string]string{}
	if webProps.config != nil {
		headers = webProps.config.Headers
	}
	return sp.Get(webProps.ToURL(), headers)
}

// Set ...
func (webProps *WebProps) Set(prop string, value string) ([]byte, error) {
	// sp := NewHTTPClient(webProps.client)
	// body := []byte(trimMultiline(`
	// 	<Request xmlns="http://schemas.microsoft.com/sharepoint/clientquery/2009" SchemaVersion="15.0.0.0" LibraryVersion="16.0.0.0" ApplicationName="Gosip">
	// 		<Actions>
	// 			<Method Name="SetFieldValue" Id="9" ObjectPathId="4">
	// 				<Parameters>
	// 					<Parameter Type="String">` + prop + `</Parameter>
	// 					<Parameter Type="String">` + value + `</Parameter>
	// 				</Parameters>
	// 			</Method>
	// 			<Method Name="Update" Id="10" ObjectPathId="2" />
	// 		</Actions>
	// 		<ObjectPaths>
	// 			<StaticProperty Id="0" TypeId="{3747adcd-a3c3-41b9-bfab-4a64dd2f1e0a}" Name="Current" />
	// 			<Property Id="2" ParentId="0" Name="Web" />
	// 			<Property Id="4" ParentId="2" Name="AllProperties" />
	// 		</ObjectPaths>
	// 	</Request>
	// `))
	// return sp.ProcessQuery(webProps.endpoint, body)
	return webProps.SetProps(map[string]string{prop: value})
}

// SetProps ...
func (webProps *WebProps) SetProps(props map[string]string) ([]byte, error) {
	site := NewSP(webProps.client).Site()
	web := NewWeb(webProps.client, strings.Split(webProps.endpoint, "/AllProperties")[0], webProps.config)
	siteR, err := site.Select("Id").Get()
	if err != nil {
		return nil, err
	}
	webR, err := web.Select("Id").Get()
	if err != nil {
		return nil, err
	}
	methods := ""
	csomIndex := 9
	for key, val := range props {
		methods += trimMultiline(`
			<Method Name="SetFieldValue" Id="` + strconv.Itoa(csomIndex) + `" ObjectPathId="4">
				<Parameters>
					<Parameter Type="String">` + key + `</Parameter>
					<Parameter Type="String">` + val + `</Parameter>
				</Parameters>
			</Method>
		`)
		csomIndex++
	}
	sp := NewHTTPClient(webProps.client)
	body := []byte(trimMultiline(`
		<Request xmlns="http://schemas.microsoft.com/sharepoint/clientquery/2009" SchemaVersion="15.0.0.0" LibraryVersion="16.0.0.0" ApplicationName="Gosip">
			<Actions>
				` + methods + `
				<Method Name="Update" Id="` + strconv.Itoa(csomIndex) + `" ObjectPathId="2" />
			</Actions>
			<ObjectPaths>
				<Identity Id="2" Name="740c6a0b-85e2-48a0-a494-e0f1759d4aa7:site:` + siteR.Data().ID + `:web:` + webR.Data().ID + `" />
				<Property Id="4" ParentId="2" Name="AllProperties" />
			</ObjectPaths>
		</Request>
	`))
	return sp.ProcessQuery(webProps.client.AuthCnfg.GetSiteURL(), body)
}
