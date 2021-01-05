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

	return csomRespChildItemsInProp(termGroups.client, termGroups.endpoint, termGroups.config, b, "Groups")
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

// Add creates new group
func (termGroups *TermGroups) Add(name string, guid string) (map[string]interface{}, error) {
	b := termGroups.csomEntry.Clone()
	b.AddObject(csom.NewObjectMethod("CreateGroup", []string{
		fmt.Sprintf(`<Parameter Type="String">%s</Parameter>`, name),
		fmt.Sprintf(`<Parameter Type="String">%s</Parameter>`, guid),
	}), nil)
	b.AddAction(csom.NewQueryWithProps([]string{}), nil)
	return csomResponse(termGroups.client, termGroups.endpoint, termGroups.config, b)
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

	return csomResponse(termGroup.client, termGroup.endpoint, termGroup.config, b)
}

// Delete deletes group object
func (termGroup *TermGroup) Delete() error {
	b := termGroup.csomBuilderEntry().Clone()
	b.AddAction(csom.NewActionMethod("DeleteObject", []string{}), nil)
	_, err := csomResponse(termGroup.client, termGroup.endpoint, termGroup.config, b)
	return err
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
