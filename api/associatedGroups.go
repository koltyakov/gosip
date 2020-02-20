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
func (groups *AssociatedGroups) Visitors() *Group {
	return NewGroup(
		groups.client,
		fmt.Sprintf("%s/AssociatedVisitorGroup", groups.endpoint),
		groups.config,
	)
}

// Members gets web associated members group API object
func (groups *AssociatedGroups) Members() *Group {
	return NewGroup(
		groups.client,
		fmt.Sprintf("%s/AssociatedMemberGroup", groups.endpoint),
		groups.config,
	)
}

// Owners gets web associated owners group API object
func (groups *AssociatedGroups) Owners() *Group {
	return NewGroup(
		groups.client,
		fmt.Sprintf("%s/AssociatedOwnerGroup", groups.endpoint),
		groups.config,
	)
}
