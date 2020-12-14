package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/koltyakov/gosip"
)

//go:generate ggen -ent Subscriptions -item Subscription -conf -coll
//go:generate ggen -ent Subscription -conf

// API reference: https://docs.microsoft.com/en-us/sharepoint/dev/apis/webhooks/lists/overview-sharepoint-list-webhooks

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
	ID string `json:"id"`
	// The service URL to send notifications to
	NotificationURL string `json:"notificationUrl"`
	// The date the notification will expire and be deleted
	ExpirationDateTime time.Time `json:"expirationDateTime"`
	// The URL of the list to receive notifications from
	Resource string `json:"resource"`
	// Optional. Opaque string passed back to the client on all notifications.
	// You can use this for validating notifications or tagging different subscriptions.
	ClientState  string `json:"clientState"`
	ResourceData string `json:"resourceData"`
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
func (subscriptions *Subscriptions) Get() ([]*SubscriptionInfo, error) {
	client := NewHTTPClient(subscriptions.client)
	resp, err := client.Get(subscriptions.endpoint, subscriptions.config)
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
func (subscriptions *Subscriptions) GetByID(subscriptionID string) *Subscription {
	return NewSubscription(
		subscriptions.client,
		fmt.Sprintf("%s('%s')", subscriptions.endpoint, subscriptionID),
		subscriptions.config,
	)
}

// Add adds/updates new subscription to a list
func (subscriptions *Subscriptions) Add(notificationURL string, expiration time.Time, clientState string) (*SubscriptionInfo, error) {
	client := NewHTTPClient(subscriptions.client)
	payload := map[string]interface{}{
		"notificationUrl":    notificationURL,
		"expirationDateTime": expiration,
		"resource":           getPriorEndpoint(subscriptions.endpoint, "/subscriptions"),
		"clientState":        clientState,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	conf := patchConfigHeaders(subscriptions.config, map[string]string{
		"Content-Type": "application/json",
	})
	resp, err := client.Post(subscriptions.endpoint, bytes.NewBuffer(body), conf)
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
func (subscription *Subscription) Get() (*SubscriptionInfo, error) {
	client := NewHTTPClient(subscription.client)
	resp, err := client.Get(subscription.endpoint, subscription.config)
	if err != nil {
		return nil, err
	}
	return subscription.parseResponse(resp)
}

// Delete deletes a subscription by its ID (GUID)
func (subscription *Subscription) Delete() error {
	client := NewHTTPClient(subscription.client)
	_, err := client.Delete(subscription.endpoint, subscription.config)
	return err
}

// Update updates a subscription
func (subscription *Subscription) Update(metadata map[string]interface{}) (*SubscriptionInfo, error) {
	client := NewHTTPClient(subscription.client)
	body, err := json.Marshal(metadata)
	if err != nil {
		return nil, err
	}
	conf := patchConfigHeaders(subscription.config, map[string]string{
		"Content-Type": "application/json",
	})
	if _, err := client.Update(subscription.endpoint, bytes.NewBuffer(body), conf); err != nil {
		return nil, err
	}
	return subscription.Get()
}

// SetExpiration sets new expiration datetime
func (subscription *Subscription) SetExpiration(expiration time.Time) (*SubscriptionInfo, error) {
	return subscription.Update(map[string]interface{}{
		"expirationDateTime": expiration,
	})
}

// SetNotificationURL sets new notification URL state
func (subscription *Subscription) SetNotificationURL(notificationURL string) (*SubscriptionInfo, error) {
	return subscription.Update(map[string]interface{}{
		"notificationUrl": notificationURL,
	})
}

// SetClientState sets new client state
func (subscription *Subscription) SetClientState(clientState string) (*SubscriptionInfo, error) {
	return subscription.Update(map[string]interface{}{
		"clientState": clientState,
	})
}

// parseResponse normalize and unmarshal subscription info response
func (subscription *Subscription) parseResponse(bytes []byte) (*SubscriptionInfo, error) {
	data := NormalizeODataItem(bytes)
	var sub *SubscriptionInfo
	if err := json.Unmarshal(data, &sub); err != nil {
		return nil, err
	}
	return sub, nil
}
