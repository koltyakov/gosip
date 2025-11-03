# Gosip - SharePoint SDK for Go (Golang)

> Authentication, HTTP client & fluent API wrapper

<!-- ![Build Status](https://koltyakov.visualstudio.com/SPNode/_apis/build/status/gosip?branchName=master) -->

[![Go Report Card](https://goreportcard.com/badge/github.com/koltyakov/gosip)](https://goreportcard.com/report/github.com/koltyakov/gosip)
[![GoDoc](https://godoc.org/github.com/koltyakov/gosip?status.svg)](https://godoc.org/github.com/koltyakov/gosip)
[![License](https://img.shields.io/github/license/koltyakov/gosip.svg)](https://github.com/koltyakov/gosip/blob/master/LICENSE)
[![codecov](https://codecov.io/gh/koltyakov/gosip/branch/master/graph/badge.svg)](https://codecov.io/gh/koltyakov/gosip)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fkoltyakov%2Fgosip.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fkoltyakov%2Fgosip?ref=badge_shield)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)

<!--suppress HtmlDeprecatedAttribute -->
<div align="center">
  <img alt="Gosip" src="https://raw.githubusercontent.com/koltyakov/gosip-docs/master/.gitbook/assets/gosip.png" />
</div>

## Main features

- Unattended authentication using different strategies.
- Fluent API syntax for SharePoint object model.
- Simplified API consumption (REST, CSOM, SOAP).
- SharePoint-aware embedded features (retries, header presets, error handling).

### Supported SharePoint versions

- SharePoint Online (SPO)
- On-Premises (2019/2016/2013)

### Supported auth strategies

- SharePoint Online:

  - Azure Certificate (App Only) [🔗](https://go.spflow.com/auth/strategies/azure-certificate-auth)
  - Azure Username/Password [🔗](https://go.spflow.com/auth/strategies/azure-creds-auth)
  - Azure Device Flow [🔗](https://go.spflow.com/auth/strategies/azure-device-flow)
  - ~~SAML based with user credentials~~ (deprecated by platform)
  - Add-In only permissions
  - ADFS user credentials
  - On-Demand auth

- SharePoint On-Premises 2019/2016/2013:
  - User credentials (NTLM)
  - ADFS user credentials (ADFS, WAP -> Basic/NTLM, WAP -> ADFS)
  - Behind a reverse proxy (Forefront TMG, WAP -> Basic/NTLM, WAP -> ADFS)
  - Form-based authentication (FBA)
  - On-Demand auth

## Installation

```bash
go get github.com/koltyakov/gosip
```

## Usage insights

1\. Understand SharePoint environment type and authentication strategy.

Let's assume it's SharePoint Online and Add-In Only permissions. Then `strategy "github.com/koltyakov/gosip/auth/addin"` subpackage should be used.

```golang
package main

import (
	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/api"
	strategy "github.com/koltyakov/gosip/auth/addin"
)
```

2\. Initiate an authentication object.

```golang
auth := &strategy.AuthCnfg{
	SiteURL:      os.Getenv("SPAUTH_SITEURL"),
	ClientID:     os.Getenv("SPAUTH_CLIENTID"),
	ClientSecret: os.Getenv("SPAUTH_CLIENTSECRET"),
}
```

AuthCnfg from different strategies contains different options relevant for a specified auth type.

The authentication options can be provided explicitly or can be read from a configuration file.

```golang
configPath := "./config/private.json"
auth := &strategy.AuthCnfg{}

err := auth.ReadConfig(configPath)
if err != nil {
	fmt.Printf("Unable to get config: %v\n", err)
	return
}
```

3\. Bind auth client with Fluent API.

```golang
client := &gosip.SPClient{AuthCnfg: auth}

sp := api.NewSP(client)

res, err := sp.Web().Select("Title").Get()
if err != nil {
	fmt.Println(err)
}

fmt.Printf("%s\n", res.Data().Title)
```

## Usage samples

### Fluent API client

Fluent API gives a simple way of constructing API endpoint calls with IntelliSense and chainable syntax.

![Fluent Sample](https://raw.githubusercontent.com/koltyakov/gosip-docs/master/.gitbook/assets/fluent.gif)

```golang
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/api"
	strategy "github.com/koltyakov/gosip/auth/addin"
)

func main() {
	// Getting auth params and client
	client, err := getAuthClient()
	if err != nil {
		log.Fatalln(err)
	}

	// Binding SharePoint API
	sp := api.NewSP(client)

	// Custom headers
	headers := map[string]string{
		"Accept": "application/json;odata=minimalmetadata",
		"Accept-Language": "de-DE,de;q=0.9",
	}
	config := &api.RequestConfig{Headers: headers}

	// Chainable request sample
	data, err := sp.Conf(config).Web().Lists().Select("Id,Title").Get()
	if err != nil {
		log.Fatalln(err)
	}

	// Response object unmarshalling (struct depends on OData mode and API method)
	res := &struct {
		Value []struct {
			ID    string `json:"Id"`
			Title string `json:"Title"`
		} `json:"value"`
	}{}

	if err := json.Unmarshal(data, &res); err != nil {
		log.Fatalf("unable to parse the response: %v", err)
	}

	for _, list := range res.Value {
		fmt.Printf("%+v\n", list)
	}

}

func getAuthClient() (*gosip.SPClient, error) {
	configPath := "./config/private.spo-addin.json"
	auth := &strategy.AuthCnfg{}
	if err := auth.ReadConfig(configPath); err != nil {
		return nil, fmt.Errorf("unable to get config: %v", err)
	}
	return &gosip.SPClient{AuthCnfg: auth}, nil
}
```

### Generic HTTP client helper

Provides generic GET/POST helpers for REST operations, reducing the amount of `http.NewRequest` scaffolded code, can be used for custom or not covered with Fluent API endpoints.

```golang
package main

import (
	"fmt"
	"log"

	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/api"
	strategy "github.com/koltyakov/gosip/auth/ntlm"
)

func main() {
	configPath := "./config/private.ntlm.json"
	auth := &strategy.AuthCnfg{}

	if err := auth.ReadConfig(configPath); err != nil {
		log.Fatalf("unable to get config: %v\n", err)
	}

	sp := api.NewHTTPClient(&gosip.SPClient{AuthCnfg: auth})

	endpoint := auth.GetSiteURL() + "/_api/web?$select=Title"

	data, err := sp.Get(endpoint, nil)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	// sp.Post(endpoint, body, nil) // generic POST
	// sp.Delete(endpoint, nil) // generic DELETE helper crafts "X-HTTP-Method"="DELETE" header
	// sp.Update(endpoint, nil) // generic UPDATE helper crafts "X-HTTP-Method"="MERGE" header
	// sp.ProcessQuery(endpoint, body) // CSOM helper (client.svc/ProcessQuery)

	fmt.Printf("response: %s\n", data)
}
```

### Low-level HTTP client usage

Low-lever SharePoint-aware HTTP client from `github.com/koltyakov/gosip` package for custom or not covered with a Fluent API client endpoints with granular control for an HTTP request, response, and http.Client parameters. The client is used internally but rarely required in consumer code.

```golang
client := &gosip.SPClient{AuthCnfg: auth}

var req *http.Request
// Initiate API request
// ...

resp, err := client.Execute(req)
if err != nil {
	fmt.Printf("Unable to request api: %v", err)
	return
}
```

SPClient has `Execute` method which is a wrapper function injecting SharePoint authentication and ending up calling `http.Client`'s `Do` method.

## Authentication strategies

Auth strategy should be selected corresponding to your SharePoint environment and its configuration.

Import path `strategy "github.com/koltyakov/gosip/auth/{strategy}"`. Where `/{strategy}` stands for a strategy auth package.

Azure AD based strategies (recommended production use with SharePoint Online):

| `/{strategy}` | Description                                       | Credentials sample(s)                                                   |
| ------------- | ------------------------------------------------- | ----------------------------------------------------------------------- |
| `/azurecert`  | Azure AD Certificate authentication               | [details](https://go.spflow.com/auth/strategies/azure-certificate-auth) |
| `/azurecreds` | Azure AD authorization with username and password | [details](https://go.spflow.com/auth/strategies/azure-creds-auth)       |
| `/azureenv`   | Azure AD environment-based authentication         | [details](https://go.spflow.com/auth/strategies/azure-environment-auth) |
| `/device`     | Azure AD Device Token authentication              | [details](https://go.spflow.com/auth/strategies/azure-device-flow)      |

Other strategies:

| `/{strategy}` | SPO                   | On-Prem | Credentials sample(s)                                       |
| ------------- | --------------------- | ------- | ----------------------------------------------------------- |
| `/adfs`       | ✅                    | ✅      | [details](https://go.spflow.com/auth/strategies/adfs)       |
| `/ondemand`   | ✅                    | ✅      | [details](https://go.spflow.com/auth/custom-auth/on-demand) |
| `/addin`      | ✅                    | ❌      | [details](https://go.spflow.com/auth/strategies/addin)      |
| `/ntlm`       | ❌                    | ✅      | [details](https://go.spflow.com/auth/strategies/ntlm)       |
| `/fba`        | ❌                    | ✅      | [details](https://go.spflow.com/auth/strategies/fba)        |
| `/tmg`        | ❌                    | ✅      | [details](https://go.spflow.com/auth/strategies/tmg)        |
| `/saml`       | ✅ -> ❌ (deprecated) | ❌      | [details](https://go.spflow.com/auth/strategies/saml)       |

Environment should configured for a specific auth strategy. E.g. you won't succeed with `adfs` in SPO if it has not setup properly.

Below are the most commonly authentication methods in more details:

### Azure AD application authentication

Azure AD application authentication is a recommended way for production use with SharePoint Online. It's based on OAuth 2.0 protocol and uses AAD application credentials for authentication.

Depending on an application type, there are different authentication strategies:

- Azure AD Certificate authentication: for headless applications, which are not able to provide user interaction, like a background service or a daemon. It uses a certificate to authenticate an application.
- Azure AD authorization with username and password: for applications which are able to provide user interaction, like a desktop application or CLI with credentials prompt. It uses a username and password to authenticate a user.
- Azure AD device token authentication: for applications which are able to provide user interaction, like a desktop application or CLI. It uses a device code to authenticate a user. It also supports multi-factor authentication.

### AddIn Only Auth

This type of authentication uses AddIn Only policy and OAuth bearer tokens for authenticating HTTP requests.

```golang
// AuthCnfg - AddIn Only auth config structure
type AuthCnfg struct {
	// SPSite or SPWeb URL, which is the context target for the API calls
	SiteURL string `json:"siteUrl"`
	// Client ID obtained when registering the AddIn
	ClientID string `json:"clientId"`
	// Client Secret obtained when registering the AddIn
	ClientSecret string `json:"clientSecret"`
	// Your SharePoint Online tenant ID (optional)
	Realm string `json:"realm"`
}
```

Realm can be left empty or filled in, which will add small performance improvement. The easiest way to find the tenant is to open SharePoint Online site collection, click `Site Settings` -> `Site App Permissions`. Taking any random app, the tenant ID (realm) is the GUID part after the `@`.

See more details of [AddIn Configuration and Permissions](https://github.com/s-kainet/node-sp-auth/wiki/SharePoint-Online-addin-only-authentication).

### NTLM Auth (NTLM handshake)

This type of authentication uses an HTTP NTLM handshake to obtain an authentication header.

```golang
// AuthCnfg - NTLM auth config structure
type AuthCnfg struct {
	// SPSite or SPWeb URL, which is the context target for the API calls
	SiteURL  string `json:"siteUrl"`
	Domain   string `json:"domain"`   // AD domain name
	Username string `json:"username"` // AD user name
	Password string `json:"password"` // AD user password
}
```

Gosip uses `github.com/Azure/go-ntlmssp` NTLM negotiator, however, a custom one also can be [provided](https://github.com/koltyakov/gosip/issues/14) in case of demand.

## Secrets encoding

When storing credential in local `private.json` files, which can be handy in local development scenarios, we strongly recommend to encode secrets such as `password` or `clientSecret` using [cpass](./cmd/cpass/README.md). Class converts a secret to an encrypted representation, which can only be decrypted on the same machine where it was generated. That reduces accidental leaks, e.g. together with git commits.

## Reference

Many auth flows have been "copied" from [node-sp-auth](https://github.com/s-kainet/node-sp-auth) library (used as a blueprint), which we intensively use in Node.js ecosystem for years.

Fluent API and wrapper syntax are inspired by [PnPjs](https://github.com/pnp/pnpjs), which is also the first-class citizen on almost all our Node.js and front-end projects with SharePoint involved.

## 📚 [Documentation](https://go.spflow.com)

## 📦 [Samples](https://github.com/koltyakov/gosip-sandbox/tree/master/samples)

## License

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fkoltyakov%2Fgosip.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fkoltyakov%2Fgosip?ref=badge_large)
