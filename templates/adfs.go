package templates

import (
	"bytes"
	"text/template"
)

// AdfsSamlWsfedTemplate : AdfsSamlWsfedTemplate template
func AdfsSamlWsfedTemplate(to, username, password, relyingParty string) (string, error) {
	type adfsSamlWsfed struct {
		To           string
		Username     string
		Password     string
		RelyingParty string
	}

	t, err := template.New("adfsSamlWsfed").Parse(`
		<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope" xmlns:a="http://www.w3.org/2005/08/addressing" xmlns:u="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd">
			<s:Header>
				<a:Action s:mustUnderstand="1">http://docs.oasis-open.org/ws-sx/ws-trust/200512/RST/Issue</a:Action>
				<a:ReplyTo>
					<a:Address>http://www.w3.org/2005/08/addressing/anonymous</a:Address>
				</a:ReplyTo>
				<a:To s:mustUnderstand="1">{{.To}}</a:To>
				<o:Security s:mustUnderstand="1" xmlns:o="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd">
					<o:UsernameToken u:Id="uuid-7b105801-44ac-4da7-aa69-a87f9db37299-1">
						<o:Username>{{.Username}}</o:Username>
						<o:Password Type="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-username-token-profile-1.0#PasswordText">{{.Password}}</o:Password>
					</o:UsernameToken>
				</o:Security>
			</s:Header>
			<s:Body>
				<trust:RequestSecurityToken xmlns:trust="http://docs.oasis-open.org/ws-sx/ws-trust/200512">
					<wsp:AppliesTo xmlns:wsp="http://schemas.xmlsoap.org/ws/2004/09/policy">
						<wsa:EndpointReference xmlns:wsa="http://www.w3.org/2005/08/addressing">
							<wsa:Address>{{.RelyingParty}}</wsa:Address>
						</wsa:EndpointReference>
					</wsp:AppliesTo>
					<trust:KeyType>http://docs.oasis-open.org/ws-sx/ws-trust/200512/Bearer</trust:KeyType>
					<trust:RequestType>http://docs.oasis-open.org/ws-sx/ws-trust/200512/Issue</trust:RequestType>
				</trust:RequestSecurityToken>
			</s:Body>
		</s:Envelope>
	`)
	if err != nil {
		return "", err
	}

	data := adfsSamlWsfed{
		To:           to,
		Username:     escapeParamString(username),
		Password:     escapeParamString(password), // + "1",
		RelyingParty: relyingParty,
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return "", err
	}

	result := compactTemplate(tpl.String())

	return result, nil
}

// AdfsSamlTokenTemplate : AdfsSamlTokenTemplate template
func AdfsSamlTokenTemplate(token []byte, notBefore, notAfter, relyingParty string) (string, error) {
	type adfsSamlToken struct {
		Token        string
		NotBefore    string
		NotOnOrAfter string
		RelyingParty string
	}

	t, err := template.New("adfsSamlToken").Parse(`
		<t:RequestSecurityTokenResponse xmlns:t="http://schemas.xmlsoap.org/ws/2005/02/trust">
			<t:Lifetime>
				<wsu:Created xmlns:wsu="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd">{{.NotBefore}}</wsu:Created>
				<wsu:Expires xmlns:wsu="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd">{{.NotOnOrAfter}}</wsu:Expires>
			</t:Lifetime>
			<wsp:AppliesTo xmlns:wsp="http://schemas.xmlsoap.org/ws/2004/09/policy">
				<wsa:EndpointReference xmlns:wsa="http://www.w3.org/2005/08/addressing">
					<wsa:Address>{{.RelyingParty}}</wsa:Address>
				</wsa:EndpointReference>
			</wsp:AppliesTo>
			<t:RequestedSecurityToken>{{.Token}}</t:RequestedSecurityToken>
			<t:TokenType>urn:oasis:names:tc:SAML:1.0:assertion</t:TokenType>
			<t:RequestType>http://schemas.xmlsoap.org/ws/2005/02/trust/Issue</t:RequestType>
			<t:KeyType>http://schemas.xmlsoap.org/ws/2005/05/identity/NoProofKey</t:KeyType>
		</t:RequestSecurityTokenResponse>
	`)
	if err != nil {
		return "", err
	}

	data := adfsSamlToken{
		Token:        string(token),
		NotBefore:    notBefore,
		NotOnOrAfter: notAfter,
		RelyingParty: relyingParty,
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return "", err
	}

	result := compactTemplate(tpl.String())

	return result, nil
}
