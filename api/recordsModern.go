package api

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// https://support.microsoft.com/en-us/office/apply-retention-labels-to-files-in-sharepoint-or-onedrive-11a6835b-ec9f-40db-8aca-6f5ef18132df

// LockRecordItem locks a record via modern API method SP.CompliancePolicy.SPPolicyStoreProxy.UnlockRecordItem()
func (records *Records) LockRecordItem() error {
	itemID, listURL, err := records.getItemContext()
	if err != nil {
		return err
	}

	client := NewHTTPClient(records.item.client)
	endpoint := fmt.Sprintf("%s/_api/SP.CompliancePolicy.SPPolicyStoreProxy.LockRecordItem()", getPriorEndpoint(records.item.endpoint, "/_api"))

	prop := map[string]interface{}{}
	prop["listUrl"] = listURL
	prop["itemId"] = itemID
	body, _ := json.Marshal(prop)

	_, err = client.Post(endpoint, bytes.NewBuffer(body), records.item.config)
	return err
}

// UnlockRecordItem unlocks a record via modern API method SP.CompliancePolicy.SPPolicyStoreProxy.UnlockRecordItem()
func (records *Records) UnlockRecordItem() error {
	itemID, listURL, err := records.getItemContext()
	if err != nil {
		return err
	}

	client := NewHTTPClient(records.item.client)
	endpoint := fmt.Sprintf("%s/_api/SP.CompliancePolicy.SPPolicyStoreProxy.UnlockRecordItem()", getPriorEndpoint(records.item.endpoint, "/_api"))

	prop := map[string]interface{}{}
	prop["listUrl"] = listURL
	prop["itemId"] = itemID
	body, _ := json.Marshal(prop)

	_, err = client.Post(endpoint, bytes.NewBuffer(body), records.item.config)
	return err
}

func (records *Records) getItemContext() (int, string, error) {
	// Get item's context
	data, err := records.item.Select("Id,ParentList/RootFolder/ServerRelativeURL").Expand("ParentList/RootFolder").Get()
	if err != nil {
		return 0, "", err
	}

	item := &struct {
		ID         int `json:"Id"`
		ParentList struct {
			RootFolder struct {
				ServerRelativeURL string `json:"ServerRelativeUrl"`
			} `json:"RootFolder"`
		} `json:"ParentList"`
	}{}

	if err := json.Unmarshal(data.Normalized(), &item); err != nil {
		return 0, "", err
	}

	return item.ID, item.ParentList.RootFolder.ServerRelativeURL, nil
}
