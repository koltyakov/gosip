package api

import (
	"encoding/json"

	"github.com/koltyakov/gosip"
)

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
func (receivers *EventReceivers) ToURL() string {
	return toURL(receivers.endpoint, receivers.modifiers)
}

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (receivers *EventReceivers) Conf(config *RequestConfig) *EventReceivers {
	receivers.config = config
	return receivers
}

// Select adds $select OData modifier
func (receivers *EventReceivers) Select(oDataSelect string) *EventReceivers {
	receivers.modifiers.AddSelect(oDataSelect)
	return receivers
}

// Filter adds $filter OData modifier
func (receivers *EventReceivers) Filter(oDataFilter string) *EventReceivers {
	receivers.modifiers.AddFilter(oDataFilter)
	return receivers
}

// Top adds $top OData modifier
func (receivers *EventReceivers) Top(oDataTop int) *EventReceivers {
	receivers.modifiers.AddTop(oDataTop)
	return receivers
}

// OrderBy adds $orderby OData modifier
func (receivers *EventReceivers) OrderBy(oDataOrderBy string, ascending bool) *EventReceivers {
	receivers.modifiers.AddOrderBy(oDataOrderBy, ascending)
	return receivers
}

// Get gets event receivers collection
func (receivers *EventReceivers) Get() ([]*EventReceiverInfo, error) {
	sp := NewHTTPClient(receivers.client)
	data, err := sp.Get(receivers.ToURL(), getConfHeaders(receivers.config))
	if err != nil {
		return nil, err
	}
	data, _ = NormalizeODataCollection(data)
	res := []*EventReceiverInfo{}
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res, nil
}
