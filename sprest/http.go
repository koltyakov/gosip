package sprest

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/koltyakov/gosip"
)

// HTTPClient HTTP methods helper
type HTTPClient struct {
	SPClient *gosip.SPClient
}

// Get - generic GET request wrapper
func (ctx *HTTPClient) Get(endpoint string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create a request: %v", err)
	}

	// Default headers
	req.Header.Set("Accept", "application/json;odata=verbose") // default to SP2013 for backwards compatibility

	// Apply custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := ctx.SPClient.Execute(req)
	if err != nil {
		return nil, fmt.Errorf("unable to request api: %v", err)
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// Post - generic POST request wrapper
func (ctx *HTTPClient) Post(endpoint string, body []byte, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("unable to create a request: %v", err)
	}

	// Default headers
	req.Header.Set("Accept", "application/json;odata=verbose") // default to SP2013 for backwards compatibility
	req.Header.Set("Content-Type", "application/json;odata=verbose;charset=utf-8")

	// Apply custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := ctx.SPClient.Execute(req)
	if err != nil {
		return nil, fmt.Errorf("unable to request api: %v", err)
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// Delete - generic DELETE request wrapper
func (ctx *HTTPClient) Delete(endpoint string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create a request: %v", err)
	}

	// Default headers
	req.Header.Set("Accept", "application/json;odata=verbose") // default to SP2013 for backwards compatibility
	req.Header.Set("Content-Type", "application/json;odata=verbose;charset=utf-8")
	req.Header.Add("X-Http-Method", "DELETE")
	req.Header.Add("If-Match", "*")

	// Apply custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := ctx.SPClient.Execute(req)
	if err != nil {
		return nil, fmt.Errorf("unable to request api: %v", err)
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// Update - generic MERGE request wrapper
func (ctx *HTTPClient) Update(endpoint string, body []byte, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("unable to create a request: %v", err)
	}

	// Default headers
	req.Header.Set("Accept", "application/json;odata=verbose") // default to SP2013 for backwards compatibility
	req.Header.Set("Content-Type", "application/json;odata=verbose;charset=utf-8")
	req.Header.Add("X-Http-Method", "MERGE")
	req.Header.Add("If-Match", "*")

	// Apply custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := ctx.SPClient.Execute(req)
	if err != nil {
		return nil, fmt.Errorf("unable to request api: %v", err)
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
