package api

import (
	"encoding/json"
	"testing"
)

func TestUser(t *testing.T) {
	checkClient(t)

	user := NewSP(spClient).Web().CurrentUser()
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

	t.Run("GetUserInfo", func(t *testing.T) {
		data, err := user.Conf(headers.verbose).Get()

		if err != nil {
			t.Error(err)
		}

		res := &struct {
			User *UserInfo `json:"d"`
		}{}

		if err := json.Unmarshal(data, &res); err != nil {
			t.Error(err)
		}

		if res.User.ID == 0 {
			t.Error("can't get user info")
		}
	})

	t.Run("GetGroups", func(t *testing.T) {
		data, err := user.Groups().Select("Id").Conf(headers.verbose).Get()

		if err != nil {
			t.Error(err)
		}

		res := &struct {
			D struct {
				Results []*GroupInfo `json:"results"`
			} `json:"d"`
		}{}

		if err := json.Unmarshal(data, &res); err != nil {
			t.Error(err)
		}
	})

}
