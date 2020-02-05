package gosip

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestEdges(t *testing.T) {
	siteURL := "http://localhost:8989"
	closer, err := startFakeServer(":8989", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// faking digest response
		if r.RequestURI == "/_api/ContextInfo" {
			fmt.Fprintf(w, `{"d":{"GetContextWebInformation":{"FormDigestValue":"","FormDigestTimeoutSeconds":120,"LibraryVersion":"FAKE"}}}`)
			return
		}
		if r.RequestURI == "/_api/wait" {
			time.Sleep(1 * time.Second)
			fmt.Fprintf(w, `{"result":"one eternity later"}`)
			return
		}
		fmt.Fprintf(w, `{ "result": "Cool alfter some retries" }`)
	}))
	if err != nil {
		t.Fatal(err)
	}
	defer closer.Close()

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
			t.Error(err)
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
		if strings.Index(err.Error(), "request canceled") == -1 {
			t.Error("request canceling failed")
		}
	})

}
