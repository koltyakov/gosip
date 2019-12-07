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
		profile, err := profiles.GetMyProperties()
		if err != nil {
			t.Error(err)
		}
		if len(profile.Data().UserProfileProperties) == 0 {
			t.Error("can't get user profile properties")
		}
	})

	t.Run("GetPropertiesFor", func(t *testing.T) {
		profile, err := profiles.GetPropertiesFor(user.Data().LoginName)
		if err != nil {
			t.Error(err)
		}
		if len(profile.Data().UserProfileProperties) == 0 {
			t.Error("can't get user profile properties")
		}
	})

	t.Run("GetUserProfilePropertyFor", func(t *testing.T) {
		if _, err := profiles.GetUserProfilePropertyFor(user.Data().LoginName, "LastName"); err != nil {
			t.Error(err)
		}
	})

	t.Run("GetOwnerUserProfile", func(t *testing.T) {
		profile, err := profiles.GetOwnerUserProfile()
		if err != nil {
			t.Error(err)
		}
		if profile.Data().AccountName == "" {
			t.Error("can't get profile")
		}
	})

	t.Run("UserProfile", func(t *testing.T) {
		profile, err := profiles.UserProfile()
		if err != nil {
			t.Error(err)
		}
		if profile.Data().AccountName == "" {
			t.Error("can't get profile")
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
