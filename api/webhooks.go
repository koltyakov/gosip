package api

import "time"

// WebhookInfo represents webhook payload object structure
type WebhookInfo struct {
	SubscriptionID string    `json:"subscriptionId"`     // :"1111111111-3ef7-4917-ada1-xxxxxxxxxxxxx",
	ClientState    string    `json:"clientState"`        // :null,
	Expiration     time.Time `json:"expirationDateTime"` // :"2020-06-14T16:22:51.2160000Z",
	Resource       string    `json:"resource"`           // :"xxxxxx-c0ba-4063-a078-xxxxxxxxx",
	TenantID       string    `json:"tenantId"`           // :"4e2a1952-1ed1-4da3-85a6-xxxxxxxxxx",
	SiteURL        string    `json:"siteUrl"`            // :"/sites/webhooktest",
	WebID          string    `json:"webId"`              // :"xxxxx-3a7c-417b-964e-39f421c55d59"
}
