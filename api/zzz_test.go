package api

import (
	"testing"
)

func TestStats(t *testing.T) {
	checkClient(t)

	t.Run("PrintStats", func(t *testing.T) {
		t.Logf("Requests: %d\n", requestCntrs.Requests)
		t.Logf("Responses: %d\n", requestCntrs.Responses)
		t.Logf("Retries (including intended): %d\n", requestCntrs.Retries)
		t.Logf("Errors (including intended): %d\n", requestCntrs.Errors)
	})

}
