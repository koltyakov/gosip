package api

import (
	"encoding/json"
	"fmt"

	"github.com/koltyakov/gosip"
)

// Utility represents SharePoint Utilities namespace API object struct
// Always use NewUtility constructor instead of &Utility{}
type Utility struct {
	client   *gosip.SPClient
	config   *RequestConfig
	endpoint string
}

// EmailProps struct for SendEmail method parametes
type EmailProps struct {
	Subject string   // Email subject
	Body    string   // Email text or HTML body
	To      []string // Slice of To email addresses to whom Email is intended to be sent
	CC      []string // Slice of CC email addresses (optional)
	BCC     []string // Slice of BCC email addresses (optional)
	From    string   // Sender email addresses (optional)
}

// NewUtility - Utility struct constructor function
func NewUtility(client *gosip.SPClient, endpoint string, config *RequestConfig) *Utility {
	return &Utility{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// SendEmail sends an email via REST API due to the provided EmailProps options
func (utility *Utility) SendEmail(options *EmailProps) error {
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

	_, err := sp.Post(endpoint, body, getConfHeaders(utility.config))
	return err
}
