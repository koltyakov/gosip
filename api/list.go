package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/koltyakov/gosip"
)

// List ...
type List struct {
	client   *gosip.SPClient
	conf     *Conf
	endpoint string
	oSelect  string
	oExpand  string
}

// Conf ...
func (list *List) Conf(conf *Conf) *List {
	list.conf = conf
	return list
}

// Select ...
func (list *List) Select(oDataSelect string) *List {
	list.oSelect = oDataSelect
	return list
}

// Expand ...
func (list *List) Expand(oDataExpand string) *List {
	list.oExpand = oDataExpand
	return list
}

// Get ...
func (list *List) Get() ([]byte, error) {
	apiURL, _ := url.Parse(list.endpoint)
	query := url.Values{}
	if list.oSelect != "" {
		query.Add("$select", trimMultiline(list.oSelect))
	}
	if list.oExpand != "" {
		query.Add("$expand", trimMultiline(list.oExpand))
	}
	apiURL.RawQuery = query.Encode()
	sp := NewHTTPClient(list.client)
	return sp.Get(apiURL.String(), getConfHeaders(list.conf))
}

// Delete ...
func (list *List) Delete() ([]byte, error) {
	sp := NewHTTPClient(list.client)
	return sp.Delete(list.endpoint, getConfHeaders(list.conf))
}

// Update ...
func (list *List) Update(body []byte) ([]byte, error) {
	sp := NewHTTPClient(list.client)
	return sp.Update(list.endpoint, body, getConfHeaders(list.conf))
}

// Items ...
func (list *List) Items() *Items {
	return &Items{
		client:   list.client,
		conf:     list.conf,
		endpoint: fmt.Sprintf("%s/items", list.endpoint),
	}
}

// GetEntityType ...
func (list *List) GetEntityType() (string, error) {
	headers := getConfHeaders(list.conf)
	headers["Accept"] = "application/json;odata=verbose"

	data, err := list.Select("ListItemEntityTypeFullName").Conf(&Conf{Headers: headers}).Get()
	if err != nil {
		return "", err
	}

	res := &struct {
		D struct {
			Results struct {
				ListItemEntityTypeFullName string `json:"ListItemEntityTypeFullName"`
			} `json:"results"`
		} `json:"d"`
	}{}

	if err := json.Unmarshal(data, &res); err != nil {
		return "", fmt.Errorf("unable to parse the response: %v", err)
	}

	return res.D.Results.ListItemEntityTypeFullName, nil
}

// Roles ...
func (list *List) Roles() *Roles {
	return &Roles{
		client:   list.client,
		conf:     list.conf,
		endpoint: list.endpoint,
	}
}
