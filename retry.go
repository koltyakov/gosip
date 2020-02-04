package gosip

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"
)

// RetryPolicies : error state requests default retry policies
var retryPolicies = map[int]int{
	401: 5,  // on 401 - Unauthorized
	429: 5,  // on 429 - Too many requests throttling error response
	500: 1,  // on 500 - Internal Server Error
	503: 10, // on 503 - Service Unavailable
	504: 5,  // on 504 - Gateway Timeout Error
}

// getRetryPolicy receives retries policy retry number
func (c *SPClient) getRetryPolicy(statusCode int) int {
	// Return defaults when no custom
	if c.RetryPolicies == nil {
		return retryPolicies[statusCode]
	}
	// Check in custom
	retries, ok := c.RetryPolicies[statusCode]
	if !ok {
		// Fallback to default
		return retryPolicies[statusCode]
	}
	return retries
}

// shouldRetry checks should the request be retried, used with specific resp.StatusCode's
func (c *SPClient) shouldRetry(req *http.Request, resp *http.Response, retries int) bool {
	noRetry := req.Header.Get("X-Gosip-NoRetry")
	if noRetry == "true" {
		return false
	}
	retry, _ := strconv.Atoi(req.Header.Get("X-Gosip-Retry"))
	if retry < retries {
		resp.Body.Close() // closing to reuse request
		retryAfter := 0
		if resp != nil && resp.StatusCode == 429 { // sometimes SPO is abusing Retry-After header on 503 errors
			retryAfter, _ = strconv.Atoi(resp.Header.Get("Retry-After"))
			fmt.Println("Retry after", retryAfter)
		}
		req.Header.Set("X-Gosip-Retry", strconv.Itoa(retry+1))
		if retryAfter != 0 {
			time.Sleep(time.Duration(retryAfter) * time.Second) // wait for Retry-After header info value
		} else {
			time.Sleep(time.Duration(100*math.Pow(2, float64(retry))) * time.Millisecond) // no Retry-After header
		}
		return true
	}
	return false
}
