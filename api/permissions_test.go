package api

import (
	"testing"
)

func TestPermissions(t *testing.T) {
	// checkClient(t)

	editPermissions := BasePermissions{
		High: 432,
		Low:  1011030767,
	}

	limitedPermissions := BasePermissions{
		High: 48,
		Low:  134287360,
	}

	t.Run("HasPermissions", func(t *testing.T) {

		if has := HasPermissions(editPermissions, PermissionKind.EditListItems); !has {
			t.Error("should have Edit permissions")
		}

		if has := HasPermissions(limitedPermissions, PermissionKind.ViewListItems); has {
			t.Error("should not have View permissions")
		}

	})

}
