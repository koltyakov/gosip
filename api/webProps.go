package api

import (
	"strconv"

	"github.com/koltyakov/gosip"
)

// WebProps represent SharePoint Web Properties API queryable collection struct
// Always use NewWebProps constructor instead of &WebProps{}
type WebProps struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers *ODataMods
}

// NewWebProps - WebProps struct constructor function
func NewWebProps(client *gosip.SPClient, endpoint string, config *RequestConfig) *WebProps {
	return &WebProps{
		client:    client,
		endpoint:  endpoint,
		config:    config,
		modifiers: NewODataMods(),
	}
}

// ToURL gets endpoint with modificators raw URL
func (webProps *WebProps) ToURL() string {
	return toURL(webProps.endpoint, webProps.modifiers)
}

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (webProps *WebProps) Conf(config *RequestConfig) *WebProps {
	webProps.config = config
	return webProps
}

// Select ...
func (webProps *WebProps) Select(oDataSelect string) *WebProps {
	webProps.modifiers.AddSelect(oDataSelect)
	return webProps
}

// Expand ...
func (webProps *WebProps) Expand(oDataExpand string) *WebProps {
	webProps.modifiers.AddExpand(oDataExpand)
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
	return webProps.SetProps(map[string]string{prop: value})
}

// SetProps ...
func (webProps *WebProps) SetProps(props map[string]string) ([]byte, error) {
	site := NewSP(webProps.client).Site()
	web := NewWeb(webProps.client, getPriorEndpoint(webProps.endpoint, "/AllProperties"), webProps.config)
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
