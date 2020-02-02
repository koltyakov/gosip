package api

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestWebs(t *testing.T) {
	checkClient(t)

	sp := NewSP(spClient)
	webs := sp.Web().Webs()
	endpoint := spClient.AuthCnfg.GetSiteURL() + "/_api/Web/Webs"
	newWebGUID := uuid.New().String()

	t.Run("Constructor", func(t *testing.T) {
		w := NewWebs(spClient, endpoint, nil)
		if _, err := w.Select("Id").Get(); err != nil {
			t.Error(err)
		}
	})

	t.Run("ToURL", func(t *testing.T) {
		if webs.ToURL() != endpoint {
			t.Errorf(
				"incorrect endpoint URL, expected \"%s\", received \"%s\"",
				endpoint,
				webs.ToURL(),
			)
		}
	})

	t.Run("AddWeb", func(t *testing.T) {
		if !heavyTests {
			t.Skip("setup SPAPI_HEAVY_TESTS env var to \"true\" to run this test")
		}
		if _, err := webs.Add("CI: "+newWebGUID, "ci_"+newWebGUID, nil); err != nil {
			t.Error(err)
		}
		data, err := webs.Select("Id,Title").Get()
		if err != nil {
			t.Error(err)
		}
		if len(data.Data()) == 0 {
			t.Error("wrong webs number")
		}
	})

	t.Run("GetWebs", func(t *testing.T) {
		data, err := webs.Select("Id,Title").Get()
		if err != nil {
			t.Error(err)
		}

		if bytes.Compare(data, data.Normalized()) == -1 {
			t.Error("wrong response normalization")
		}
	})

	t.Run("DeleteWeb", func(t *testing.T) {
		if !heavyTests {
			t.Skip("setup SPAPI_HEAVY_TESTS env var to \"true\" to run this test")
		}
		createdWebURL := spClient.AuthCnfg.GetSiteURL() + "/ci_" + newWebGUID
		subWeb := NewWeb(spClient, createdWebURL, nil)
		if err := subWeb.Delete(); err != nil {
			t.Error(err)
		}
	})

	t.Run("CreateFieldAsXML", func(t *testing.T) {
		title := strings.Replace(uuid.New().String(), "-", "", -1)
		schemaXML := `<Field Type="Text" DisplayName="` + title + `" MaxLength="255" Name="` + title + `" Title="` + title + `"></Field>`
		if _, err := sp.Web().Fields().CreateFieldAsXML(schemaXML, 0); err != nil {
			t.Error(err)
		}
		if err := sp.Web().Fields().GetByInternalNameOrTitle(title).Delete(); err != nil {
			t.Error(err)
		}
	})

}
