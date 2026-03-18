package entraid

import (
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
)

// CloudFromSiteURL returns the Azure cloud configuration based on the SharePoint domain.
func CloudFromSiteURL(siteURL string) cloud.Configuration {
	u, err := url.Parse(siteURL)
	if err != nil {
		return cloud.AzurePublic
	}
	switch {
	case strings.HasSuffix(u.Host, ".sharepoint.us"):
		return cloud.AzureGovernment
	case strings.HasSuffix(u.Host, ".sharepoint.cn"):
		return cloud.AzureChina
	default:
		return cloud.AzurePublic
	}
}
