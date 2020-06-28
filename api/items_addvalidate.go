package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

// AddValidateResp - add validate using path response type with helper processor methods
type AddValidateResp []byte

// ValidateAddOptions AddValidateUpdateItemUsingPath method options
type ValidateAddOptions struct {
	DecodedPath       string
	NewDocumentUpdate bool
	CheckInComment    string
}

// AddValidate adds new item in this list using AddValidateUpdateItemUsingPath method.
// formValues fingerprints https://github.com/koltyakov/sp-sig-20180705-demo/blob/master/src/03-pnp/FieldTypes.md#field-data-types-fingerprints-sample
func (items *Items) AddValidate(formValues map[string]string, options *ValidateAddOptions) (AddValidateResp, error) {
	endpoint := fmt.Sprintf("%s/AddValidateUpdateItemUsingPath()", getPriorEndpoint(items.endpoint, "/items"))
	client := NewHTTPClient(items.client)
	type formValue struct {
		FieldName  string `json:"FieldName"`
		FieldValue string `json:"FieldValue"`
	}
	var formValuesArray []*formValue
	for n, v := range formValues {
		formValuesArray = append(formValuesArray, &formValue{
			FieldName:  n,
			FieldValue: v,
		})
	}
	payload := map[string]interface{}{"formValues": formValuesArray}
	if options != nil {
		payload["bNewDocumentUpdate"] = options.NewDocumentUpdate
		payload["checkInComment"] = options.CheckInComment
		if options.DecodedPath != "" {
			payload["listItemCreateInfo"] = map[string]interface{}{
				"__metadata": map[string]string{"type": "SP.ListItemCreationInformationUsingPath"},
				"FolderPath": map[string]interface{}{
					"__metadata": map[string]string{"type": "SP.ResourcePath"},
					"DecodedUrl": checkGetRelativeURL(options.DecodedPath, items.endpoint),
				},
			}
		}
	}
	body, _ := json.Marshal(payload)
	return client.Post(endpoint, bytes.NewBuffer(body), items.config)
}

/* AddValidate response helpers */

// Data unmarshals AddValidate response
func (avResp *AddValidateResp) Data() []map[string]interface{} {
	var d []map[string]interface{}
	r := &struct {
		D struct {
			AddValidateUpdateItemUsingPath struct {
				Results []map[string]interface{} `json:"results"`
			} `json:"AddValidateUpdateItemUsingPath"`
		} `json:"d"`
		Value []map[string]interface{} `json:"value"`
	}{}
	_ = json.Unmarshal(*avResp, &r)
	if r.Value != nil {
		return r.Value
	}
	if r.D.AddValidateUpdateItemUsingPath.Results != nil {
		return r.D.AddValidateUpdateItemUsingPath.Results
	}
	return d
}

// ID gets created item's ID from the response
func (avResp *AddValidateResp) ID() int {
	dd := avResp.Data()
	for _, d := range dd {
		if d["FieldName"] == "Id" {
			d, _ := strconv.Atoi(d["FieldValue"].(string))
			return d
		}
	}
	return 0
}
