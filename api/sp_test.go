package api

import (
	"context"
	"testing"
)

func TestSP(t *testing.T) {
	checkClient(t)

	t.Run("ToURL", func(t *testing.T) {
		sp := NewSP(spClient)
		if sp.ToURL() != spClient.AuthCnfg.GetSiteURL() {
			t.Errorf(
				"incorrect site URL, expected \"%s\", received \"%s\"",
				spClient.AuthCnfg.GetSiteURL(),
				sp.ToURL(),
			)
		}
	})

	t.Run("Web", func(t *testing.T) {
		sp := NewSP(spClient)
		if sp.Web() == nil {
			t.Errorf("failed to get Web object")
		}
	})

	t.Run("Site", func(t *testing.T) {
		sp := NewSP(spClient)
		if sp.Site() == nil {
			t.Errorf("failed to get Site object")
		}
	})

	t.Run("Utility", func(t *testing.T) {
		sp := NewSP(spClient)
		if sp.Utility() == nil {
			t.Errorf("failed to get Utility object")
		}
	})

	t.Run("Search", func(t *testing.T) {
		sp := NewSP(spClient)
		if sp.Search() == nil {
			t.Errorf("failed to get Search object")
		}
	})

	t.Run("Profiles", func(t *testing.T) {
		sp := NewSP(spClient)
		if sp.Profiles() == nil {
			t.Errorf("failed to get Profiles object")
		}
	})

	t.Run("Taxonomy", func(t *testing.T) {
		sp := NewSP(spClient)
		if sp.Taxonomy() == nil {
			t.Errorf("failed to get Taxonomy object")
		}
	})

	t.Run("Metadata", func(t *testing.T) {
		sp := NewSP(spClient)
		if _, err := sp.Metadata(context.Background()); err != nil {
			t.Error(err)
		}
	})

}
