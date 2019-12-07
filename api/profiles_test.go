package api

import (
	"testing"
)

func TestProfiles(t *testing.T) {
	checkClient(t)

	sp := NewSP(spClient)
	profiles := sp.Profiles()
	user, err := sp.Web().CurrentUser().Get()
	if err != nil {
		t.Error(err)
	}

	t.Run("Get", func(t *testing.T) {
		if _, err := profiles.Get(); err != nil {
			t.Error(err)
		}
	})

	t.Run("GetMyProperties", func(t *testing.T) {
		if _, err := profiles.GetMyProperties(); err != nil {
			t.Error(err)
		}
	})

	t.Run("GetPropertiesFor", func(t *testing.T) {
		if _, err := profiles.GetPropertiesFor(user.Data().LoginName); err != nil {
			t.Error(err)
		}
	})

	t.Run("GetUserProfilePropertyFor", func(t *testing.T) {
		if _, err := profiles.GetUserProfilePropertyFor(user.Data().LoginName, "LastName"); err != nil {
			t.Error(err)
		}
	})

	t.Run("GetOwnerUserProfile", func(t *testing.T) {
		if _, err := profiles.GetOwnerUserProfile(); err != nil {
			t.Error(err)
		}
	})

	t.Run("UserProfile", func(t *testing.T) {
		if _, err := profiles.UserProfile(); err != nil {
			t.Error(err)
		}
	})

	t.Run("SetSingleValueProfileProperty", func(t *testing.T) {
		if _, err := profiles.SetSingleValueProfileProperty(user.Data().LoginName, "AboutMe", "Updated from Gosip"); err != nil {
			t.Error(err)
		}
	})

	t.Run("SetMultiValuedProfileProperty", func(t *testing.T) {
		tags := []string{"#ci", "#demo", "#test"}
		if _, err := profiles.SetMultiValuedProfileProperty(user.Data().LoginName, "SPS-HashTags", tags); err != nil {
			t.Error(err)
		}
	})

}
