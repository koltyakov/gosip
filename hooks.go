package gosip

import (
	"net/http"
	"time"
)

// HookHandlers struct to configure events handlers
type HookHandlers struct {
	OnError    func(event *HookEvent) // when error appeared
	OnRetry    func(event *HookEvent) // before retry request
	OnRequest  func(event *HookEvent) // before request is sent
	OnResponse func(event *HookEvent) // after response is received
}

// HookEvent hook event parameters struct
type HookEvent struct {
	Request    *http.Request
	StartedAt  time.Time
	StatusCode int
	Error      error
}

// onError on error hook handler
func (c *SPClient) onError(req *http.Request, startAt time.Time, statusCode int, err error) {
	if c.Hooks != nil && c.Hooks.OnError != nil {
		c.Hooks.OnError(&HookEvent{
			Request:    req,
			StartedAt:  startAt,
			StatusCode: statusCode,
			Error:      err,
		})
	}
}

// onRetry on retry hook handler
func (c *SPClient) onRetry(req *http.Request, startAt time.Time, statusCode int, err error) {
	if c.Hooks != nil && c.Hooks.OnRetry != nil {
		c.Hooks.OnRetry(&HookEvent{
			Request:    req,
			StartedAt:  startAt,
			StatusCode: statusCode,
			Error:      err,
		})
	}
}

// onResponse on response hook handler
func (c *SPClient) onResponse(req *http.Request, startAt time.Time, statusCode int, err error) {
	if c.Hooks != nil && c.Hooks.OnResponse != nil {
		c.Hooks.OnResponse(&HookEvent{
			Request:    req,
			StartedAt:  startAt,
			StatusCode: statusCode,
			Error:      err,
		})
	}
}

// onRequest on response hook handler
func (c *SPClient) onRequest(req *http.Request, startAt time.Time, statusCode int, err error) {
	if c.Hooks != nil && c.Hooks.OnRequest != nil {
		c.Hooks.OnRequest(&HookEvent{
			Request:    req,
			StartedAt:  startAt,
			StatusCode: statusCode,
			Error:      err,
		})
	}
}
