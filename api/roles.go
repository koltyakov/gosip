package api

import (
	"fmt"

	"github.com/koltyakov/gosip"
)

// RoleTypeKinds defined standard role type kinds
var RoleTypeKinds = struct {
	None          int
	Guest         int
	Reader        int
	Contributor   int
	WebDesigner   int
	Administrator int
	Editor        int
	System        int
}{
	None:          0,
	Guest:         1,
	Reader:        2,
	Contributor:   3,
	WebDesigner:   4,
	Administrator: 5,
	Editor:        6,
	System:        7,
}

// Roles  represent SharePoint Permissions Roles API queryable collection struct
// Always use NewRoles constructor instead of &Roles{}
type Roles struct {
	client   *gosip.SPClient
	config   *RequestConfig
	endpoint string
}

// NewRoles - Roles struct constructor function
func NewRoles(client *gosip.SPClient, endpoint string, config *RequestConfig) *Roles {
	return &Roles{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// ResetInheritance ...
func (permissions *Roles) ResetInheritance() error {
	sp := NewHTTPClient(permissions.client)
	endpoint := fmt.Sprintf("%s/ResetRoleInheritance", permissions.endpoint)
	_, err := sp.Post(endpoint, nil, getConfHeaders(permissions.config))
	return err
}

// BreakInheritance ...
func (permissions *Roles) BreakInheritance(copyRoleAssigments bool, clearSubScopes bool) error {
	sp := NewHTTPClient(permissions.client)
	endpoint := fmt.Sprintf(
		"%s/BreakRoleInheritance(copyroleassignments=%t,clearsubscopes=%t)",
		permissions.endpoint,
		copyRoleAssigments,
		clearSubScopes,
	)
	_, err := sp.Post(endpoint, nil, getConfHeaders(permissions.config))
	return err
}

// AddAssigment ...
func (permissions *Roles) AddAssigment(principalID int, roleDefID int) error {
	sp := NewHTTPClient(permissions.client)
	endpoint := fmt.Sprintf(
		"%s/RoleAssignments/AddRoleAssignment(principalid=%d,roledefid=%d)",
		permissions.endpoint,
		principalID,
		roleDefID,
	)
	_, err := sp.Post(endpoint, nil, getConfHeaders(permissions.config))
	return err
}

// RemoveAssigment ...
func (permissions *Roles) RemoveAssigment(principalID int, roleDefID int) error {
	sp := NewHTTPClient(permissions.client)
	endpoint := fmt.Sprintf(
		"%s/RoleAssignments/RemoveRoleAssignment(principalid=%d,roledefid=%d)",
		permissions.endpoint,
		principalID,
		roleDefID,
	)
	_, err := sp.Post(endpoint, nil, getConfHeaders(permissions.config))
	return err
}
