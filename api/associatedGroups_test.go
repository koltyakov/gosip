package api

import (
	"testing"
)

func TestAssociatedGroups(t *testing.T) {
	checkClient(t)

	sp := NewSP(spClient)

	t.Run("Visitors", func(t *testing.T) {
		group, err := sp.Web().AssociatedGroups().Visitors().Get()
		if err != nil {
			t.Error(err)
		}
		if group.Data().ID == 0 {
			t.Error("can't get visitors group")
		}
	})

	t.Run("Members", func(t *testing.T) {
		group, err := sp.Web().AssociatedGroups().Members().Get()
		if err != nil {
			t.Error(err)
		}
		if group.Data().ID == 0 {
			t.Error("can't get members group")
		}
	})

	t.Run("Owners", func(t *testing.T) {
		group, err := sp.Web().AssociatedGroups().Owners().Get()
		if err != nil {
			t.Error(err)
		}
		if group.Data().ID == 0 {
			t.Error("can't get owners group")
		}
	})

}
