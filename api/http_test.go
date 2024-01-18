package api

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestHttp(t *testing.T) {
	checkClient(t)

	t.Run("CSOMErrorHandling", func(t *testing.T) {
		client := NewHTTPClient(spClient)
		body := []byte(TrimMultiline(`
			<Request xmlns="http://schemas.microsoft.com/sharepoint/clientquery/2009" SchemaVersion="15.0.0.0" LibraryVersion="15.0.0.0" ApplicationName="Gosip">
				<Actions>
					<Query Id="4" ObjectPathId="3">
						<Query SelectAllProperties="true">
							<Properties />
						</Query>
					</Query>
				</Actions>
				<ObjectPaths>
					<StaticProperty Id="0" TypeId="{3747adcd-a3c3-41b9-bfab-4a64dd2f1e0a}" Name="Current" />
					<Property Id="2" ParentId="0" Name="Web" />
				</ObjectPaths>
			</Request>
		`))
		if _, err := client.ProcessQuery(context.Background(), spClient.AuthCnfg.GetSiteURL(), bytes.NewBuffer(body), nil); err != nil {
			if !strings.Contains(err.Error(), "Microsoft.SharePoint.Client.InvalidClientQueryException") {
				t.Error(err)
			}
		}
	})

}
