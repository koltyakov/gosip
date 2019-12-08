package api

import (
	"net/url"

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
	sp := NewHTTPClient(webProps.client)
	body := []byte(trimMultiline(`
		<Request xmlns="http://schemas.microsoft.com/sharepoint/clientquery/2009" SchemaVersion="15.0.0.0" LibraryVersion="16.0.0.0" ApplicationName="Gosip">
			<Actions>
				<Method Name="SetFieldValue" Id="9" ObjectPathId="4">
					<Parameters>
						<Parameter Type="String">` + prop + `</Parameter>
						<Parameter Type="String">` + value + `</Parameter>
					</Parameters>
				</Method>
				<Method Name="Update" Id="10" ObjectPathId="2" />
			</Actions>
			<ObjectPaths>
				<StaticProperty Id="0" TypeId="{3747adcd-a3c3-41b9-bfab-4a64dd2f1e0a}" Name="Current" />
				<Property Id="2" ParentId="0" Name="Web" />
				<Property Id="4" ParentId="2" Name="AllProperties" />
			</ObjectPaths>
		</Request>
	`))
	return sp.ProcessQuery(webProps.endpoint, body)
}

// ToDo:
// Write Props with CSOM
