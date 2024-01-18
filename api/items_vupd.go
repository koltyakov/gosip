package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

// UpdateValidateResp - update validate response type with helper processor methods
type UpdateValidateResp []byte

// ValidateUpdateOptions ValidateUpdateListItem request options
type ValidateUpdateOptions struct {
	NewDocumentUpdate bool
	CheckInComment    string
}

// UpdateValidateFieldResult field result struct
type UpdateValidateFieldResult struct {
	ErrorCode    int
	ErrorMessage string
	FieldName    string
	FieldValue   string
	HasException bool
	ItemID       int `json:"ItemId"`
}

// UpdateValidate updates an item in this list using ValidateUpdateListItem method.
// formValues fingerprints https://github.com/koltyakov/sp-sig-20180705-demo/blob/master/src/03-pnp/FieldTypes.md#field-data-types-fingerprints-sample
func (item *Item) UpdateValidate(ctx context.Context, formValues map[string]string, options *ValidateUpdateOptions) (UpdateValidateResp, error) {
	endpoint := fmt.Sprintf("%s/ValidateUpdateListItem", item.endpoint)
	client := NewHTTPClient(item.client)
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
	}
	body, _ := json.Marshal(payload)

	var res UpdateValidateResp
	var err error

	res, err = client.Post(ctx, endpoint, bytes.NewBuffer(body), item.config)
	if err != nil {
		return res, err
	}

	var errs []error
	for _, f := range res.Data() {
		if f.HasException {
			errs = append(errs, fmt.Errorf("%s: %s", f.FieldName, f.ErrorMessage))
		}
	}
	if len(errs) > 0 {
		return res, fmt.Errorf("%v", errs)
	}

	return res, nil
}

/* UpdateValidate response helpers */

// Data unmarshals UpdateValidate response
func (uvResp *UpdateValidateResp) Data() []UpdateValidateFieldResult {
	var d []UpdateValidateFieldResult
	r := &struct {
		D struct {
			ValidateUpdateListItem struct {
				Results []UpdateValidateFieldResult `json:"results"`
			} `json:"ValidateUpdateListItem"`
		} `json:"d"`
		Value []UpdateValidateFieldResult `json:"value"`
	}{}
	_ = json.Unmarshal(*uvResp, &r)
	if r.Value != nil {
		return r.Value
	}
	if r.D.ValidateUpdateListItem.Results != nil {
		return r.D.ValidateUpdateListItem.Results
	}
	return d
}

// Value gets updated item's value from the response
func (uvResp *UpdateValidateResp) Value(fieldName string) string {
	dd := uvResp.Data()
	for _, d := range dd {
		if d.FieldName == fieldName {
			return d.FieldValue
		}
	}
	return ""
}
