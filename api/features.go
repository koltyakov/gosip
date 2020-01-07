package api

import (
	"encoding/json"
	"fmt"

	"github.com/koltyakov/gosip"
)

// Features represent SharePoint Webs & Site Features API queryable collection struct
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

// Get gets features collection (IDs)
func (features *Features) Get() ([]*FeatureInfo, error) {
	sp := NewHTTPClient(features.client)
	data, err := sp.Get(features.endpoint, getConfHeaders(features.config))
	if err != nil {
		return nil, err
	}
	data, _ = NormalizeODataCollection(data)
	res := []*FeatureInfo{}
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// Add activates a feature by its ID (GUID) in the parent container (Site or Web)
func (features *Features) Add(featureID string, force bool) error {
	endpoint := fmt.Sprintf("%s/Add", features.endpoint)
	sp := NewHTTPClient(features.client)
	body := []byte(fmt.Sprintf(`{"featdefScope":0,"featureId":"%s","force":%t}`, featureID, force))
	_, err := sp.Post(endpoint, body, getConfHeaders(features.config))
	return err
}

// Remove deactivates a feature by its ID (GUID) in the parent container (Site or Web)
func (features *Features) Remove(featureID string, force bool) error {
	endpoint := fmt.Sprintf("%s/Remove", features.endpoint)
	sp := NewHTTPClient(features.client)
	body := []byte(fmt.Sprintf(`{"featureId":"%s","force":%t}`, featureID, force))
	_, err := sp.Post(endpoint, body, getConfHeaders(features.config))
	return err
}
