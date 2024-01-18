package api

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestSubscriptions(t *testing.T) {
	checkClient(t)

	if envCode == "2013" || envCode == "2016" {
		t.Skip("is not supported with SP " + envCode)
	}

	notificationURL := os.Getenv("GOSIP_TESTS_WEBHOOKS_URL")
	if len(notificationURL) == 0 {
		t.Skip(`provide "GOSIP_TESTS_WEBHOOKS_URL" environment variable to enable these tests`)
	}

	web := NewSP(spClient).Web()
	list := web.Lists().GetByTitle("Site Pages")
	subID := ""

	t.Run("Add", func(t *testing.T) {
		expiration := time.Now().Add(60 * time.Second)
		sub, err := list.Subscriptions().Add(context.Background(), notificationURL, expiration, "")
		if err != nil {
			t.Error(err)
			return
		}
		if sub.ID == "" {
			t.Error("can't parse subscription add response")
			return
		}
		subID = sub.ID
	})

	t.Run("Add#ExpirationInThePast", func(t *testing.T) {
		expiration := time.Now().AddDate(0, 0, -1)
		if _, err := list.Subscriptions().Add(context.Background(), notificationURL, expiration, ""); err == nil {
			t.Error("should fail due to expiration limitation 'not in the past'")
		}
	})

	t.Run("Add#MoreThan6Months", func(t *testing.T) {
		expiration := time.Now().AddDate(0, 6, 10)
		if _, err := list.Subscriptions().Add(context.Background(), notificationURL, expiration, ""); err == nil {
			t.Error("should fail due to expiration limitation 'no more than 6 month'")
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		sub, err := list.Subscriptions().GetByID(subID).Get(context.Background())
		if err != nil {
			t.Error(err)
			return
		}
		if sub.ID == "" {
			t.Error("can't parse subscription add response")
			return
		}
		subID = sub.ID
	})

	t.Run("GetByID#WrongID", func(t *testing.T) {
		if _, err := list.Subscriptions().GetByID("WrongID").Get(context.Background()); err == nil {
			t.Error("should fail with wrong id")
		}
	})

	t.Run("SetExpiration", func(t *testing.T) {
		sub, err := list.Subscriptions().GetByID(subID).Get(context.Background())
		if err != nil {
			t.Error(err)
			return
		}
		newExpiration := sub.ExpirationDateTime.Add(60 * time.Second)
		subUpd, err := list.Subscriptions().GetByID(subID).SetExpiration(context.Background(), newExpiration)
		if err != nil {
			t.Error(err)
			return
		}
		if newExpiration != subUpd.ExpirationDateTime {
			t.Error("can't set expiration date time")
		}
	})

	t.Run("SetClientState", func(t *testing.T) {
		subUpd, err := list.Subscriptions().GetByID(subID).SetClientState(context.Background(), "client state")
		if err != nil {
			t.Error(err)
			return
		}
		if subUpd.ClientState != "client state" {
			t.Error("can't set client state")
		}
	})

	t.Run("SetNotificationURL#WrongURL", func(t *testing.T) {
		if _, err := list.Subscriptions().GetByID(subID).SetNotificationURL(context.Background(), "wrong-url"); err == nil {
			t.Error("should fail with wrong URL")
		}
	})

	t.Run("GetSubscriptions", func(t *testing.T) {
		sub, err := list.Subscriptions().Get(context.Background())
		if err != nil {
			t.Error(err)
		}
		if len(sub) == 0 {
			t.Error("incorrect number of subscriptions")
		}
	})

	t.Run("Delete", func(t *testing.T) {
		resp, err := list.Subscriptions().Get(context.Background())
		if err != nil {
			t.Error(err)
		}
		for _, s := range resp {
			if err := list.Subscriptions().GetByID(s.ID).Delete(context.Background()); err != nil {
				t.Error(err)
			}
		}
		if resp, err := list.Subscriptions().Get(context.Background()); err != nil || len(resp) != 0 {
			t.Error("can't delete subscription(s)")
		}
	})

}
