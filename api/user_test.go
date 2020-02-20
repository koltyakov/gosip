package api

import (
	"bytes"
	"testing"
)

func TestUser(t *testing.T) {
	checkClient(t)

	sp := NewSP(spClient)
	user := sp.Web().CurrentUser()
	endpoint := spClient.AuthCnfg.GetSiteURL() + "/_api/Web/CurrentUser"

	t.Run("Constructor", func(t *testing.T) {
		u := NewUser(spClient, endpoint, nil)
		if _, err := u.Select("Id").Get(); err != nil {
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

	t.Run("GetUserInfo", func(t *testing.T) {
		data, err := user.Get()
		if err != nil {
			t.Error(err)
		}

		if data.Data().ID == 0 {
			t.Error("can't get user info")
		}

		if bytes.Compare(data, data.Normalized()) == -1 {
			t.Error("wrong response normalization")
		}
	})

	t.Run("GetGroups", func(t *testing.T) {
		if _, err := user.Groups().Select("Id").Get(); err != nil {
			t.Error(err)
		}
	})

}
