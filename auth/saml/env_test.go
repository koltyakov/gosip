package saml

import (
	"testing"
)

func TestSPOHostingEnvCases(t *testing.T) {

	t.Run("ResolveSPOHostingEnv", func(t *testing.T) {
		if resolveSPOEnv("https://contoso.sharepoint.com") != spoProd {
			t.Error("should be PROD")
		}
		if resolveSPOEnv("https://contoso.com") != spoProd {
			t.Error("should be PROD")
		}
		if resolveSPOEnv("//contoso.com") != spoProd {
			t.Error("should be PROD")
		}
		if resolveSPOEnv("https://contoso.sharepoint.de") != spoGerman {
			t.Error("should be GERMAN")
		}
		if resolveSPOEnv("https://contoso.sharepoint.cn") != spoChina {
			t.Error("should be CHINA")
		}
		if resolveSPOEnv("https://contoso.sharepoint-mil.us") != spoUSGov {
			t.Error("should be USGOV")
		}
		if resolveSPOEnv("https://contoso.sharepoint.us") != spoUSDef {
			t.Error("should be USDEF")
		}
	})

}
