package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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

// NewHTTPClient creates an instance of httpClient
func NewHTTPClient(spClient *gosip.SPClient) *HTTPClient {
	return &HTTPClient{sp: spClient}
}

// Get - generic GET request wrapper
func (client *HTTPClient) Get(endpoint string, conf *RequestConfig) ([]byte, error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create a request: %v", err)
	}

	// Apply context
	if conf != nil && conf.Context != nil {
		req = req.WithContext(conf.Context)
	}

	// Default headers
	req.Header.Set("Accept", "application/json;odata=verbose") // default to SP2013 for backwards compatibility

	// Apply custom headers
	if conf != nil && conf.Headers != nil {
		for key, value := range conf.Headers {
			req.Header.Set(key, value)
		}
	}

	resp, err := client.sp.Execute(req)
	if err != nil {
		return nil, fmt.Errorf("unable to request api: %v", err)
	}
	defer shut(resp.Body)

	return ioutil.ReadAll(resp.Body)
}

// Post - generic POST request wrapper
func (client *HTTPClient) Post(endpoint string, body io.Reader, conf *RequestConfig) ([]byte, error) {
	// req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	req, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("unable to create a request: %v", err)
	}

	// Apply context
	if conf != nil && conf.Context != nil {
		req = req.WithContext(conf.Context)
	}

	// Default headers
	req.Header.Set("Accept", "application/json;odata=verbose") // default to SP2013 for backwards compatibility
	req.Header.Set("Content-Type", "application/json;odata=verbose;charset=utf-8")

	// Apply custom headers
	if conf != nil && conf.Headers != nil {
		for key, value := range conf.Headers {
			req.Header.Set(key, value)
		}
	}

	resp, err := client.sp.Execute(req)
	if err != nil {
		return nil, fmt.Errorf("unable to request api: %v", err)
	}
	defer shut(resp.Body)

	return ioutil.ReadAll(resp.Body)
}

// Delete - generic DELETE request wrapper
func (client *HTTPClient) Delete(endpoint string, conf *RequestConfig) ([]byte, error) {
	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create a request: %v", err)
	}

	// Apply context
	if conf != nil && conf.Context != nil {
		req = req.WithContext(conf.Context)
	}

	// Default headers
	req.Header.Set("Accept", "application/json;odata=verbose") // default to SP2013 for backwards compatibility
	req.Header.Set("Content-Type", "application/json;odata=verbose;charset=utf-8")
	req.Header.Add("X-Http-Method", "DELETE")
	req.Header.Add("If-Match", "*")

	// Apply custom headers
	if conf != nil && conf.Headers != nil {
		for key, value := range conf.Headers {
			req.Header.Set(key, value)
		}
	}

	resp, err := client.sp.Execute(req)
	if err != nil {
		return nil, fmt.Errorf("unable to request api: %v", err)
	}
	defer shut(resp.Body)

	return ioutil.ReadAll(resp.Body)
}

// Update - generic MERGE request wrapper
func (client *HTTPClient) Update(endpoint string, body io.Reader, conf *RequestConfig) ([]byte, error) {
	// req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	req, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("unable to create a request: %v", err)
	}

	// Apply context
	if conf != nil && conf.Context != nil {
		req = req.WithContext(conf.Context)
	}

	// Default headers
	req.Header.Set("Accept", "application/json;odata=verbose") // default to SP2013 for backwards compatibility
	req.Header.Set("Content-Type", "application/json;odata=verbose;charset=utf-8")
	req.Header.Add("X-Http-Method", "MERGE")
	req.Header.Add("If-Match", "*")

	// Apply custom headers
	if conf != nil {
		for key, value := range conf.Headers {
			req.Header.Set(key, value)
		}
	}

	resp, err := client.sp.Execute(req)
	var headers map[string]string
	if conf != nil {
		headers = conf.Headers
	}
	if err != nil && headers != nil {
		return nil, fmt.Errorf("unable to request api: %v", err)
	}
	defer shut(resp.Body)

	return ioutil.ReadAll(resp.Body)
}

// ProcessQuery - CSOM requests helper
func (client *HTTPClient) ProcessQuery(endpoint string, body io.Reader, conf *RequestConfig) ([]byte, error) {
	if strings.Index(strings.ToLower(endpoint), strings.ToLower("/_vti_bin/client.svc/ProcessQuery")) == -1 {
		endpoint = fmt.Sprintf("%s/_vti_bin/client.svc/ProcessQuery", getPriorEndpoint(endpoint, "/_api"))
	}

	// req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	req, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("unable to create a request: %v", err)
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
		return nil, fmt.Errorf("unable to request api: %v", err)
	}
	defer shut(resp.Body)

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

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
