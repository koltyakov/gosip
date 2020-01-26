package csom

import (
	"testing"
)

func TestCurrent(t *testing.T) {
	c := &current{}

	t.Run("String", func(t *testing.T) {
		static := `<StaticProperty Id="0" TypeId="{3747adcd-a3c3-41b9-bfab-4a64dd2f1e0a}" Name="Current" />`
		if c.String() != static {
			t.Error("wrong static property")
		}
	})

	t.Run("ID", func(t *testing.T) {
		if c.GetID() != 0 {
			t.Error("wrong ID")
		}
		c.SetID(1)
		if c.GetID() != 0 {
			t.Error("wrong ID, SetID should be ignored")
		}
	})

	t.Run("ParentID", func(t *testing.T) {
		if c.GetParentID() != -1 {
			t.Error("wrong ParentID")
		}
		c.SetParentID(1)
		if c.GetParentID() != -1 {
			t.Error("wrong ParentID, SetParentID should be ignored")
		}
	})

}
