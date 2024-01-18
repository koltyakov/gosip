package api

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestHttpRetry(t *testing.T) {
	checkClient(t)

	sp := NewSP(spClient)

	t.Run("ShouldForceRetry", func(t *testing.T) {
		guid := uuid.New().String()
		if _, err := sp.Web().GetFolder("Shared Documents/"+guid).Folders().Add(context.Background(), "123"); err == nil {
			t.Error("should not succeeded, but force a retries")
		}
	})

}
