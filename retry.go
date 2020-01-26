package gosip

import (
	"math"
	"net/http"
	"strconv"
	"time"
)

// RetryPolicies : error state requests default retry policies
var retryPolicies = map[int]int{
	401: 5, // on 401 :: Unauthorized
	429: 5, // on 429 :: Too many requests throttling error response
	500: 1, // on 500 :: Internal Server Error
	503: 5, // on 503 :: Service Unavailable
}

// getRetryPolicy receives retries policy retry number
func (c *SPClient) getRetryPolicy(statusCode int) int {
	// Apply default policies
	if c.RetryPolicies == nil {
		c.RetryPolicies = retryPolicies
	} else {
		// Append defaults to custom
		for status, retries := range retryPolicies {
			if _, ok := c.RetryPolicies[status]; !ok {
				c.RetryPolicies[status] = retries
			}
		}
	}
	return c.RetryPolicies[statusCode]
}

// shouldRetry checks should the request be retried, used with specific resp.StatusCode's
func (c *SPClient) shouldRetry(req *http.Request, resp *http.Response, retries int) bool {
	noRetry := req.Header.Get("X-Gosip-NoRetry")
	if noRetry == "true" {
		return false
	}
	retry, _ := strconv.Atoi(req.Header.Get("X-Gosip-Retry"))
	if retry < retries {
		retryAfter := 0
		if resp != nil {
			retryAfter, _ = strconv.Atoi(resp.Header.Get("Retry-After"))
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

func cloneHeader(h http.Header) http.Header {
	h2 := make(http.Header, len(h))
	for k, vv := range h {
		vv2 := make([]string, len(vv))
		copy(vv2, vv)
		h2[k] = vv2
	}
	return h2
}
