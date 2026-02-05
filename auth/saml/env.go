package saml

import (
	"net/url"
	"strings"
)

type spoEnv int32

const (
	spoProd spoEnv = iota
	spoGerman
	spoChina
	spoUSGov
	spoUSDef
	spoBleu  // France sovereign cloud
	spoDelos // Germany sovereign cloud
	spoGovSG // Singapore sovereign cloud
)

// resolveSPOEnv resolves SPO hosting environment type
func resolveSPOEnv(siteURL string) spoEnv {
	parsedURL, err := url.Parse(siteURL)
	if err != nil {
		return spoProd
	}

	if strings.Contains(parsedURL.Host, ".sharepoint.com") {
		return spoProd
	}
	if strings.Contains(parsedURL.Host, ".sharepoint.de") {
		return spoGerman
	}
	if strings.Contains(parsedURL.Host, ".sharepoint.cn") {
		return spoChina
	}
	if strings.Contains(parsedURL.Host, ".sharepoint-mil.us") {
		return spoUSGov
	}
	if strings.Contains(parsedURL.Host, ".sharepoint.us") {
		return spoUSDef
	}

	// TODO: Add new sovereign cloud domain patterns when SharePoint domains are announced
	// if strings.Contains(parsedURL.Host, ".TBD") {  // Bleu (France)
	//     return spoBleu
	// }
	// if strings.Contains(parsedURL.Host, ".TBD") {  // Delos (Germany)
	//     return spoDelos
	// }
	// if strings.Contains(parsedURL.Host, ".TBD") {  // GovSG (Singapore)
	//     return spoGovSG
	// }

	return spoProd // ToDo: Research how to identify Office 365 Dedicated
}
