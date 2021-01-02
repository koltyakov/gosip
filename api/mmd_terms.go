package api

import (
	"fmt"
	"strings"

	"github.com/koltyakov/gosip/csom"
)

// Terms struct
type Terms struct {
	client   *HTTPClient
	config   *RequestConfig
	endpoint string

	csomEntry   csom.Builder
	selectProps []string
}

// csomBuilderEntry gets CSOM builder entry
func (terms *Terms) csomBuilderEntry() csom.Builder {
	b := terms.csomEntry.Clone()
	return b
}

// Select adds select props to terms collection query
func (terms *Terms) Select(props string) *Terms {
	terms.selectProps = appendTaxonomyProp(terms.selectProps, props)
	return terms
}

// GetAll gets all terms collection metadata
func (terms *Terms) GetAll() ([]map[string]interface{}, error) {
	var props []string
	for _, prop := range terms.selectProps {
		propertyXML := prop
		if strings.Index(prop, "<") == -1 {
			propertyXML = fmt.Sprintf(`<Property Name="%s" SelectAll="true" />`, prop)
		}
		props = append(props, propertyXML)
	}

	b := terms.csomBuilderEntry().Clone()
	b.AddObject(csom.NewObjectMethod("GetAllTerms", []string{}), nil)
	b.AddAction(csom.NewQueryWithChildProps(props), nil)

	termsResp, err := getCSOMResponse(terms.client, terms.endpoint, terms.config, b)
	if err != nil {
		return nil, err
	}

	items, ok := termsResp["_Child_Items_"].([]interface{})
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

// GetByID gets a term by its GUID
func (terms *Terms) GetByID(termGUID string) *Term {
	return &Term{
		client:      terms.client,
		config:      terms.config,
		endpoint:    terms.endpoint,
		id:          trimTaxonomyGUID(termGUID),
		csomEntry:   terms.csomEntry,
		selectProps: []string{},
	}
}

// Add creates new term
func (terms *Terms) Add(name string, lang int, guid string) (map[string]interface{}, error) {
	b := terms.csomBuilderEntry().Clone()
	b.AddObject(csom.NewObjectMethod("CreateTerm", []string{
		fmt.Sprintf(`<Parameter Type="String">%s</Parameter>`, name),
		fmt.Sprintf(`<Parameter Type="Number">%d</Parameter>`, lang),
		fmt.Sprintf(`<Parameter Type="String">%s</Parameter>`, guid),
	}), nil)
	b.AddAction(csom.NewQueryWithProps([]string{}), nil)

	return getCSOMResponse(terms.client, terms.endpoint, terms.config, b)
}

/* Term */

// Term struct
type Term struct {
	client   *HTTPClient
	config   *RequestConfig
	endpoint string

	id string

	csomEntry   csom.Builder
	selectProps []string
}

// csomBuilderEntry gets CSOM builder entry
func (term *Term) csomBuilderEntry() csom.Builder {
	b := term.csomEntry.Clone()
	b.AddObject(csom.NewObjectMethod("GetTerm", []string{
		fmt.Sprintf(`<Parameter Type="String">%s</Parameter>`, term.id),
	}), nil)
	return b
}

// Select adds select props to term query
func (term *Term) Select(props string) *Term {
	term.selectProps = appendTaxonomyProp(term.selectProps, props)
	return term
}

// Get gets term metadata
func (term *Term) Get() (map[string]interface{}, error) {
	var props []string
	for _, prop := range term.selectProps {
		propertyXML := prop
		if strings.Index(prop, "<") == -1 {
			propertyXML = fmt.Sprintf(`<Property Name="%s" SelectAll="true" />`, prop)
		}
		props = append(props, propertyXML)
	}

	b := term.csomBuilderEntry().Clone()
	b.AddAction(csom.NewQueryWithProps(props), nil)

	return getCSOMResponse(term.client, term.endpoint, term.config, b)
}

// Update sets term's properties
func (term *Term) Update(properties map[string]interface{}) (map[string]interface{}, error) {
	b := term.csomBuilderEntry().Clone()
	objects := b.GetObjects() // get parent from all objects
	termObject := objects[len(objects)-1]
	var scalarProperties []string
	for prop, value := range properties {
		valueXML := fmt.Sprintf("%s", value)
		if strings.Index(prop, "<") == -1 {
			valueXML = fmt.Sprintf(`<Parameter Type="String">%s</Parameter>`, value)
		}
		b.AddAction(csom.NewSetProperty(prop, valueXML), termObject)
		scalarProperties = append(scalarProperties, fmt.Sprintf(`<Property Name="%s" ScalarProperty="true" />`, prop))
	}
	b.AddAction(csom.NewQueryWithProps(scalarProperties), termObject)
	return getCSOMResponse(term.client, term.endpoint, term.config, b)
}

// Delete deletes term object
func (term *Term) Delete() error {
	b := term.csomBuilderEntry().Clone()
	b.AddAction(csom.NewActionMethod("DeleteObject", []string{}), nil)
	_, err := getCSOMResponse(term.client, term.endpoint, term.config, b)
	return err
}
