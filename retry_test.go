package gosip

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestRetry(t *testing.T) {
	siteURL := "http://localhost:8989"
	closer, err := startFakeServer(":8989", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// faking digest response
		if r.RequestURI == "/_api/ContextInfo" {
			fmt.Fprintf(w, `{"d":{"GetContextWebInformation":{"FormDigestValue":"FAKE","FormDigestTimeoutSeconds":120,"LibraryVersion":"FAKE"}}}`)
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

}
