package api

import (
	"bytes"
	"testing"

	"github.com/koltyakov/gosip/csom"
)

func TestCsomRequest(t *testing.T) {
	checkClient(t)

	sp := NewHTTPClient(spClient)

	b := csom.NewBuilder()

	b.AddObject(csom.NewObject(`<Property Id="{{.ID}}" ParentId="{{.ParentID}}" Name="Web" />`), nil)
	b.AddAction(csom.NewAction(`
		<Query Id="{{.ID}}" ObjectPathId="{{.ObjectID}}">
			<Query SelectAllProperties="true">
				<Properties />
			</Query>
		</Query>
	`), nil)

	csomXML, err := b.Compile()
	if err != nil {
		t.Error(err)
	}

	if _, err := sp.ProcessQuery(spClient.AuthCnfg.GetSiteURL(), bytes.NewBuffer([]byte(csomXML))); err != nil {
		t.Error(err)
	}
}
