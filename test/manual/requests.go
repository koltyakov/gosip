package manual

import (
	"fmt"

	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/api"
)

// CheckBasicPost : try creating an item
func CheckBasicPost(client *gosip.SPClient) (string, error) {
	sp := api.NewHTTPClient(client)
	endpoint := client.AuthCnfg.GetSiteURL() + "/_api/web/lists/getByTitle('Custom')/items"
	body := `{"__metadata":{"type":"SP.Data.CustomListItem"},"Title":"Test"}`

	data, err := sp.Post(endpoint, []byte(body), nil)
	if err != nil {
		return "", fmt.Errorf("unable to read a response: %v", err)
	}

	return string(data), nil
}
