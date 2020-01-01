package api

import (
	"encoding/json"
	"fmt"

	"github.com/koltyakov/gosip"
)

// Features represent SharePoint Webs & SIte Features API queryable collection struct
// Always use NewFeatures constructor instead of &Features{}
type Features struct {
	client   *gosip.SPClient
	config   *RequestConfig
	endpoint string
}

// FeatureInfo - features API response payload structure
type FeatureInfo struct {
	DefinitionID string `json:"DefinitionId"`
}

// NewFeatures - Features struct constructor function
func NewFeatures(client *gosip.SPClient, endpoint string, config *RequestConfig) *Features {
	return &Features{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// Get ...
func (features *Features) Get() ([]*FeatureInfo, error) {
	sp := NewHTTPClient(features.client)
	data, err := sp.Get(features.endpoint, getConfHeaders(features.config))
	if err != nil {
		return nil, err
	}
	data, _ = parseODataCollectionPlain(data)
	res := []*FeatureInfo{}
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// GetByID ...
// func (features *Features) GetByID(featureID string) (*FeatureInfo, error) {
// 	sp := NewHTTPClient(features.client)
// 	endpoint := fmt.Sprintf("%s('%s')", features.endpoint, featureID)
// 	data, err := sp.Get(endpoint, getConfHeaders(features.config))
// 	if err != nil {
// 		return nil, err
// 	}
// 	data = parseODataItem(data)
// 	res := &FeatureInfo{}
// 	if err := json.Unmarshal(data, &res); err != nil {
// 		return nil, err
// 	}
// 	return res, nil
// }

// Add ...
func (features *Features) Add(featureID string, force bool) ([]byte, error) {
	endpoint := fmt.Sprintf("%s/Add", features.endpoint)
	sp := NewHTTPClient(features.client)
	body := []byte(fmt.Sprintf(`{"featdefScope":0,"featureId":"%s","force":%t}`, featureID, force))
	return sp.Post(endpoint, body, getConfHeaders(features.config))
}

// Remove ...
func (features *Features) Remove(featureID string, force bool) ([]byte, error) {
	endpoint := fmt.Sprintf("%s/Remove", features.endpoint)
	sp := NewHTTPClient(features.client)
	body := []byte(fmt.Sprintf(`{"featureId":"%s","force":%t}`, featureID, force))
	return sp.Post(endpoint, body, getConfHeaders(features.config))
}
