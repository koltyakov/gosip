package api

import (
	"testing"
)

func TestUser(t *testing.T) {
	checkClient(t)

	sp := NewSP(spClient)
	user := sp.Web().CurrentUser()
	endpoint := spClient.AuthCnfg.GetSiteURL() + "/_api/Web/CurrentUser"

	t.Run("Constructor", func(t *testing.T) {
		users := NewUser(spClient, endpoint, nil)
		if _, err := users.Select("Id").Get(); err != nil {
			t.Error(err)
		}
	})

	t.Run("ToURL", func(t *testing.T) {
		if user.ToURL() != endpoint {
			t.Errorf(
				"incorrect endpoint URL, expected \"%s\", received \"%s\"",
				endpoint,
				user.ToURL(),
			)
		}
	})

	t.Run("Conf", func(t *testing.T) {
		user.config = nil
		user.Conf(headers.verbose)
		if user.config != headers.verbose {
			t.Errorf("failed to apply config")
		}
	})

	t.Run("Modifiers", func(t *testing.T) {
		user := sp.Web().CurrentUser()
		mods := user.Select("*").Expand("*").modifiers
		if mods == nil || len(mods.mods) != 2 {
			t.Error("can't add modifiers")
		}
	})

	t.Run("GetUserInfo", func(t *testing.T) {
		data, err := user.Get()
		if err != nil {
			t.Error(err)
		}

		if data.Data().ID == 0 {
			t.Error("can't get user info")
		}
	})

	t.Run("GetGroups", func(t *testing.T) {
		if _, err := user.Groups().Select("Id").Get(); err != nil {
			t.Error(err)
		}
	})

	// ToDo:
	// Update

}
