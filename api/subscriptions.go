package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/koltyakov/gosip"
)

//go:generate ggen -ent Subscriptions -item Subscription -conf -coll
//go:generate ggen -ent Subscription -conf

// Subscriptions represent SharePoint lists Subscriptions API queryable collection struct
// Always use NewSubscriptions constructor instead of &Subscriptions{}
type Subscriptions struct {
	client   *gosip.SPClient
	config   *RequestConfig
	endpoint string
}

// Subscription represent SharePoint lists Subscription API
// Always use NewSubscription constructor instead of &Subscription{}
type Subscription struct {
	client   *gosip.SPClient
	config   *RequestConfig
	endpoint string
}

// SubscriptionInfo list subscription info
type SubscriptionInfo struct {
	ID                 string    `json:"id"`
	NotificationURL    string    `json:"notificationUrl"`
	ExpirationDateTime time.Time `json:"expirationDateTime"`
	Resource           string    `json:"resource"`
	ClientState        string    `json:"clientState"`
	ResourceData       string    `json:"resourceData"`
}

// NewSubscriptions - Subscriptions struct constructor function
func NewSubscriptions(client *gosip.SPClient, endpoint string, config *RequestConfig) *Subscriptions {
	return &Subscriptions{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// NewSubscription - subscription struct constructor function
func NewSubscription(client *gosip.SPClient, endpoint string, config *RequestConfig) *Subscription {
	return &Subscription{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// Get gets list subscriptions response collection
func (ss *Subscriptions) Get() ([]*SubscriptionInfo, error) {
	client := NewHTTPClient(ss.client)
	resp, err := client.Get(ss.endpoint, ss.config)
	if err != nil {
		return nil, err
	}
	data, _ := NormalizeODataCollection(resp)
	var subs []*SubscriptionInfo
	if err := json.Unmarshal(data, &subs); err != nil {
		return nil, err
	}
	return subs, nil
}

// GetByID gets list subscription by its ID (GUID)
func (ss *Subscriptions) GetByID(subscriptionID string) *Subscription {
	return NewSubscription(
		ss.client,
		fmt.Sprintf("%s('%s')", ss.endpoint, subscriptionID),
		ss.config,
	)
}

// Add adds/updates new subscription to a list
func (ss *Subscriptions) Add(notificationURL string, expiration time.Time, clientState string) (*SubscriptionInfo, error) {
	client := NewHTTPClient(ss.client)
	payload := map[string]interface{}{
		"notificationUrl":    notificationURL,
		"expirationDateTime": expiration,
		"resource":           strings.Replace(ss.endpoint, getPriorEndpoint(ss.endpoint, "_api"), "", -1),
		"clientState":        clientState,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	conf := patchConfigHeaders(ss.config, map[string]string{
		"Content-Type": "application/json",
	})
	resp, err := client.Post(ss.endpoint, bytes.NewBuffer(body), conf)
	if err != nil {
		return nil, err
	}
	data := NormalizeODataItem(resp)
	var sub *SubscriptionInfo
	if err := json.Unmarshal(data, &sub); err != nil {
		return nil, err
	}
	return sub, nil
}

// Get gets subscription info
func (s *Subscription) Get() (*SubscriptionInfo, error) {
	client := NewHTTPClient(s.client)
	resp, err := client.Get(s.endpoint, s.config)
	if err != nil {
		return nil, err
	}
	return s.parseResponse(resp)
}

// Delete deletes a subscription by its ID (GUID)
func (s *Subscription) Delete() error {
	client := NewHTTPClient(s.client)
	_, err := client.Delete(s.endpoint, s.config)
	return err
}

// Update updates a subscription
func (s *Subscription) Update(metadata map[string]interface{}) (*SubscriptionInfo, error) {
	client := NewHTTPClient(s.client)
	body, err := json.Marshal(metadata)
	if err != nil {
		return nil, err
	}
	conf := patchConfigHeaders(s.config, map[string]string{
		"Content-Type": "application/json",
	})
	if _, err := client.Update(s.endpoint, bytes.NewBuffer(body), conf); err != nil {
		return nil, err
	}
	return s.Get()
}

// SetExpiration sets new expiration datetime
func (s *Subscription) SetExpiration(expiration time.Time) (*SubscriptionInfo, error) {
	return s.Update(map[string]interface{}{
		"expirationDateTime": expiration,
	})
}

// SetNotificationURL sets new notification URL state
func (s *Subscription) SetNotificationURL(notificationURL string) (*SubscriptionInfo, error) {
	return s.Update(map[string]interface{}{
		"notificationUrl": notificationURL,
	})
}

// SetClientState sets new client state
func (s *Subscription) SetClientState(clientState string) (*SubscriptionInfo, error) {
	return s.Update(map[string]interface{}{
		"clientState": clientState,
	})
}

// parseResponse normalize and unmarshal subscription info response
func (s *Subscription) parseResponse(bytes []byte) (*SubscriptionInfo, error) {
	data := NormalizeODataItem(bytes)
	var sub *SubscriptionInfo
	if err := json.Unmarshal(data, &sub); err != nil {
		return nil, err
	}
	return sub, nil
}
