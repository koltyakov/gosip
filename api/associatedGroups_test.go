package api

import (
	"context"
	"testing"
)

func TestAssociatedGroups(t *testing.T) {
	checkClient(t)

	sp := NewSP(spClient)

	t.Run("Visitors", func(t *testing.T) {
		group, err := sp.Web().AssociatedGroups().Visitors().Get(context.Background())
		if err != nil {
			t.Error(err)
		}
		if group.Data().ID == 0 {
			t.Error("can't get visitors group")
		}
	})

	t.Run("Members", func(t *testing.T) {
		group, err := sp.Web().AssociatedGroups().Members().Get(context.Background())
		if err != nil {
			t.Error(err)
		}
		if group.Data().ID == 0 {
			t.Error("can't get members group")
		}
	})

	t.Run("Owners", func(t *testing.T) {
		group, err := sp.Web().AssociatedGroups().Owners().Get(context.Background())
		if err != nil {
			t.Error(err)
		}
		if group.Data().ID == 0 {
			t.Error("can't get owners group")
		}
	})

}
