package addin

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

	return spoProd // ToDo: Research how to identify Office 365 Dedicated
}
