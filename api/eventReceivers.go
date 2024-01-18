package api

import (
	"context"
	"encoding/json"

	"github.com/koltyakov/gosip"
)

//go:generate ggen -ent EventReceivers -conf -coll -mods Select,Filter,Top,OrderBy

// EventReceivers represent SharePoint EventReceivers API queryable collection struct
// Always use NewEventReceivers constructor instead of &EventReceivers{}
type EventReceivers struct {
	client    *gosip.SPClient
	config    *RequestConfig
	endpoint  string
	modifiers *ODataMods
}

// EventReceiverInfo - event receivers API response payload structure
type EventReceiverInfo struct {
	EventType        int    `json:"EventType"`
	ReceiverAssembly string `json:"ReceiverAssembly"`
	ReceiverClass    string `json:"ReceiverClass"`
	ReceiverID       string `json:"ReceiverId"`
	ReceiverName     string `json:"ReceiverName"`
	ReceiverURL      string `json:"ReceiverUrl"`
	SequenceNumber   int    `json:"SequenceNumber"`
	Synchronization  int    `json:"Synchronization"`
}

// NewEventReceivers - EventReceivers struct constructor function
func NewEventReceivers(client *gosip.SPClient, endpoint string, config *RequestConfig) *EventReceivers {
	return &EventReceivers{
		client:    client,
		endpoint:  endpoint,
		config:    config,
		modifiers: NewODataMods(),
	}
}

// ToURL gets endpoint with modificators raw URL
func (eventReceivers *EventReceivers) ToURL() string {
	return toURL(eventReceivers.endpoint, eventReceivers.modifiers)
}

// Get gets event receivers collection
func (eventReceivers *EventReceivers) Get(ctx context.Context) ([]*EventReceiverInfo, error) {
	client := NewHTTPClient(eventReceivers.client)
	data, err := client.Get(ctx, eventReceivers.ToURL(), eventReceivers.config)
	if err != nil {
		return nil, err
	}
	data, _ = NormalizeODataCollection(data)
	var res []*EventReceiverInfo
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res, nil
}
