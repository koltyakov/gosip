package gosip

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	siteURL := "http://localhost:8989"
	closer, err := startFakeServer(":8989", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// faking digest response
		if r.RequestURI == "/_api/ContextInfo" {
			fmt.Fprintf(w, `{"d":{"GetContextWebInformation":{"FormDigestValue":"FAKE","FormDigestTimeoutSeconds":120,"LibraryVersion":"FAKE"}}}`)
			return
		}
		// retry after
		if r.RequestURI == "/_api/retryafter" && r.Header.Get("X-Gosip-Retry") == "1" {
			w.Header().Add("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{ "error": "Body is not backed off" }`))
			return
		}
		// ntlm retry
		if r.RequestURI == "/_api/ntlm" && r.Header.Get("X-Gosip-Retry") == "" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{ "error": "NTLM force retry" }`))
			return
		}
		// context cancel
		if r.RequestURI == "/_api/contextcancel" && r.Header.Get("X-Gosip-Retry") == "" {
			w.Header().Add("Retry-After", "5")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{ "error": "context cancel" }`))
			return
		}
		if r.Body != nil {
			defer r.Body.Close()
			data, _ := ioutil.ReadAll(r.Body)
			if r.RequestURI == "/_api/post/keepbody" && r.Header.Get("X-Gosip-Retry") == "1" {
				if fmt.Sprintf("%s", data) != "none-empty" {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{ "error": "Body is not backed off" }`))
					return
				}
			}
		}
		// backoff after 2 retries
		if r.Header.Get("X-Gosip-Retry") == "2" {
			fmt.Fprintf(w, `{ "result": "Cool alfter some retries" }`)
			return
		}
		// intentional 503
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(`{ "error": "503 Retry Please" }`))
	}))
	if err != nil {
		t.Fatal(err)
	}
	defer closer.Close()

	t.Run("GetRequest", func(t *testing.T) {
		client := &SPClient{
			AuthCnfg:      &AnonymousCnfg{SiteURL: siteURL},
			RetryPolicies: map[int]int{503: 3},
		}

		req, err := http.NewRequest("GET", client.AuthCnfg.GetSiteURL()+"/_api/get", nil)
		if err != nil {
			t.Fatal(err)
		}

		rsp, err := client.Execute(req)
		if err != nil {
			t.Error(err)
		}
		defer rsp.Body.Close()

		if rsp.StatusCode != 200 {
			t.Error("can't retry a request")
		}
	})

	t.Run("PostRequest", func(t *testing.T) {
		client := &SPClient{
			AuthCnfg:      &AnonymousCnfg{SiteURL: siteURL},
			RetryPolicies: map[int]int{503: 3},
		}

		req, err := http.NewRequest("POST", client.AuthCnfg.GetSiteURL()+"/_api/post/keepbody", bytes.NewBuffer([]byte("none-empty")))
		if err != nil {
			t.Fatal(err)
		}

		rsp, err := client.Execute(req)
		if err != nil {
			t.Error(err)
		}
		defer rsp.Body.Close()

		if rsp.StatusCode != 200 {
			t.Error("can't retry a request")
		}
	})

	t.Run("PostRequestEmptyBody", func(t *testing.T) {
		client := &SPClient{
			AuthCnfg:      &AnonymousCnfg{SiteURL: siteURL},
			RetryPolicies: map[int]int{503: 3},
		}

		req, err := http.NewRequest("POST", client.AuthCnfg.GetSiteURL()+"/_api/post", nil)
		if err != nil {
			t.Fatal(err)
		}

		rsp, err := client.Execute(req)
		if err != nil {
			t.Error(err)
		}
		defer rsp.Body.Close()

		if rsp.StatusCode != 200 {
			t.Error("can't retry a request")
		}
	})

	t.Run("PostRequestShould503", func(t *testing.T) {
		client := &SPClient{
			AuthCnfg:      &AnonymousCnfg{SiteURL: siteURL},
			RetryPolicies: map[int]int{503: 1},
		}

		req, err := http.NewRequest("POST", client.AuthCnfg.GetSiteURL()+"/_api/post", bytes.NewBuffer([]byte("none-empty")))
		if err != nil {
			t.Fatal(err)
		}

		rsp, _ := client.Execute(req)
		defer rsp.Body.Close()

		if rsp.StatusCode != 503 {
			t.Error("should receive 503")
		}
	})

	t.Run("DisableRetry", func(t *testing.T) {
		client := &SPClient{
			AuthCnfg: &AnonymousCnfg{SiteURL: siteURL},
		}

		req, err := http.NewRequest("GET", client.AuthCnfg.GetSiteURL()+"/_api/get", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Add("X-Gosip-NoRetry", "true")

		rsp, _ := client.Execute(req)
		defer rsp.Body.Close()

		if rsp.StatusCode != 503 {
			t.Error("should receive 503")
		}
	})

	t.Run("RetryAfter", func(t *testing.T) {
		client := &SPClient{
			AuthCnfg: &AnonymousCnfg{SiteURL: siteURL},
		}

		req, err := http.NewRequest("GET", client.AuthCnfg.GetSiteURL()+"/_api/retryafter", nil)
		if err != nil {
			t.Fatal(err)
		}

		beforeReq := time.Now()
		if _, err := client.Execute(req); err != nil {
			t.Error(err)
		}

		dur := time.Now().Sub(beforeReq)
		if dur < 1*time.Second {
			t.Error("retry after is ignored")
		}
	})

	t.Run("NtlmRetry", func(t *testing.T) {
		client := &SPClient{
			AuthCnfg: &AnonymousCnfg{
				SiteURL:  siteURL,
				Strategy: "ntlm",
			},
		}

		req, err := http.NewRequest("GET", client.AuthCnfg.GetSiteURL()+"/_api/ntlm", nil)
		if err != nil {
			t.Fatal(err)
		}

		if _, err := client.Execute(req); err != nil {
			t.Error(err)
		}
	})

	t.Run("ContextCancel", func(t *testing.T) {
		client := &SPClient{
			AuthCnfg: &AnonymousCnfg{SiteURL: siteURL},
		}

		req, err := http.NewRequest("GET", client.AuthCnfg.GetSiteURL()+"/_api/contextcancel", nil)
		if err != nil {
			t.Fatal(err)
		}
		ctx, cancel := context.WithCancel(context.Background())
		req = req.WithContext(ctx)

		beforeReq := time.Now()

		go func() {
			select {
			case <-time.After(900 * time.Millisecond):
				cancel()
			}
		}()

		client.Execute(req) // should be canceled with a context after 900 milliseconds

		dur := time.Now().Sub(beforeReq)
		if dur > 1*time.Second {
			t.Error("context canceling failed")
		}
	})

	t.Run("ContextCancel", func(t *testing.T) {
		client := &SPClient{
			AuthCnfg: &AnonymousCnfg{SiteURL: siteURL},
		}

		req, err := http.NewRequest("GET", client.AuthCnfg.GetSiteURL()+"/_api/contextcancel_2", nil)
		if err != nil {
			t.Fatal(err)
		}
		ctx, cancel := context.WithCancel(context.Background())
		req = req.WithContext(ctx)

		cancel()

		_, err = client.Execute(req) // should be prevented due to already closed context

		if strings.Index(err.Error(), "context canceled") == -1 {
			t.Error("context canceling failed")
		}
	})
}
