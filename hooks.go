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
