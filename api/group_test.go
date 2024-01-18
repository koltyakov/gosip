package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestGroup(t *testing.T) {
	checkClient(t)

	web := NewSP(spClient).Web()
	newGroupName := uuid.New().String()
	group := &GroupInfo{}

	user, err := web.CurrentUser().Select("Id,LoginName").Get(context.Background())
	if err != nil {
		t.Error(err)
	}

	type groupOwner struct {
		Owner struct {
			ID int
		}
	}

	t.Run("Add", func(t *testing.T) {
		data, err := web.SiteGroups().Add(context.Background(), newGroupName, nil)
		if err != nil {
			t.Error(err)
		}
		group = data.Data()
	})

	t.Run("Get", func(t *testing.T) {
		data, err := web.SiteGroups().GetByName(newGroupName).Get(context.Background())
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
		if _, err := web.SiteGroups().GetByID(group.ID).Update(context.Background(), body); err != nil {
			t.Error(err)
		}
	})

	t.Run("AddUser", func(t *testing.T) {
		if err := web.SiteGroups().GetByID(group.ID).AddUser(context.Background(), user.Data().LoginName); err != nil {
			t.Error(err)
		}
	})

	t.Run("RemoveUser", func(t *testing.T) {
		if err := web.SiteGroups().GetByID(group.ID).RemoveUser(context.Background(), user.Data().LoginName); err != nil {
			t.Error(err)
		}
	})

	t.Run("AddUserByID", func(t *testing.T) {
		if err := web.SiteGroups().GetByID(group.ID).AddUserByID(context.Background(), user.Data().ID); err != nil {
			t.Error(err)
		}
	})

	t.Run("RemoveUserByID", func(t *testing.T) {
		if err := web.SiteGroups().GetByID(group.ID).RemoveUserByID(context.Background(), user.Data().ID); err != nil {
			t.Error(err)
		}
	})

	t.Run("SetAsOwner/User", func(t *testing.T) {
		au, err := web.SiteUsers().Select("Id").Filter(fmt.Sprintf("Id ne %d and PrincipalType eq 1", user.Data().ID)).Top(1).Get(context.Background())
		if err != nil {
			t.Error(err)
		}
		if len(au) == 0 {
			return
		}
		g := web.SiteGroups().GetByID(group.ID)
		if err := g.SetOwner(context.Background(), au.Data()[0].Data().ID); err != nil {
			t.Error(err)
		}
		o, err := g.Select("Owner/Id").Expand("Owner").Get(context.Background())
		if err != nil {
			t.Error(err)
		}
		var owner *groupOwner
		if err := json.Unmarshal(o.Normalized(), &owner); err != nil {
			t.Error(err)
		}
		if owner.Owner.ID == user.Data().ID {
			t.Error("can't set a user as group owner")
		}
	})

	t.Run("SetAsOwner/Group", func(t *testing.T) {
		mg, err := web.AssociatedGroups().Members().Select("Id").Get(context.Background())
		if err != nil {
			t.Error(err)
		}
		g := web.SiteGroups().GetByID(group.ID)
		if err := g.SetOwner(context.Background(), mg.Data().ID); err != nil {
			t.Error(err)
		}
		o, err := g.Select("Owner/Id").Expand("Owner").Get(context.Background())
		if err != nil {
			t.Error(err)
		}
		var owner *groupOwner
		if err := json.Unmarshal(o.Normalized(), &owner); err != nil {
			t.Error(err)
		}
		if owner.Owner.ID == user.Data().ID {
			t.Error("can't set a user as group owner")
		}
	})

	t.Run("SetUserAsOwner", func(t *testing.T) {
		if envCode == "2013" {
			t.Skip("is not supported with SP 2013")
		}

		au, err := web.SiteUsers().Select("Id").Filter(fmt.Sprintf("Id ne %d", user.Data().ID)).OrderBy("Id", false).Top(1).Get(context.Background())
		if err != nil {
			t.Error(err)
		}
		if len(au) == 0 {
			return
		}
		g := web.SiteGroups().GetByID(group.ID)
		if err := g.SetUserAsOwner(context.Background(), au.Data()[0].Data().ID); err != nil {
			t.Error(err)
		}
	})

	t.Run("RemoveByID", func(t *testing.T) {
		if err := web.SiteGroups().RemoveByID(context.Background(), group.ID); err != nil {
			t.Error(err)
		}
	})

	t.Run("Modifiers", func(t *testing.T) {
		if _, err := web.AssociatedGroups().Visitors().Users().Get(context.Background()); err != nil {
			t.Error(err)
		}
	})

}
