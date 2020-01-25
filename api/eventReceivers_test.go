package api

import (
	"testing"
)

func TestEventReceivers(t *testing.T) {
	checkClient(t)

	sp := NewSP(spClient)

	t.Run("Get/Site", func(t *testing.T) {
		receivers, err := sp.Site().EventReceivers().Top(1).Get()
		if err != nil {
			t.Error(err)
		}
		if receivers[0].ReceiverID == "" {
			t.Error("can't get event receivers")
		}
	})

	t.Run("Get/Web", func(t *testing.T) {
		receivers, err := sp.Web().EventReceivers().Top(1).Get()
		if err != nil {
			t.Error(err)
		}
		if receivers[0].ReceiverID == "" {
			t.Error("can't get event receivers")
		}
	})

	t.Run("Modifiers", func(t *testing.T) {
		er := sp.Web().EventReceivers()
		mods := er.Select("*").Filter("*").Top(1).OrderBy("*", true).modifiers
		if mods == nil || len(mods.mods) != 4 {
			t.Error("wrong number of modifiers")
		}
	})

}
