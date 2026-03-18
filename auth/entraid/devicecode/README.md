# Entra ID Device Code Auth

Interactive authentication using the [device code flow](https://learn.microsoft.com/en-us/entra/identity-platform/v2-oauth2-device-code). The user is prompted to open a URL in a browser and enter a code to sign in. Supports MFA.

## Azure App Registration

1. Create or use an existing app registration in [Azure Portal](https://portal.azure.com/#blade/Microsoft_AAD_RegisteredApps/ApplicationsListBlade).
2. Under **Authentication > Settings**, set **Allow public client flows** to **Enabled**.
3. Under **API permissions**, add the *SharePoint* permissions your application requires (e.g. `Sites.Read.All`).

## Configuration

Create a JSON config file (e.g. `private.json`):

```json
{
  "siteUrl": "https://contoso.sharepoint.com/sites/yoursite",
  "tenantId": "00000000-0000-0000-0000-000000000000",
  "clientId": "00000000-0000-0000-0000-000000000000",
}
```

| Field | Description |
|-------|-------------|
| `siteUrl` | SharePoint site URL |
| `tenantId` | Entra ID tenant ID |
| `clientId` | Entra ID app (client) ID |

### Custom device code message handler

By default the device code message is printed to stdout. Set `UserPrompt` to handle it differently:

```go
auth := &devicecode.AuthCnfg{
 Base: entraid.Base{ /* ... */ },
 UserPrompt: func(_ context.Context, msg azidentity.DeviceCodeMessage) error {
  return os.WriteFile("device_code.txt", []byte(msg.Message), 0600) // write the device_code prompt to a file
 },
}
```

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
 "github.com/koltyakov/gosip/auth/entraid/devicecode"
)

func main() {
 auth := &devicecode.AuthCnfg{
  Base: entraid.Base{
   SiteURL:  "https://contoso.sharepoint.com/sites/yoursite",
   TenantID: "00000000-0000-0000-0000-000000000000",
   ClientID: "00000000-0000-0000-0000-000000000000",
  }
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
auth := &devicecode.AuthCnfg{
	Base: entraid.Base{
		SiteURL:  "https://contoso.sharepoint.com/sites/yoursite",
		TenantID: "00000000-0000-0000-0000-000000000000",
		ClientID: "00000000-0000-0000-0000-000000000000",
	},
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
 "github.com/koltyakov/gosip/auth/entraid/devicecode"
)

func main() {
 auth := &devicecode.AuthCnfg{}
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
