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

	t.Run("Conf", func(t *testing.T) {
		g := web.AssociatedGroups().Visitors()
		hs := map[string]*RequestConfig{
			"nometadata":      HeadersPresets.Nometadata,
			"minimalmetadata": HeadersPresets.Minimalmetadata,
			"verbose":         HeadersPresets.Verbose,
		}
		for key, preset := range hs {
			u := g.Conf(preset)
			if u.config != preset {
				t.Errorf("can't %v config", key)
			}
		}
	})

	t.Run("Modifiers", func(t *testing.T) {
		g := web.AssociatedGroups().Visitors()
		mods := g.Select("*").Expand("*").modifiers
		if mods == nil || len(mods.mods) != 2 {
			t.Error("wrong number of modifiers")
		}
	})

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
