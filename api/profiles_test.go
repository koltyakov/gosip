package api

import (
	"bytes"
	"testing"
)

func TestProfiles(t *testing.T) {
	checkClient(t)

	sp := NewSP(spClient)
	profiles := sp.Profiles()
	user, err := sp.Web().CurrentUser().Get()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("GetMyProperties", func(t *testing.T) {
		profile, err := profiles.GetMyProperties()
		if err != nil {
			t.Fatal(err)
		}
		if len(profile.Data().UserProfileProperties) == 0 {
			t.Error("can't get user profile properties")
		}
	})

	t.Run("GetPropertiesFor", func(t *testing.T) {
		props, err := profiles.GetPropertiesFor(user.Data().LoginName)
		if err != nil {
			t.Fatal(err)
		}
		if len(props.Data().UserProfileProperties) == 0 {
			t.Error("can't get user profile properties")
		}
		if bytes.Compare(props, props.Normalized()) == -1 {
			t.Error("wrong response normalization")
		}
	})

	t.Run("GetUserProfilePropertyFor", func(t *testing.T) {
		accountName, err := sp.Profiles().Conf(headers.verbose).
			GetUserProfilePropertyFor(user.Data().LoginName, "AccountName")
		if err != nil {
			t.Fatal(err)
		}
		if accountName == "" {
			t.Error("wrong property value")
		}
		if envCode != "2013" {
			accountName, err = sp.Profiles().Conf(headers.minimalmetadata).
				GetUserProfilePropertyFor(user.Data().LoginName, "AccountName")
			if err != nil {
				t.Fatal(err)
			}
			if accountName == "" {
				t.Error("wrong property value")
			}
			accountName, err = sp.Profiles().Conf(headers.nometadata).
				GetUserProfilePropertyFor(user.Data().LoginName, "AccountName")
			if err != nil {
				t.Fatal(err)
			}
			if accountName == "" {
				t.Error("wrong property value")
			}
		}
	})

	t.Run("GetOwnerUserProfile", func(t *testing.T) {
		profile, err := profiles.GetOwnerUserProfile()
		if err != nil {
			t.Fatal(err)
		}
		if profile.Data().AccountName == "" {
			t.Error("can't get profile")
		}
	})

	t.Run("UserProfile", func(t *testing.T) {
		profile, err := profiles.UserProfile()
		if err != nil {
			t.Fatal(err)
		}
		if profile.Data().AccountName == "" {
			t.Error("can't get profile")
		}
		if bytes.Compare(profile, profile.Normalized()) == -1 {
			t.Error("wrong response normalization")
		}
	})

	t.Run("UserProfile", func(t *testing.T) {
		profile, err := profiles.UserProfile()
		if err != nil {
			t.Fatal(err)
		}
		if _, err := sp.Profiles().HideSuggestion(profile.Data().AccountName); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("SetSingleValueProfileProperty", func(t *testing.T) {
		if err := profiles.SetSingleValueProfileProperty(user.Data().LoginName, "AboutMe", "Updated from Gosip"); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("SetMultiValuedProfileProperty", func(t *testing.T) {
		tags := []string{"#ci", "#demo", "#test"}
		if err := profiles.SetMultiValuedProfileProperty(user.Data().LoginName, "SPS-HashTags", tags); err != nil {
			t.Fatal(err)
		}
	})

}
