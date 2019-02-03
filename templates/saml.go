package templates

import (
	"bytes"
	"text/template"
)

// OnlineSamlWsfedTemplate : OnlineSamlWsfedTemplate template
func OnlineSamlWsfedTemplate(endpoint, username, password string) (string, error) {
	type onlineSamlWsfed struct {
		Endpoint string
		Username string
		Password string
	}

	template, err := template.New("onlineSamlWsfed").Parse(`
		<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope" xmlns:a="http://www.w3.org/2005/08/addressing" xmlns:u="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd">
			<s:Header>
				<a:Action s:mustUnderstand="1">http://schemas.xmlsoap.org/ws/2005/02/trust/RST/Issue</a:Action>
				<a:ReplyTo>
					<a:Address>http://www.w3.org/2005/08/addressing/anonymous</a:Address>
				</a:ReplyTo>
				<a:To s:mustUnderstand="1">https://login.microsoftonline.com/extSTS.srf</a:To>
				<o:Security s:mustUnderstand="1" xmlns:o="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd">
					<o:UsernameToken>
						<o:Username>{{.Username}}</o:Username>
						<o:Password>{{.Password}}</o:Password>
					</o:UsernameToken>
				</o:Security>
			</s:Header>
			<s:Body>
				<t:RequestSecurityToken xmlns:t="http://schemas.xmlsoap.org/ws/2005/02/trust">
					<wsp:AppliesTo xmlns:wsp="http://schemas.xmlsoap.org/ws/2004/09/policy">
						<a:EndpointReference>
							<a:Address>{{.Endpoint}}</a:Address>
						</a:EndpointReference>
					</wsp:AppliesTo>
					<t:KeyType>http://schemas.xmlsoap.org/ws/2005/05/identity/NoProofKey</t:KeyType>
					<t:RequestType>http://schemas.xmlsoap.org/ws/2005/02/trust/Issue</t:RequestType>
					<t:TokenType>urn:oasis:names:tc:SAML:1.0:assertion</t:TokenType>
				</t:RequestSecurityToken>
			</s:Body>
		</s:Envelope>
	`)
	if err != nil {
		return "", err
	}

	data := onlineSamlWsfed{
		Endpoint: endpoint,
		Username: escapeParamString(username),
		Password: escapeParamString(password),
	}

	var tpl bytes.Buffer
	if err := template.Execute(&tpl, data); err != nil {
		return "", err
	}

	result := compactTemplate(tpl.String())

	return result, nil
}

// OnlineSamlWsfedAdfsTemplate : OnlineSamlWsfedAdfsTemplate template
func OnlineSamlWsfedAdfsTemplate(endpoint, token string) (string, error) {
	type onlineSamlWsfedAdfs struct {
		Endpoint string
		Token    string
	}

	template, err := template.New("onlineSamlWsfedAdfs").Parse(`
		<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope" xmlns:a="http://www.w3.org/2005/08/addressing" xmlns:u="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd">
			<s:Header>
				<a:Action s:mustUnderstand="1">http://schemas.xmlsoap.org/ws/2005/02/trust/RST/Issue</a:Action>
				<a:ReplyTo>
					<a:Address>http://www.w3.org/2005/08/addressing/anonymous</a:Address>
				</a:ReplyTo>
				<a:To s:mustUnderstand="1">https://login.microsoftonline.com/extSTS.srf</a:To>
				<o:Security s:mustUnderstand="1" xmlns:o="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd">{{.Token}}</o:Security>
			</s:Header>
			<s:Body>
				<t:RequestSecurityToken xmlns:t="http://schemas.xmlsoap.org/ws/2005/02/trust">
					<wsp:AppliesTo xmlns:wsp="http://schemas.xmlsoap.org/ws/2004/09/policy">
						<a:EndpointReference>
							<a:Address>{{.Endpoint}}</a:Address>
						</a:EndpointReference>
					</wsp:AppliesTo>
					<t:KeyType>http://schemas.xmlsoap.org/ws/2005/05/identity/NoProofKey</t:KeyType>
					<t:RequestType>http://schemas.xmlsoap.org/ws/2005/02/trust/Issue</t:RequestType>
					<t:TokenType>urn:oasis:names:tc:SAML:1.0:assertion</t:TokenType>
				</t:RequestSecurityToken>
			</s:Body>
		</s:Envelope>
	`)
	if err != nil {
		return "", err
	}

	data := onlineSamlWsfedAdfs{
		Endpoint: endpoint,
		Token:    token,
	}

	var tpl bytes.Buffer
	if err := template.Execute(&tpl, data); err != nil {
		return "", err
	}

	result := compactTemplate(tpl.String())

	return result, nil
}
