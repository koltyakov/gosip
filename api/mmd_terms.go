package api

import (
	"context"
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

// Get gets child terms
func (terms *Terms) Get(ctx context.Context) ([]map[string]interface{}, error) {
	b := terms.csomEntry.Clone()
	b.AddAction(csom.NewQueryWithProps([]string{
		`<Property Name="Terms" SelectAll="true" />`,
	}), nil)

	return csomRespChildItemsInProp(ctx, terms.client, terms.endpoint, terms.config, b, "Terms")
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
func (terms *Terms) Add(ctx context.Context, name string, guid string, lcid int) (map[string]interface{}, error) {
	b := terms.csomBuilderEntry().Clone()
	b.AddObject(csom.NewObjectMethod("CreateTerm", []string{
		fmt.Sprintf(`<Parameter Type="String">%s</Parameter>`, name),
		fmt.Sprintf(`<Parameter Type="Number">%d</Parameter>`, lcid),
		fmt.Sprintf(`<Parameter Type="String">%s</Parameter>`, guid),
	}), nil)
	b.AddAction(csom.NewQueryWithProps([]string{}), nil)
	return csomResponse(ctx, terms.client, terms.endpoint, terms.config, b)
}

/* Term */
// API Reference: https://docs.microsoft.com/en-us/dotnet/api/microsoft.sharepoint.taxonomy.term?view=sharepoint-server

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
func (term *Term) Get(ctx context.Context) (map[string]interface{}, error) {
	var props []string
	for _, prop := range term.selectProps {
		propertyXML := prop
		if !strings.Contains(prop, "<") {
			propertyXML = fmt.Sprintf(`<Property Name="%s" SelectAll="true" />`, prop)
		}
		props = append(props, propertyXML)
	}

	b := term.csomBuilderEntry().Clone()
	b.AddAction(csom.NewQueryWithProps(props), nil)

	return csomResponse(ctx, term.client, term.endpoint, term.config, b)
}

// Update sets term's properties
func (term *Term) Update(ctx context.Context, properties map[string]interface{}) (map[string]interface{}, error) {
	b := term.csomBuilderEntry().Clone()
	objects := b.GetObjects() // get parent from all objects
	termObject := objects[len(objects)-1]
	// var scalarProperties []string
	for prop, value := range properties {
		valueXML := fmt.Sprintf("%s", value)
		if !strings.Contains(prop, "<") {
			valueXML = fmt.Sprintf(`<Parameter Type="String">%s</Parameter>`, value)
		}
		b.AddAction(csom.NewSetProperty(prop, valueXML), termObject)
		// scalarProperties = append(scalarProperties, fmt.Sprintf(`<Property Name="%s" ScalarProperty="true" />`, prop))
	}
	b.AddAction(csom.NewQueryWithProps([]string{}), termObject) // scalarProperties
	return csomResponse(ctx, term.client, term.endpoint, term.config, b)
}

// Delete deletes term object
func (term *Term) Delete(ctx context.Context) error {
	b := term.csomBuilderEntry().Clone()
	b.AddAction(csom.NewActionMethod("DeleteObject", []string{}), nil)
	_, err := csomResponse(ctx, term.client, term.endpoint, term.config, b)
	return err
}

// Deprecate deprecates/activates a term
func (term *Term) Deprecate(ctx context.Context, deprecate bool) error {
	b := term.csomBuilderEntry().Clone()
	b.AddAction(csom.NewActionMethod("Deprecate", []string{
		fmt.Sprintf(`<Parameter Type="Boolean">%t</Parameter>`, deprecate),
	}), nil)
	_, err := csomResponse(ctx, term.client, term.endpoint, term.config, b)
	return err
}

// Move moves a term to a new location
// use empty `termGUID` to move to a root term store level
func (term *Term) Move(ctx context.Context, termSetGUID string, termGUID string) error {
	termSetGUID = trimTaxonomyGUID(termSetGUID)
	termGUID = trimTaxonomyGUID(termGUID)

	b := term.csomBuilderEntry().Clone()
	objs := b.GetObjects()

	storeObj := objs[2] // 3rd object is always term store
	childTermObj := objs[len(objs)-1]

	parentObj, _ := b.AddObject(csom.NewObjectMethod("GetTermSet", []string{
		fmt.Sprintf(`<Parameter Type="String">%s</Parameter>`, termSetGUID),
	}), storeObj)

	if len(termGUID) > 0 {
		parentObj, _ = b.AddObject(csom.NewObjectMethod("GetTerm", []string{
			fmt.Sprintf(`<Parameter Type="String">%s</Parameter>`, termGUID),
		}), storeObj)
	}

	_, _ = b.Compile() // generate ID numbers

	b.AddAction(csom.NewActionMethod("Move", []string{
		fmt.Sprintf(`<Parameter ObjectPathId="%d" />`, parentObj.GetID()),
	}), childTermObj)

	_, err := csomResponse(ctx, term.client, term.endpoint, term.config, b)
	return err
}

// Terms gets sub-terms object instance
func (term *Term) Terms() *Terms {
	return &Terms{
		client:   term.client,
		endpoint: term.endpoint,
		config:   term.config,

		csomEntry:   term.csomBuilderEntry().Clone(),
		selectProps: []string{},
	}
}
