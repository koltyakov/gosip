package adfs

import (
	"net/url"
	"testing"
	"time"
)

func TestHelpersEdgeCases(t *testing.T) {

	t.Run("GetAuth/EmptySiteURL", func(t *testing.T) {
		cnfg := &AuthCnfg{SiteURL: ""}
		if _, _, err := GetAuth(cnfg); err == nil {
			t.Error("empty SiteURL should not go")
		}
	})

	t.Run("adfsAuthFlow/EmptyAdfsURL", func(t *testing.T) {
		cnfg := &AuthCnfg{AdfsURL: ""}
		if _, _, err := GetAuth(cnfg); err == nil {
			t.Error("empty AdfsURL should not go")
		}
	})

	t.Run("CleanAuthCache", func(t *testing.T) {
		cnfg := &AuthCnfg{
			AdfsURL:  "https://contoso.sharepoint.com",
			Username: "username",
			Password: "password",
		}
		parsedURL, _ := url.Parse(cnfg.SiteURL)
		cacheKey := parsedURL.Host + "@adfs@" + cnfg.Username + "@" + cnfg.Password
		storage.Set(cacheKey, "token", 1*time.Minute)

		if err := cnfg.CleanAuthCache(); err != nil {
			t.Errorf("can't clean auth cache: %s", err)
		}

		if _, found := storage.Get(cacheKey); found {
			t.Error("auth cache was not cleaned")
		}
	})

}
