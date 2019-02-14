# gosip - SharePoint HTTP client for Go

<p align="center">
  <img src="./assets/gosip.png" />
</p>

> This is early draft version. Lot's of improvements and breaking changes are expected in the near future until GA version is proclaimed.

Don't throw bananas at me, it's my first steps and experiments in Golang field. =)

## Main features

`gosip` allows you to perform SharePoint unattended (without user interaction) http authentication with Go (Golang) using different authentication strategies.

Supported SharePoint versions:

- SharePoint Online (SPO)
- On-Prem: 2019, 2016, and 2013

Authentication strategies:

- SharePoint 2013, 2016, 2019:
  - ADFS user credentials (ADFS, WAP -> Basic, WAP -> ADFS)
  - Behind reverse proxy (TMG, WAP -> Basic, WAP -> ADFS)
  - Form-based authentication (FBA)
- SharePoint Online:
  - Addin only permissions
  - SAML based with user credentials
  - ADFS user credentials

## Installation

```bash
go get github.com/koltyakov/gosip
```

## Usage insights

1\. Understand SharePoint environment type and authentication strategy.

Let's assume it's, SharePoint Online and Addin Only permissions. Then `github.com/koltyakov/gosip/auth/addin` subpackage should be used.

2\. Initiate authentication object.

```golang
auth := &addin.AuthCnfg{
	SiteURL:      os.Getenv("SPAUTH_SITEURL"),
	ClientID:     os.Getenv("SPAUTH_CLIENTID"),
	ClientSecret: os.Getenv("SPAUTH_CLIENTSECRET"),
}
```

AuthCnfg's from different strategies contains different options.

The authentication options can be provided explicitly or can be read from a configuration file.

```golang
configPath := "./config/private.adfs.json"
auth := &adfs.AuthCnfg{}

err := auth.ReadConfig(configPath)
if err != nil {
	fmt.Printf("Unable to get config: %v\n", err)
	return
}
```

3\. Use SP awared HTTP client from `github.com/koltyakov/gosip` package.

```golang
client := &gosip.SPClient{
	AuthCnfg: auth,
}

var req *http.Request
// Initiate API request
// ...

resp, err := client.Execute(req)
if err != nil {
	fmt.Printf("Unable to request api: %v", err)
	return
}
```

SPClient has `Execute` method which is a wrapper function injecting SharePoint authentication and ending up calling http.Client's `Do` method.

## Usage samples

### Addin Only Permissions

```golang
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/auth/addin"
)

func main() {
	auth := &addin.AuthCnfg{
		SiteURL:      os.Getenv("SPAUTH_SITEURL"),
		ClientID:     os.Getenv("SPAUTH_CLIENTID"),
		ClientSecret: os.Getenv("SPAUTH_CLIENTSECRET"),
	}

	client := &gosip.SPClient{
		AuthCnfg: auth,
	}

	apiEndpoint := auth.GetSiteURL() + "/_api/web?$select=Title"
	req, err := http.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		log.Fatalf("unable to create a request: %v\n", err)
	}

	req.Header.Set("Accept", "application/json;odata=minimalmetadata")

	resp, err := client.Execute(req)
	if err != nil {
		log.Fatalf("unable to request api: %v\n", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("unable to read a response: %v\n", err)
	}

	// No JSON unmarshalling for simplicity
	fmt.Printf("response: %s\n", data)
}
```

### ADFS with config reader

```golang
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/auth/adfs"
)

func main() {
	configPath := "./config/private.adfs.json"
	auth := &adfs.AuthCnfg{}

	err := auth.ReadConfig(configPath)
	if err != nil {
		log.Fatalf("unable to get config: %v\n", err)
	}

	client := &gosip.SPClient{
		AuthCnfg: auth,
	}

	apiEndpoint := auth.GetSiteURL() + "/_api/web?$select=Title"
	req, err := http.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		log.Fatalf("unable to create a request: %v\n", err)
	}

	req.Header.Set("Accept", "application/json;odata=verbose")

	resp, err := client.Execute(req)
	if err != nil {
		log.Fatalf("unable to request api: %v\n", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("unable to read a response: %v\n", err)
	}

	fmt.Printf("response: %s\n", data)
}
```

### Basic auth (NTML)

```golang
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/auth/basic"
)

func main() {
	configPath := "./config/private.basic.json"
	auth := &basic.AuthCnfg{}

	err := auth.ReadConfig(configPath)
	if err != nil {
		log.Fatalf("unable to get config: %v\n", err)
	}

	client := &gosip.SPClient{
		AuthCnfg: auth,
	}

	apiEndpoint := auth.GetSiteURL() + "/_api/web?$select=Title"
	req, err := http.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		log.Fatalf("unable to create a request: %v", err)
	}

	req.Header.Set("Accept", "application/json;odata=verbose")

	resp, err := client.Execute(req)
	if err != nil {
		log.Fatalf("unable to request api: %v\n", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("unable to read a response: %v\n", err)
	}

	fmt.Printf("response: %s\n", data)
}
```

## Tests

### Run automated tests

Create auth credentials store files in `./config` folder for corresponding strategies:

- private.addin.json
- private.adfs.json
- private.basic.json
- private.fba.json
- private.saml.json
- private.tmg.json

Auth configs should have the same structure as [node-sp-auth's](https://github.com/s-kainet/node-sp-auth) configs.

```bash
go test ./...
```

### Run manual test

Modify `cmd/gosip/main.go` to include required scenarios and run:

```bash
go run cmd/gosip/main.go
```

## Reference

A lot of stuff for auth flows have been "copied" from [node-sp-auth](https://github.com/s-kainet/node-sp-auth) library (used as blueprint), which we intensively use in Node.js ecosystem for years.
