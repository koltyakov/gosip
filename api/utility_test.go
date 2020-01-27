package api

import (
	"testing"
)

func TestUtility(t *testing.T) {
	t.Parallel()
	checkClient(t)

	sp := NewSP(spClient)

	t.Run("SendEmail", func(t *testing.T) {
		if spClient.AuthCnfg.GetStrategy() == "addin" {
			t.Skip("not supported by addin auth")
		}
		user, err := sp.Web().CurrentUser().Get()
		if err != nil {
			t.Error(err)
		}
		if user.Data().Email != "" {
			if err := sp.Utility().SendEmail(&EmailProps{
				Subject: "Gosip SendEmail utility test",
				Body:    "Feel free to delete the email",
				To:      []string{user.Data().Email},
			}); err != nil {
				t.Error(err)
			}
		}
	})

	t.Run("SendEmail", func(t *testing.T) {
		if spClient.AuthCnfg.GetStrategy() == "addin" {
			t.Skip("not supported by addin auth")
		}
		user, err := sp.Web().CurrentUser().Get()
		if err != nil {
			t.Error(err)
		}
		if user.Data().Email != "" {
			if err := sp.Utility().SendEmail(&EmailProps{
				Subject: "Gosip SendEmail utility test",
				Body:    "Feel free to delete the email",
				To:      []string{user.Data().Email},
				From:    user.Data().Email,
				CC:      []string{user.Data().Email},
				BCC:     []string{user.Data().Email},
			}); err != nil {
				t.Error(err)
			}
		}
	})

}
