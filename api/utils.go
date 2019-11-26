package api

import (
	"net/url"
	"strings"
)

// GetConfHeaders resolves headers from config overrides
func getConfHeaders(conf *Conf) map[string]string {
	headers := map[string]string{}
	if conf != nil {
		headers = conf.Headers
	}
	return headers
}

// TrimMultiline - trims multiline
func trimMultiline(multi string) string {
	res := ""
	for _, line := range strings.Split(multi, "\n") {
		res += strings.Trim(line, "\t")
	}
	return res
}

// GetRelativeURL out of an absolute one
func getRelativeURL(absURL string) string {
	url, _ := url.Parse(absURL)
	return url.Path
}
