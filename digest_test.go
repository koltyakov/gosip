package gosip

import (
	"fmt"
	"net/http"
	"testing"
)

func TestDigest(t *testing.T) {
	siteURL := "http://localhost:8989"
	digestTriggered := false
	closer, err := startFakeServer(":8989", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// faking digest response
		if r.RequestURI == "/_api/ContextInfo" {
			digestTriggered = true
			_, _ = fmt.Fprintf(w, `{"d":{"GetContextWebInformation":{"FormDigestValue":"FAKE","FormDigestTimeoutSeconds":120,"LibraryVersion":"FAKE"}}}`)
			return
		}
		_, _ = fmt.Fprintf(w, `{ "result": "Cool alfter some retries" }`)
	}))
	if err != nil {
		t.Fatal(err)
	}
	defer shut(closer)

	t.Run("ShouldTriggerDigest", func(t *testing.T) {
		client := &SPClient{
			AuthCnfg: &AnonymousCnfg{SiteURL: siteURL},
		}

		req, err := http.NewRequest("POST", client.AuthCnfg.GetSiteURL()+"/_api/post", nil)
		if err != nil {
			t.Fatal(err)
		}

		if _, err := client.Execute(req); err != nil {
			t.Error(err)
		}

		if !digestTriggered {
			t.Error("digest is not retrieved")
		}
	})

}
