package api

import (
	"fmt"
	"strings"

	"github.com/koltyakov/gosip/csom"
)

// TermSets term sets struct
type TermSets struct {
	client   *HTTPClient
	config   *RequestConfig
	endpoint string

	csomEntry csom.Builder
}

// csomBuilderEntry gets CSOM builder entry
func (termSets *TermSets) csomBuilderEntry() csom.Builder {
	b := termSets.csomEntry.Clone()
	return b
}

// Get gets term sets metadata
func (termSets *TermSets) Get() ([]map[string]interface{}, error) {
	b := termSets.csomBuilderEntry().Clone()
	b.AddAction(csom.NewQueryWithProps([]string{
		`<Property Name="TermSets" SelectAll="true" />`,
	}), nil)

	csomResp, err := getCSOMResponse(termSets.client, termSets.endpoint, termSets.config, b)
	if err != nil {
		return nil, err
	}

	groups, ok := csomResp["TermSets"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("can't get term sets from term group")
	}

	items, ok := groups["_Child_Items_"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("can't get child items from term sets")
	}

	var resItems []map[string]interface{}
	for _, item := range items {
		resItem, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("can't get item from term set")
		}
		resItems = append(resItems, resItem)
	}

	return resItems, nil
}

// GetByID gets term set object by its GUID
func (termSets *TermSets) GetByID(setGUID string) *TermSet {
	return &TermSet{
		client:   termSets.client,
		endpoint: termSets.endpoint,
		config:   termSets.config,

		id: trimTaxonomyGUID(setGUID),

		csomEntry:   termSets.csomBuilderEntry().Clone(),
		selectProps: []string{},
	}
}

/* Term Sets */

// TermSet term set struct
type TermSet struct {
	client   *HTTPClient
	config   *RequestConfig
	endpoint string

	id string

	csomEntry   csom.Builder
	selectProps []string
}

// csomBuilderEntry gets CSOM builder entry
func (termSet *TermSet) csomBuilderEntry() csom.Builder {
	b := termSet.csomEntry.Clone()
	b.AddObject(csom.NewObjectMethod("GetTermSet", []string{
		fmt.Sprintf(`<Parameter Type="String">%s</Parameter>`, termSet.id),
	}), nil)
	return b
}

// Select adds select props to term set query
func (termSet *TermSet) Select(props string) *TermSet {
	termSet.selectProps = appendTaxonomyProp(termSet.selectProps, props)
	return termSet
}

// Get gets term set metadata
func (termSet *TermSet) Get() (map[string]interface{}, error) {
	var props []string
	for _, prop := range termSet.selectProps {
		propertyXML := prop
		if strings.Index(prop, "<") == -1 {
			propertyXML = fmt.Sprintf(`<Property Name="%s" SelectAll="true" />`, prop)
		}
		props = append(props, propertyXML)
	}

	b := termSet.csomBuilderEntry().Clone()
	b.AddAction(csom.NewQueryWithProps(props), nil)

	return getCSOMResponse(termSet.client, termSet.endpoint, termSet.config, b)
}

// Terms gets terms object instance
func (termSet *TermSet) Terms() *Terms {
	return &Terms{
		client:   termSet.client,
		endpoint: termSet.endpoint,
		config:   termSet.config,

		csomEntry:   termSet.csomBuilderEntry().Clone(),
		selectProps: []string{},
	}
}
