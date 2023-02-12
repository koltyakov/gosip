# Azure AD Device Auth Flow Sample

The sample shows Gosip [custom auth](https://go.spflow.com/auth/custom-auth) with [AAD Device Token Authorization](https://docs.microsoft.com/en-us/azure/go/azure-sdk-go-authorization#use-device-token-authentication).

## Custom auth implementation

Checkout [the code](./device.go).

## Azure App registration

1\. Create or use existing app registration

2\. Make sure that the app is configured to support device flow

- Authentication settings
  - Public client/native (mobile & desktop)
  - Suggested Redirect URIs for public clients (mobile, desktop) - https://login.microsoftonline.com/common/oauth2/nativeclient - checked
  - Default client type - Yes - for Device code flow [Learn more](https://go.microsoft.com/fwlink/?linkid=2094804)
- App permissions
  - Azure Service Management :: user_impersonation
  - SharePoint :: based on your application requirements
- Manifest
  - oauth2AllowIdTokenImplicitFlow - true
  - oauth2AllowImplicitFlow - true
- etc. based on application needs

## Auth configuration and usage

```golang
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/api"
	strategy "github.com/koltyakov/gosip/auth/device"
)

func main() {

	authCnfg := &strategy.AuthCnfg{
		SiteURL:  os.Getenv("SPAUTH_SITEURL"),
		ClientID: os.Getenv("SPAUTH_AAD_CLIENTID"),
		TenantID: os.Getenv("SPAUTH_AAD_TENANTID"),
	}

	client := &gosip.SPClient{AuthCnfg: authCnfg}
	sp := api.NewSP(client)

	res, err := sp.Web().Select("Title").Get()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Site title: %s\n", res.Data().Title)

}
```
