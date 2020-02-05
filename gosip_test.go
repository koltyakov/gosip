package gosip

import (
	"net/http"
	"testing"
)

func TestEdges(t *testing.T) {

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

}
