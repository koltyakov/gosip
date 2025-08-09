package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/koltyakov/gosip"
)

// HTTPClient HTTP methods helper
type HTTPClient struct {
	sp *gosip.SPClient
}

// RequestConfig struct
type RequestConfig struct {
	Headers map[string]string
	Context context.Context
}

// HeadersPresets : SP REST OData headers presets
var HeadersPresets = struct {
	Verbose         *RequestConfig
	Minimalmetadata *RequestConfig
	Nometadata      *RequestConfig
}{
	Verbose: &RequestConfig{
		Headers: map[string]string{
			"Accept":          "application/json;odata=verbose",
			"Content-Type":    "application/json;odata=verbose;charset=utf-8",
			"Accept-Language": "en-US,en;q=0.9",
		},
	},
	Minimalmetadata: &RequestConfig{
		Headers: map[string]string{
			"Accept":          "application/json;odata=minimalmetadata",
			"Content-Type":    "application/json;odata=verbose;charset=utf-8",
			"Accept-Language": "en-US,en;q=0.9",
		},
	},
	Nometadata: &RequestConfig{
		Headers: map[string]string{
			"Accept":          "application/json;odata=nometadata",
			"Content-Type":    "application/json;odata=verbose;charset=utf-8",
			"Accept-Language": "en-US,en;q=0.9",
		},
	},
}

// applyDefaults applies context and header defaults for SP REST helpers
func applyDefaults(req *http.Request, conf *RequestConfig, needsBodyCT bool) {
	// Apply context
	if conf != nil && conf.Context != nil {
		req = req.WithContext(conf.Context)
	}
	// Default headers (Accept always; Content-Type only when body is expected)
	if req.Header.Get("Accept") == "" {
		req.Header.Set("Accept", gosip.DefaultAcceptVerbose) // default to SP2013 for backwards compatibility
	}
	if needsBodyCT && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", gosip.DefaultContentTypeVerbose)
	}
	// Apply custom headers last to allow overrides
	if conf != nil && conf.Headers != nil {
		for key, value := range conf.Headers {
			req.Header.Set(key, value)
		}
	}
}

// NewHTTPClient creates an instance of httpClient
func NewHTTPClient(spClient *gosip.SPClient) *HTTPClient {
	return &HTTPClient{sp: spClient}
}

// Get - generic GET request wrapper
func (client *HTTPClient) Get(endpoint string, conf *RequestConfig) ([]byte, error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create a request: %w", err)
	}
	applyDefaults(req, conf, false)

	resp, err := client.sp.Execute(req)
	if err != nil {
		return nil, fmt.Errorf("unable to request api: %w", err)
	}
	defer shut(resp.Body)

	return io.ReadAll(resp.Body)
}

// Post - generic POST request wrapper
func (client *HTTPClient) Post(endpoint string, body io.Reader, conf *RequestConfig) ([]byte, error) {
	// req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	req, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("unable to create a request: %w", err)
	}

	applyDefaults(req, conf, true)

	resp, err := client.sp.Execute(req)
	if err != nil {
		return nil, fmt.Errorf("unable to request api: %w", err)
	}
	defer shut(resp.Body)

	return io.ReadAll(resp.Body)
}

// Delete - generic DELETE request wrapper
func (client *HTTPClient) Delete(endpoint string, conf *RequestConfig) ([]byte, error) {
	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create a request: %w", err)
	}

	applyDefaults(req, conf, true)
	req.Header.Add("X-HTTP-Method", "DELETE")
	req.Header.Add("If-Match", "*")

	resp, err := client.sp.Execute(req)
	if err != nil {
		return nil, fmt.Errorf("unable to request api: %w", err)
	}
	defer shut(resp.Body)

	return io.ReadAll(resp.Body)
}

// Update - generic MERGE request wrapper
func (client *HTTPClient) Update(endpoint string, body io.Reader, conf *RequestConfig) ([]byte, error) {
	// req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	req, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("unable to create a request: %w", err)
	}

	applyDefaults(req, conf, true)
	req.Header.Add("X-HTTP-Method", "MERGE")
	req.Header.Add("If-Match", "*")

	resp, err := client.sp.Execute(req)
	if err != nil {
		return nil, fmt.Errorf("unable to request api: %w", err)
	}
	defer shut(resp.Body)

	return io.ReadAll(resp.Body)
}

// ProcessQuery - CSOM requests helper
func (client *HTTPClient) ProcessQuery(endpoint string, body io.Reader, conf *RequestConfig) ([]byte, error) {
	if !strings.Contains(strings.ToLower(endpoint), strings.ToLower("/_vti_bin/client.svc/ProcessQuery")) {
		endpoint = fmt.Sprintf("%s/_vti_bin/client.svc/ProcessQuery", getPriorEndpoint(endpoint, "/_api"))
	}

	// req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	req, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("unable to create a request: %w", err)
	}

	// Apply context
	if conf != nil && conf.Context != nil {
		req = req.WithContext(conf.Context)
	}

	// CSOM headers
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", `text/xml;charset="UTF-8"`)
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	// // Apply custom headers
	// if conf != nil && conf.Headers != nil {
	// 	for key, value := range conf.Headers {
	// 		req.Header.Set(key, value)
	// 	}
	// }

	resp, err := client.sp.Execute(req)
	if err != nil {
		return nil, fmt.Errorf("unable to request api: %w", err)
	}
	defer shut(resp.Body)

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// https://stackoverflow.com/questions/31398044/got-error-invalid-character-ï-looking-for-beginning-of-value-from-json-unmar
	data = bytes.TrimPrefix(data, []byte("\xef\xbb\xbf")) // removing BOM

	var arrRes []interface{}
	if err := json.Unmarshal(data, &arrRes); err != nil {
		return data, err
	}

	res := &struct {
		SchemaVersion  string `json:"SchemaVersion"`
		LibraryVersion string `json:"LibraryVersion"`
		ErrorInfo      *struct {
			ErrorMessage  string `json:"ErrorMessage"`
			ErrorValue    string `json:"ErrorValue"`
			ErrorCode     int    `json:"ErrorCode"`
			ErrorTypeName string `json:"ErrorTypeName"`
		} `json:"ErrorInfo"`
		TraceCorrelationID string `json:"TraceCorrelationId"`
	}{}

	arrEl1, err := json.Marshal(arrRes[0])
	if err != nil {
		return data, err
	}

	if err := json.Unmarshal(arrEl1, &res); err != nil {
		return data, err
	}

	if res.ErrorInfo != nil {
		return data, fmt.Errorf(
			"%s (Code: %d, %s, Correlation ID: %s)",
			res.ErrorInfo.ErrorMessage,
			res.ErrorInfo.ErrorCode,
			res.ErrorInfo.ErrorTypeName,
			res.TraceCorrelationID,
		)
	}

	return data, nil
}
