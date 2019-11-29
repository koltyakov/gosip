package api

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
)

func TestGroups(t *testing.T) {
	checkClient(t)

	groups := NewSP(spClient).Web().SiteGroups()
	endpoint := spClient.AuthCnfg.GetSiteURL() + "/_api/Web/SiteGroups"
	newGroupName := uuid.New().String()
	group := &GroupInfo{}

	t.Run("Constructor", func(t *testing.T) {
		groups := NewGroups(spClient, endpoint, nil)
		if _, err := groups.Select("Id").Get(); err != nil {
			t.Error(err)
		}
	})

	t.Run("ToURL", func(t *testing.T) {
		if groups.ToURL() != endpoint {
			t.Errorf(
				"incorrect endpoint URL, expected \"%s\", received \"%s\"",
				endpoint,
				groups.ToURL(),
			)
		}
	})

	t.Run("Conf", func(t *testing.T) {
		groups.config = nil
		groups.Conf(headers.verbose)
		if groups.config != headers.verbose {
			t.Errorf("failed to apply config")
		}
	})

	t.Run("GetGroups", func(t *testing.T) {
		data, err := groups.Select("Id").Top(5).Get()
		if err != nil {
			t.Error(err)
		}

		res := &struct {
			D struct {
				Results []interface{} `json:"results"`
			} `json:"d"`
		}{}

		if err := json.Unmarshal(data, &res); err != nil {
			t.Error(err)
		}

		if len(res.D.Results) == 0 {
			t.Error("can't get users")
		}
	})

	t.Run("GetGroup", func(t *testing.T) {
		data, err := groups.Select("Id,Title").Top(1).Conf(headers.verbose).Get()

		if err != nil {
			t.Error(err)
		}

		res := &struct {
			D struct {
				Groups []*GroupInfo `json:"results"`
			} `json:"d"`
		}{}

		if err := json.Unmarshal(data, &res); err != nil {
			t.Error(err)
		}

		group = res.D.Groups[0]
	})

	t.Run("GetByID", func(t *testing.T) {
		if group.ID == 0 {
			t.Skip("no group ID to use in the test")
		}

		data, err := groups.GetByID(group.ID).Select("Id").Get()
		if err != nil {
			t.Error(err)
		}

		res := &struct {
			Group *GroupInfo `json:"d"`
		}{}

		if err := json.Unmarshal(data, &res); err != nil {
			t.Error(err)
		}

		if res.Group.ID != group.ID {
			t.Errorf(
				"incorrect group ID, expected \"%d\", received \"%d\"",
				group.ID,
				res.Group.ID,
			)
		}
	})

	t.Run("GetByName", func(t *testing.T) {
		if group.Title == "" {
			t.Skip("no group Title to use in the test")
		}

		data, err := groups.GetByName(group.Title).Select("Title").Get()
		if err != nil {
			t.Error(err)
		}

		res := &struct {
			Group *GroupInfo `json:"d"`
		}{}

		if err := json.Unmarshal(data, &res); err != nil {
			t.Error(err)
		}

		if res.Group.Title != group.Title {
			t.Errorf(
				"incorrect group Title, expected \"%s\", received \"%s\"",
				group.Title,
				res.Group.Title,
			)
		}
	})

	t.Run("Add", func(t *testing.T) {
		data, err := groups.Conf(headers.verbose).Add(newGroupName, nil)
		if err != nil {
			t.Error(err)
		}

		res := &struct {
			Group *GroupInfo `json:"d"`
		}{}

		if err := json.Unmarshal(data, &res); err != nil {
			t.Error(err)
		}
		group = res.Group
	})

	t.Run("RemoveByLoginName", func(t *testing.T) {
		if _, err := groups.RemoveByLoginName(group.LoginName); err != nil {
			t.Error(err)
		}
	})

}
