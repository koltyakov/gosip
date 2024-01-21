package api

import (
	"context"
	"testing"
)

func TestRoleDefinitions(t *testing.T) {
	checkClient(t)

	web := NewSP(spClient).Web()
	roleDef := &RoleDefInfo{}

	t.Run("GetByType", func(t *testing.T) {
		res, err := web.RoleDefinitions().GetByType(context.Background(), RoleTypeKinds.Contributor)
		if err != nil {
			t.Error(err)
		}
		if res.Name == "" {
			t.Error("can't find role definition")
		}
		roleDef = res
	})

	t.Run("GetByID", func(t *testing.T) {
		if roleDef.ID == 0 {
			t.Skip("no role definition ID provided")
		}

		res, err := web.RoleDefinitions().GetByID(context.Background(), roleDef.ID)
		if err != nil {
			t.Error(err)
		}
		if res.Name == "" {
			t.Error("can't find role definition")
		}
		roleDef = res
	})

	t.Run("GetByName", func(t *testing.T) {
		if roleDef.Name == "" {
			t.Skip("no role definition Name provided")
		}

		res, err := web.RoleDefinitions().GetByName(context.Background(), roleDef.Name)
		if err != nil {
			t.Error(err)
		}
		if res.Name == "" {
			t.Error("can't find role definition")
		}
		roleDef = res
	})

	t.Run("Get", func(t *testing.T) {
		res, err := web.RoleDefinitions().Get(context.Background())
		if err != nil {
			t.Error(err)
		}
		if len(res) == 0 {
			t.Error("can't get role definitions")
		}
	})

}
