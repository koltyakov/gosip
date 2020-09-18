package api

import (
	"fmt"

	"github.com/koltyakov/gosip"
)

//go:generate ggen -ent AssociatedGroups -conf -coll

// AssociatedGroups web associated groups scope constructor
// Always use NewAssociatedGroups constructor instead of &AssociatedGroups{}
type AssociatedGroups struct {
	client   *gosip.SPClient
	config   *RequestConfig
	endpoint string
}

// NewAssociatedGroups - AssociatedGroups struct constructor function
func NewAssociatedGroups(client *gosip.SPClient, endpoint string, config *RequestConfig) *AssociatedGroups {
	return &AssociatedGroups{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// Visitors gets web associated visitors group API object
func (associatedGroups *AssociatedGroups) Visitors() *Group {
	return NewGroup(
		associatedGroups.client,
		fmt.Sprintf("%s/AssociatedVisitorGroup", associatedGroups.endpoint),
		associatedGroups.config,
	)
}

// Members gets web associated members group API object
func (associatedGroups *AssociatedGroups) Members() *Group {
	return NewGroup(
		associatedGroups.client,
		fmt.Sprintf("%s/AssociatedMemberGroup", associatedGroups.endpoint),
		associatedGroups.config,
	)
}

// Owners gets web associated owners group API object
func (associatedGroups *AssociatedGroups) Owners() *Group {
	return NewGroup(
		associatedGroups.client,
		fmt.Sprintf("%s/AssociatedOwnerGroup", associatedGroups.endpoint),
		associatedGroups.config,
	)
}
