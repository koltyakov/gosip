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

	return csomRespChildItemsInProp(termSets.client, termSets.endpoint, termSets.config, b, "TermSets")
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

// GetByName gets term sets by a name and LCID, searches within term store
func (termSets *TermSets) GetByName(termSetName string, lcid int) ([]map[string]interface{}, error) {
	b := csom.NewBuilder()
	objs := termSets.csomBuilderEntry().GetObjects()

	b.AddObject(csom.NewObject(objs[1].Template()), nil) // GetTaxonomySession
	b.AddObject(csom.NewObject(objs[2].Template()), nil) // GetDefaultSiteCollectionTermStore or TermStores
	if strings.Contains(objs[2].Template(), "TermStores") {
		b.AddObject(csom.NewObject(objs[3].Template()), nil) // GetById or GetByName
	}

	b.AddObject(csom.NewObjectMethod("GetTermSetsByName", []string{
		fmt.Sprintf(`<Parameter Type="String">%s</Parameter>`, termSetName),
		fmt.Sprintf(`<Parameter Type="Number">%d</Parameter>`, lcid),
	}), nil)
	b.AddAction(csom.NewQueryWithChildProps([]string{}), nil)
	return csomRespChildItems(termSets.client, termSets.endpoint, termSets.config, b)
}

// Add creates new term set
func (termSets *TermSets) Add(name string, guid string, lcid int) (map[string]interface{}, error) {
	b := termSets.csomBuilderEntry().Clone()
	b.AddObject(csom.NewObjectMethod("CreateTermSet", []string{
		fmt.Sprintf(`<Parameter Type="String">%s</Parameter>`, name),
		fmt.Sprintf(`<Parameter Type="String">%s</Parameter>`, guid),
		fmt.Sprintf(`<Parameter Type="Number">%d</Parameter>`, lcid),
	}), nil)
	b.AddAction(csom.NewQueryWithProps([]string{}), nil)
	return csomResponse(termSets.client, termSets.endpoint, termSets.config, b)
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
		if !strings.Contains(prop, "<") {
			propertyXML = fmt.Sprintf(`<Property Name="%s" SelectAll="true" />`, prop)
		}
		props = append(props, propertyXML)
	}

	b := termSet.csomBuilderEntry().Clone()
	b.AddAction(csom.NewQueryWithProps(props), nil)

	return csomResponse(termSet.client, termSet.endpoint, termSet.config, b)
}

// Delete deletes term set object
func (termSet *TermSet) Delete() error {
	b := termSet.csomBuilderEntry().Clone()
	b.AddAction(csom.NewActionMethod("DeleteObject", []string{}), nil)
	_, err := csomResponse(termSet.client, termSet.endpoint, termSet.config, b)
	return err
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

// GetAllTerms gets all terms collection metadata
func (termSet *TermSet) GetAllTerms() ([]map[string]interface{}, error) {
	var props []string
	for _, prop := range termSet.selectProps {
		propertyXML := prop
		if !strings.Contains(prop, "<") {
			propertyXML = fmt.Sprintf(`<Property Name="%s" SelectAll="true" />`, prop)
		}
		props = append(props, propertyXML)
	}

	b := termSet.csomBuilderEntry().Clone()
	b.AddObject(csom.NewObjectMethod("GetAllTerms", []string{}), nil)
	b.AddAction(csom.NewQueryWithChildProps(props), nil)

	return csomRespChildItems(termSet.client, termSet.endpoint, termSet.config, b)
}
