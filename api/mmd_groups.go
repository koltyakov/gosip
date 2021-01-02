package api

import (
	"fmt"
	"strings"

	"github.com/koltyakov/gosip/csom"
)

// TermGroups term groups struct
type TermGroups struct {
	client   *HTTPClient
	config   *RequestConfig
	endpoint string

	csomEntry csom.Builder
}

// Get gets term groups metadata
func (termGroups *TermGroups) Get() ([]map[string]interface{}, error) {
	b := termGroups.csomEntry.Clone()
	b.AddAction(csom.NewQueryWithProps([]string{
		`<Property Name="Groups" SelectAll="true" />`,
	}), nil)

	csomResp, err := getCSOMResponse(termGroups.client, termGroups.endpoint, termGroups.config, b)
	if err != nil {
		return nil, err
	}

	groups, ok := csomResp["Groups"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("can't get groups from term store")
	}

	items, ok := groups["_Child_Items_"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("can't get child items from groups")
	}

	var resItems []map[string]interface{}
	for _, item := range items {
		resItem, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("can't get item from groups")
		}
		resItems = append(resItems, resItem)
	}

	return resItems, nil
}

// GetByID gets term group object by its GUID
func (termGroups *TermGroups) GetByID(groupGUID string) *TermGroup {
	return &TermGroup{
		client:   termGroups.client,
		endpoint: termGroups.endpoint,
		config:   termGroups.config,
		id:       trimTaxonomyGUID(groupGUID),

		csomEntry:   termGroups.csomEntry.Clone(),
		selectProps: []string{},
	}
}

/* Term Group */

// TermGroup term group struct
type TermGroup struct {
	client   *HTTPClient
	config   *RequestConfig
	endpoint string

	id string

	csomEntry   csom.Builder
	termStore   *TermStore
	selectProps []string
}

// csomBuilderEntry gets CSOM builder entry
func (termGroup *TermGroup) csomBuilderEntry() csom.Builder {
	b := termGroup.csomEntry.Clone()
	b.AddObject(csom.NewObjectMethod("GetGroup", []string{
		fmt.Sprintf(`<Parameter Type="String">%s</Parameter>`, termGroup.id),
	}), nil)
	return b
}

// Select adds select props to term store query
func (termGroup *TermGroup) Select(props string) *TermGroup {
	termGroup.selectProps = appendTaxonomyProp(termGroup.selectProps, props)
	return termGroup
}

// Get gets term group metadata
func (termGroup *TermGroup) Get() (map[string]interface{}, error) {
	var props []string
	for _, prop := range termGroup.selectProps {
		propertyXML := prop
		if strings.Index(prop, "<") == -1 {
			propertyXML = fmt.Sprintf(`<Property Name="%s" SelectAll="true" />`, prop)
		}
		props = append(props, propertyXML)
	}

	b := termGroup.csomBuilderEntry()
	b.AddAction(csom.NewQueryWithProps(props), nil)

	return getCSOMResponse(termGroup.client, termGroup.endpoint, termGroup.config, b)
}

// Sets gets term sets object for current term group
func (termGroup *TermGroup) Sets() *TermSets {
	return &TermSets{
		client:    termGroup.client,
		endpoint:  termGroup.endpoint,
		config:    termGroup.config,
		csomEntry: termGroup.csomBuilderEntry(),
	}
}
