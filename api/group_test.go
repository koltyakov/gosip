package api

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
)

func TestGroup(t *testing.T) {
	checkClient(t)

	web := NewSP(spClient).Web()
	newGroupName := uuid.New().String()
	group := &GroupInfo{}

	t.Run("Add", func(t *testing.T) {
		data, err := web.SiteGroups().Conf(headers.verbose).Add(newGroupName, nil)
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

	t.Run("Update", func(t *testing.T) {
		metadata := make(map[string]interface{})
		metadata["__metadata"] = map[string]string{
			"type": "SP.Group",
		}
		metadata["Description"] = "It's a test group" // ToDo: check if update works
		body, _ := json.Marshal(metadata)
		if _, err := web.SiteGroups().GetByID(group.ID).Update(body); err != nil {
			t.Error(err)
		}
	})

	t.Run("AddUser", func(t *testing.T) {
		data, err := web.CurrentUser().Select("LoginName").Conf(headers.verbose).Get()
		if err != nil {
			t.Error(err)
		}

		res := &struct {
			User *UserInfo `json:"d"`
		}{}

		if err := json.Unmarshal(data, &res); err != nil {
			t.Error(err)
		}

		if _, err := web.SiteGroups().GetByID(group.ID).AddUser(res.User.LoginName); err != nil {
			t.Error(err)
		}
	})

	t.Run("RemoveByID", func(t *testing.T) {
		if _, err := web.SiteGroups().RemoveByID(group.ID); err != nil {
			t.Error(err)
		}
	})

}
