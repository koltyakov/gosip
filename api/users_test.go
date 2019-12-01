package api

import (
	"encoding/json"
	"testing"
)

func TestUsers(t *testing.T) {
	checkClient(t)

	users := NewSP(spClient).Web().SiteUsers()
	endpoint := spClient.AuthCnfg.GetSiteURL() + "/_api/Web/SiteUsers"
	user := &UserInfo{}

	t.Run("Constructor", func(t *testing.T) {
		users := NewUsers(spClient, endpoint, nil)
		if _, err := users.Select("Id").Get(); err != nil {
			t.Error(err)
		}
	})

	t.Run("ToURL", func(t *testing.T) {
		if users.ToURL() != endpoint {
			t.Errorf(
				"incorrect endpoint URL, expected \"%s\", received \"%s\"",
				endpoint,
				users.ToURL(),
			)
		}
	})

	t.Run("Conf", func(t *testing.T) {
		users.config = nil
		users.Conf(headers.verbose)
		if users.config != headers.verbose {
			t.Errorf("failed to apply config")
		}
	})

	t.Run("GetUsers", func(t *testing.T) {
		data, err := users.Select("Id").Top(5).Get()
		if err != nil {
			t.Error(err)
		}

		res := &struct {
			D struct {
				Results []UserInfo `json:"results"`
			} `json:"d"`
		}{}

		if err := json.Unmarshal(data, &res); err != nil {
			t.Error(err)
		}

		if len(res.D.Results) == 0 {
			t.Error("can't get users")
		}
	})

	t.Run("GetUser", func(t *testing.T) {
		data, err := NewSP(spClient).Web().CurrentUser().Conf(headers.verbose).Get()

		if err != nil {
			t.Error(err)
		}

		res := &struct {
			User *UserInfo `json:"d"`
		}{}

		if err := json.Unmarshal(data, &res); err != nil {
			t.Error(err)
		}

		user = res.User
	})

	t.Run("GetByID", func(t *testing.T) {
		if user.ID == 0 {
			t.Skip("no user ID to use in the test")
		}

		data, err := users.GetByID(user.ID).Select("Id").Get()
		if err != nil {
			t.Error(err)
		}

		res := &struct {
			User *UserInfo `json:"d"`
		}{}

		if err := json.Unmarshal(data, &res); err != nil {
			t.Error(err)
		}

		if res.User.ID != user.ID {
			t.Errorf(
				"incorrect user ID, expected \"%d\", received \"%d\"",
				user.ID,
				res.User.ID,
			)
		}
	})

	t.Run("GetByLoginName", func(t *testing.T) {
		if envCode == "2013" {
			t.Skip("is not supported with SP 2013")
		}
		if user.LoginName == "" {
			t.Skip("no user LoginName to use in the test")
		}

		data, err := users.GetByLoginName(user.LoginName).Select("LoginName").Get()
		if err != nil {
			t.Error(err)
		}

		res := &struct {
			User *UserInfo `json:"d"`
		}{}

		if err := json.Unmarshal(data, &res); err != nil {
			t.Error(err)
		}

		if res.User.LoginName != user.LoginName {
			t.Errorf(
				"incorrect user LoginName, expected \"%s\", received \"%s\"",
				user.LoginName,
				res.User.LoginName,
			)
		}
	})

	t.Run("GetByEmail", func(t *testing.T) {
		if user.Email == "" {
			t.Skip("no user Email to use in the test")
		}

		data, err := users.GetByEmail(user.Email).Select("Email").Get()
		if err != nil {
			t.Error(err)
		}

		res := &struct {
			User *UserInfo `json:"d"`
		}{}

		if err := json.Unmarshal(data, &res); err != nil {
			t.Error(err)
		}

		if res.User.Email != user.Email {
			t.Errorf(
				"incorrect user Email, expected \"%s\", received \"%s\"",
				user.Email,
				res.User.Email,
			)
		}
	})

}
