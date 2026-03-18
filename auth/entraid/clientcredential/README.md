# Entra ID Client Credential (Certificate) Auth

App-only authentication using an Entra ID app registration with a client certificate. This strategy does not require user interaction and is suitable for background services and daemons.

## Azure App Registration

1. Create or use an existing app registration in [Azure Portal](https://portal.azure.com/#blade/Microsoft_AAD_RegisteredApps/ApplicationsListBlade).
2. Under **Certificates & secrets**, upload your certificate (.cer/.pem public key).
3. Under **API permissions**, add the SharePoint permissions your application requires (e.g. `Sites.Read.All`). Grant admin consent.

## Configuration

Create a JSON config file (e.g. `private.json`):

```json
{
  "siteUrl": "https://contoso.sharepoint.com/sites/yoursite",
  "tenantId": "00000000-0000-0000-0000-000000000000",
  "clientId": "00000000-0000-0000-0000-000000000000",
  "certPath": "./certs/cert.pfx",
  "certPass": "certificate-password"
}
```

| Field | Description |
|-------|-------------|
| `siteUrl` | SharePoint site URL |
| `tenantId` | Entra ID tenant ID |
| `clientId` | Entra ID app (client) ID |
| `certPath` | Path to `.pfx` certificate file (absolute or relative to config file) |
| `certPass` | Certificate password (supports encrypted values via `cpass`) |

### Encrypting secrets with cpass

The `certPass` field supports encrypted values via [cpass](../../../cmd/cpass/README.md). Encrypted secrets can only be decrypted on the machine where they were generated, reducing accidental leaks.

```bash
go run ./cmd/cpass/main.go -secret "certificate-password"
```

Use the encrypted output in your `private.json`:

```json
{
  "certPass": "encrypted-value-here"
}
```

To use a custom master key instead of the machine-bound default:

```bash
go run ./cmd/cpass/main.go -secret "certificate-password" -master "my-master-key"
```

Then set the master key on the auth config before reading:

```go
auth := &clientcredential.AuthCnfg{}
auth.SetMasterkey("my-master-key")
if err := auth.ReadConfig("private.json"); err != nil {
	log.Fatal(err)
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
 "github.com/koltyakov/gosip/auth/entraid/clientcredential"
)

func main() {
 auth := &clientcredential.AuthCnfg{
  Base: entraid.Base{
   SiteURL:  "https://contoso.sharepoint.com/sites/yoursite",
   TenantID: "00000000-0000-0000-0000-000000000000",
   ClientID: "00000000-0000-0000-0000-000000000000",
  },
  CertPath: "/absolute/path/to/cert.pfx",
  CertPass: "certificate-password",
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

Save an inline configuration to a file for reuse. Secrets are automatically encrypted via `cpass`:

```go
auth := &clientcredential.AuthCnfg{
	Base: entraid.Base{
		SiteURL:  "https://contoso.sharepoint.com/sites/yoursite",
		TenantID: "00000000-0000-0000-0000-000000000000",
		ClientID: "00000000-0000-0000-0000-000000000000",
	},
	CertPath: "/path/to/cert.pfx",
	CertPass: "certificate-password",
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
 "github.com/koltyakov/gosip/auth/entraid/clientcredential"
)

func main() {
 auth := &clientcredential.AuthCnfg{}
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
