package api

import (
	"testing"
)

func TestFeatures(t *testing.T) {
	checkClient(t)

	sp := NewSP(spClient)
	mdsFeatureID := "87294c72-f260-42f3-a41b-981a2ffce37a"

	sp.Web().Features().Remove(mdsFeatureID, true)

	t.Run("GetSiteFeatures", func(t *testing.T) {
		if _, err := sp.Site().Features().Get(); err != nil {
			t.Error(err)
		}
	})

	t.Run("Get", func(t *testing.T) {
		if _, err := sp.Web().Features().Get(); err != nil {
			t.Error(err)
		}
	})

	t.Run("Add", func(t *testing.T) {
		if _, err := sp.Web().Features().Add(mdsFeatureID, true); err != nil {
			t.Error(err)
		}
	})

	t.Run("Remove", func(t *testing.T) {
		if _, err := sp.Web().Features().Remove(mdsFeatureID, true); err != nil {
			t.Error(err)
		}
	})

}
