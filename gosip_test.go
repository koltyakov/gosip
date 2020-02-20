package gosip

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestEdges(t *testing.T) {
	siteURL := "http://localhost:8989/sub" // sub URI to avoid digest caching when running tests in parallel
	closer, err := startFakeServer(":8989", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// faking digest response
		if r.RequestURI == "/sub/_api/ContextInfo" {
			_, _ = fmt.Fprintf(w, `{"d":{"GetContextWebInformation":{"FormDigestValue":"","FormDigestTimeoutSeconds":120,"LibraryVersion":"FAKE"}}}`)
			return
		}
		if r.RequestURI == "/sub/_api/wait" {
			time.Sleep(1 * time.Second)
			_, _ = fmt.Fprintf(w, `{"result":"one eternity later"}`)
			return
		}
		_, _ = fmt.Fprintf(w, `{ "result": "Cool alfter some retries" }`)
	}))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = closer.Close() }()

	t.Run("EmptyURLShouldFail", func(t *testing.T) {
		client := &SPClient{
			AuthCnfg: &AnonymousCnfg{SiteURL: ""},
		}

		req, err := http.NewRequest("POST", client.AuthCnfg.GetSiteURL()+"/_api/post", nil)
		if err != nil {
			t.Fatal(err)
		}

		if _, err := client.Execute(req); err == nil {
			t.Error(err)
		}
	})

	t.Run("ImcorrectConfigShouldFail", func(t *testing.T) {
		client := &SPClient{
			AuthCnfg:   &AnonymousCnfg{},
			ConfigPath: "incorrect",
		}

		req, err := http.NewRequest("POST", client.AuthCnfg.GetSiteURL()+"/_api/post", nil)
		if err != nil {
			t.Fatal(err)
		}

		if _, err := client.Execute(req); err == nil {
			t.Error(err)
		}
	})

	t.Run("SetAuthReturnError", func(t *testing.T) {
		client := &SPClient{
			AuthCnfg: &AnonymousCnfg{
				SiteURL: "http://restricted",
			},
		}

		req, err := http.NewRequest("POST", client.AuthCnfg.GetSiteURL()+"/_api/post", nil)
		if err != nil {
			t.Fatal(err)
		}

		if _, err := client.Execute(req); err == nil {
			t.Error(err)
		}
	})

	t.Run("OnDigestFailed", func(t *testing.T) {
		client := &SPClient{
			AuthCnfg: &AnonymousCnfg{SiteURL: siteURL},
		}

		req, err := http.NewRequest("POST", client.AuthCnfg.GetSiteURL()+"/_api/faildigest", nil)
		if err != nil {
			t.Fatal(err)
		}

		if _, err := client.Execute(req); err == nil {
			t.Error("should fail to retrieve digest")
		}
	})

	t.Run("ClientTimeout", func(t *testing.T) {
		client := &SPClient{
			AuthCnfg: &AnonymousCnfg{SiteURL: siteURL},
		}
		client.Timeout = 1 * time.Millisecond

		req, err := http.NewRequest("GET", client.AuthCnfg.GetSiteURL()+"/_api/wait", nil)
		if err != nil {
			t.Fatal(err)
		}

		_, err = client.Execute(req) // should fail after a timeout
		if err != nil && strings.Index(err.Error(), "request canceled") == -1 {
			t.Error("request canceling failed")
		}
	})

}
