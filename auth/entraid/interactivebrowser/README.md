# Entra ID Interactive Browser Auth

Interactive authentication that opens a browser window for the user to sign in. Supports MFA.

## Azure App Registration

1. Create or use an existing app registration in [Azure Portal](https://portal.azure.com/#blade/Microsoft_AAD_RegisteredApps/ApplicationsListBlade).
2. Under **Authentication**, add a **Mobile and desktop applications** platform with a redirect URI (e.g. `http://localhost:8400`).
3. Under **Authentication > Settings**, set **Allow public client flows** to **Yes**.
4. Under **API permissions**, add the SharePoint permissions your application requires (e.g. `Sites.Read.All`).

## Configuration

Create a JSON config file (e.g. `private.json`):

```json
{
  "siteUrl": "https://contoso.sharepoint.com/sites/yoursite",
  "tenantId": "00000000-0000-0000-0000-000000000000",
  "clientId": "00000000-0000-0000-0000-000000000000",
  "redirectUrl": "http://localhost:8400"
}
```

| Field | Description |
|-------|-------------|
| `siteUrl` | SharePoint site URL |
| `tenantId` | Entra ID tenant ID |
| `clientId` | Entra ID app (client) ID |
| `redirectUrl` | OAuth redirect URI registered in your Entra ID app |

## Usage

### Inline

```go
package main

import (
 "fmt"
 "log"

 "github.com/koltyakov/gosip"
 "github.com/koltyakov/gosip/api"
 "github.com/koltyakov/gosip/auth/entraid"
 "github.com/koltyakov/gosip/auth/entraid/interactivebrowser"
)

func main() {
 auth := &interactivebrowser.AuthCnfg{
  Base: entraid.Base{
   SiteURL:  "https://contoso.sharepoint.com/sites/yoursite",
   TenantID: "00000000-0000-0000-0000-000000000000",
   ClientID: "00000000-0000-0000-0000-000000000000",
  },
  RedirectURL: "http://localhost:8400",
 }

 client := &gosip.SPClient{AuthCnfg: auth}
 sp := api.NewSP(client)

 res, err := sp.Web().Select("Title").Get()
 if err != nil {
  log.Fatal(err)
 }

 fmt.Printf("Site title: %s\n", res.Data().Title)
}
```

### Writing config from inline

Save an inline configuration to a file for reuse:

```go
auth := &interactivebrowser.AuthCnfg{
	Base: entraid.Base{
		SiteURL:  "https://contoso.sharepoint.com/sites/yoursite",
		TenantID: "00000000-0000-0000-0000-000000000000",
		ClientID: "00000000-0000-0000-0000-000000000000",
	},
	RedirectURL: "http://localhost:8400",
}
if err := auth.WriteConfig("private.json"); err != nil {
	log.Fatal(err)
}
```

### Config file

```go
package main

import (
 "fmt"
 "log"

 "github.com/koltyakov/gosip"
 "github.com/koltyakov/gosip/api"
 "github.com/koltyakov/gosip/auth/entraid/interactivebrowser"
)

func main() {
 auth := &interactivebrowser.AuthCnfg{}
 if err := auth.ReadConfig("private.json"); err != nil {
  log.Fatalf("failed to read config: %v", err)
 }

 client := &gosip.SPClient{AuthCnfg: auth}
 sp := api.NewSP(client)

 res, err := sp.Web().Select("Title").Get()
 if err != nil {
  log.Fatal(err)
 }

 fmt.Printf("Site title: %s\n", res.Data().Title)
}
```

## Troubleshooting

**AADSTS7000218 error**: Enable **Allow public client flows** under Authentication > Advanced settings in your Entra ID app registration, and ensure the platform is set to **Mobile and desktop applications** (not **Web**).
