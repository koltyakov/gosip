# Sovereign Cloud Support

This document describes the sovereign cloud support added to gosip authentication strategies.

## Overview

Support has been added for three new Microsoft sovereign clouds:

- **Azure Bleu** (France) - Azure AD: `login.sovcloud-identity.fr`
- **Azure Delos** (Germany) - Azure AD: `login.sovcloud-identity.de`
- **Azure GovSG** (Singapore) - Azure AD: `login.sovcloud-identity.sg`

## Supported Authentication Strategies

The following authentication strategies now support sovereign clouds:

1. **azurecert** - Certificate-based authentication
2. **azurecreds** - Username/password authentication
3. **device** - Device flow authentication
4. **saml** - SAML authentication
5. **addin** - Add-in authentication

## Configuration

### Explicit Cloud Configuration

For **azurecert**, **azurecreds**, and **device** authentication, you can explicitly specify the cloud environment using the `cloud` field:

```json
{
  "siteUrl": "https://your-sharepoint-url",
  "tenantId": "your-tenant-id",
  "clientId": "your-client-id",
  "cloud": "bleu"
}
```

#### Available Cloud Values

- `public` - Commercial cloud (login.microsoftonline.com) - default
- `usgovernment` - US Government cloud (login.microsoftonline.us)
- `china` - China cloud (login.chinacloudapi.cn)
- `germany` - Germany cloud (login.microsoftonline.de)
- `bleu` - France sovereign cloud (login.sovcloud-identity.fr)
- `delos` - Germany sovereign cloud (login.sovcloud-identity.de)
- `govsg` - Singapore sovereign cloud (login.sovcloud-identity.sg)

### Auto-Detection

For **saml** and **addin** authentication, cloud environments are auto-detected from the SharePoint URL pattern:

- `*.sharepoint.com` → Public cloud
- `*.sharepoint.us` → US Government cloud
- `*.sharepoint.cn` → China cloud
- `*.sharepoint.de` → Germany cloud

**Note**: SharePoint domain patterns for the new sovereign clouds (Bleu, Delos, GovSG) are not yet publicly known. Once Microsoft announces these domains, auto-detection will be enabled.

## Usage Examples

### Azure Certificate Auth (azurecert)

```json
{
  "siteUrl": "https://contoso-france.TBD",
  "tenantId": "e4d43069-8ecb-49c4-8178-5bec83c53e9d",
  "clientId": "628cc712-c9a4-48f0-a059-af64bdbb4be5",
  "certPath": "cert.pfx",
  "certPass": "password",
  "cloud": "bleu"
}
```

### Azure Credentials Auth (azurecreds)

```json
{
  "siteUrl": "https://contoso-germany.TBD",
  "tenantId": "e4d43069-8ecb-49c4-8178-5bec83c53e9d",
  "clientId": "628cc712-c9a4-48f0-a059-af64bdbb4be5",
  "username": "user@contoso.de",
  "password": "password",
  "cloud": "delos"
}
```

### Device Flow Auth (device)

```json
{
  "siteUrl": "https://contoso-singapore.TBD",
  "clientId": "61367a97-562c-4372-a9ee-b35307abdd26",
  "tenantId": "3f83fe32-29b2-488e-8c3f-c8b7a2e19a2f",
  "cloud": "govsg"
}
```

## Implementation Details

### Certificate Auth (azurecert)

- Added `AzureCloud` type and constants
- Added optional `Cloud` field to `AuthCnfg`
- Updated `getAADEndpoint()` to check explicit cloud config first, then auto-detect
- Added `getCloudEndpoint()` helper function for cloud-to-endpoint mapping

### Credentials Auth (azurecreds)

- Added `AzureCloud` type and constants
- Added optional `Cloud` field to `AuthCnfg`
- Set `AADEndpoint` on auth config when cloud is explicitly specified
- Added `getCloudEndpoint()` helper function

### Device Flow Auth (device)

- Added `AzureCloud` type and constants
- Added optional `Cloud` field to `AuthCnfg`
- Set `AADEndpoint` on device flow config when cloud is explicitly specified
- Added `getCloudEndpoint()` helper function

### SAML Auth (saml)

- Added new `spoEnv` constants: `spoBleu`, `spoDelos`, `spoGovSG`
- Added endpoint mappings in `loginEndpoints` map
- Commented auto-detection logic (awaiting SharePoint domain announcements)

### Add-in Auth (addin)

- Added new `spoEnv` constants: `spoBleu`, `spoDelos`, `spoGovSG`
- Added endpoint mappings in `accEndpoints` map
- Commented auto-detection logic (awaiting SharePoint domain announcements)

## Migration Guide

Existing configurations continue to work without changes. The `cloud` field is optional and defaults to auto-detection or public cloud.

To use a sovereign cloud:

1. Add the `cloud` field to your configuration
2. Set it to the appropriate cloud value
3. Update your `siteUrl` to match your organization's SharePoint URL in that cloud

## Future Updates

When Microsoft announces the SharePoint domain patterns for the new sovereign clouds, the following updates will be made:

1. Update auto-detection logic in `azurecert/azurecert.go`
2. Enable commented detection code in `saml/env.go`
3. Enable commented detection code in `addin/env.go`
4. Update this documentation with actual domain patterns

## References

- [GitHub Issue #92](https://github.com/koltyakov/gosip/issues/92)
- [MSAL.NET PR #5671](https://github.com/AzureAD/microsoft-authentication-library-for-dotnet/pull/5671)
- [MSAL.NET PR #5709](https://github.com/AzureAD/microsoft-authentication-library-for-dotnet/pull/5709)
