package api

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
)

func TestRoles(t *testing.T) {
	checkClient(t)

	web := NewSP(spClient).Web()
	newListTitle := uuid.New().String()

	// Pre-configuration
	if _, err := web.Lists().Add(context.Background(), newListTitle, nil); err != nil {
		t.Errorf("can't create a list to test permissions: %s", err)
	}
	list := web.Lists().GetByTitle(newListTitle)
	userID, err := getCurrentUserID()
	if err != nil {
		t.Error(err)
	}
	roleDef, err := web.RoleDefinitions().GetByType(context.Background(), RoleTypeKinds.Contributor)
	if err != nil {
		t.Error(err)
	}

	t.Run("BreakInheritance", func(t *testing.T) {
		if err := list.Roles().BreakInheritance(context.Background(), false, true); err != nil {
			t.Error(err)
		}
	})

	t.Run("HasUniqueAssignments", func(t *testing.T) {
		_ = list.Roles().BreakInheritance(context.Background(), false, true)
		hasUniqueAssignments, err := list.Roles().HasUniqueAssignments(context.Background())
		if err != nil {
			t.Error(err)
		}
		if !hasUniqueAssignments {
			t.Error("should have had unique role assignments")
		}
	})

	t.Run("AddAssigment", func(t *testing.T) {
		if err := list.Roles().AddAssigment(context.Background(), userID, roleDef.ID); err != nil {
			t.Error(err)
		}
	})

	t.Run("RemoveAssigment", func(t *testing.T) {
		if err := list.Roles().RemoveAssigment(context.Background(), userID, roleDef.ID); err != nil {
			t.Error(err)
		}
	})

	t.Run("ResetInheritance", func(t *testing.T) {
		if err := list.Roles().ResetInheritance(context.Background()); err != nil {
			t.Error(err)
		}
	})

	// Post-configurations
	if err := list.Delete(context.Background()); err != nil {
		t.Errorf("can't delete a list: %s", err)
	}

}

func getCurrentUserID() (int, error) {
	web := NewSP(spClient).Web()
	data, err := web.CurrentUser().Select("Id").Conf(headers.verbose).Get(context.Background())
	if err != nil {
		return 0, err
	}
	res := &struct {
		User *UserInfo `json:"d"`
	}{}
	if err := json.Unmarshal(data, &res); err != nil {
		return 0, err
	}
	return res.User.ID, nil
}
