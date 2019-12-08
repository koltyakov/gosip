package api

import (
	"encoding/json"
	"fmt"

	"github.com/koltyakov/gosip"
)

// Utility ...
type Utility struct {
	client   *gosip.SPClient
	config   *RequestConfig
	endpoint string
}

// EmailProps ...
type EmailProps struct {
	Subject string
	Body    string
	To      []string
	CC      []string
	BCC     []string
	From    string
}

// NewUtility ...
func NewUtility(client *gosip.SPClient, endpoint string, config *RequestConfig) *Utility {
	return &Utility{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// SendEmail ...
func (utility *Utility) SendEmail(options *EmailProps) ([]byte, error) {
	endpoint := fmt.Sprintf(
		"%s/_api/SP.Utilities.Utility.SendEmail",
		getPriorEndpoint(utility.endpoint, "/_api"),
	)
	sp := NewHTTPClient(utility.client)

	properties := map[string]interface{}{}
	properties["__metadata"] = map[string]string{"type": "SP.Utilities.EmailProperties"}
	properties["Subject"] = options.Subject
	properties["Body"] = options.Body
	if options.From != "" {
		properties["From"] = options.From
	}
	if len(options.To) > 0 {
		properties["To"] = map[string][]string{"results": options.To}
	}
	if len(options.CC) > 0 {
		properties["CC"] = map[string][]string{"results": options.CC}
	}
	if len(options.BCC) > 0 {
		properties["BCC"] = map[string][]string{"results": options.BCC}
	}
	props, _ := json.Marshal(properties)
	JSONProps := fmt.Sprintf("%s", props)
	body := []byte(trimMultiline(`{ "properties": ` + JSONProps + `}`))

	return sp.Post(endpoint, body, getConfHeaders(utility.config))
}
