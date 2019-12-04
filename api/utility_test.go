package api

import (
	"testing"
)

func TestUtility(t *testing.T) {

	sp := NewSP(spClient)

	t.Run("SendEmail", func(t *testing.T) {
		user, err := sp.Web().CurrentUser().Get()
		if err != nil {
			t.Error(err)
		}
		if user.Data().Email != "" {
			if _, err := sp.Utility().SendEmail(&EmailProps{
				Subject: "Gosip SendEmail utility test",
				Body:    "Feel free to delete the email",
				To:      []string{user.Data().Email},
			}); err != nil {
				t.Error(err)
			}
		}
	})

}
