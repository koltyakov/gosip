package api

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
)

func TestGroup(t *testing.T) {
	t.Parallel()
	checkClient(t)

	web := NewSP(spClient).Web()
	newGroupName := uuid.New().String()
	group := &GroupInfo{}

	t.Run("Add", func(t *testing.T) {
		data, err := web.SiteGroups().Add(newGroupName, nil)
		if err != nil {
			t.Error(err)
		}
		group = data.Data()
	})

	t.Run("Get", func(t *testing.T) {
		data, err := web.SiteGroups().GetByName(newGroupName).Get()
		if err != nil {
			t.Error(err)
		}
		if bytes.Compare(data, data.Normalized()) == -1 {
			t.Error("response normalization error")
		}
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
		user, err := web.CurrentUser().Select("LoginName").Get()
		if err != nil {
			t.Error(err)
		}
		if err := web.SiteGroups().GetByID(group.ID).AddUser(user.Data().LoginName); err != nil {
			t.Error(err)
		}
	})

	t.Run("RemoveUser", func(t *testing.T) {
		user, err := web.CurrentUser().Select("LoginName").Get()
		if err != nil {
			t.Error(err)
		}
		if err := web.SiteGroups().GetByID(group.ID).RemoveUser(user.Data().LoginName); err != nil {
			t.Error(err)
		}
	})

	t.Run("AddUserByID", func(t *testing.T) {
		user, err := web.CurrentUser().Select("Id").Get()
		if err != nil {
			t.Error(err)
		}
		if err := web.SiteGroups().GetByID(group.ID).AddUserByID(user.Data().ID); err != nil {
			t.Error(err)
		}
	})

	t.Run("RemoveUserByID", func(t *testing.T) {
		user, err := web.CurrentUser().Select("Id").Get()
		if err != nil {
			t.Error(err)
		}
		if err := web.SiteGroups().GetByID(group.ID).RemoveUserByID(user.Data().ID); err != nil {
			t.Error(err)
		}
	})

	t.Run("RemoveByID", func(t *testing.T) {
		if err := web.SiteGroups().RemoveByID(group.ID); err != nil {
			t.Error(err)
		}
	})

	t.Run("Modifiers", func(t *testing.T) {
		if _, err := web.AssociatedGroups().Visitors().Users().Get(); err != nil {
			t.Error(err)
		}
	})

}
