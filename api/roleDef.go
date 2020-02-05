package api

import (
	"encoding/json"
	"fmt"

	"github.com/koltyakov/gosip"
)

// RoleDefinitions represents SharePoint permissions Role Definitions API queryable object struct
// Always use NewRoleDefinitions constructor instead of &RoleDefinitions{}
type RoleDefinitions struct {
	client   *gosip.SPClient
	config   *RequestConfig
	endpoint string
}

// NewRoleDefinitions - RoleDefinitions struct constructor function
func NewRoleDefinitions(client *gosip.SPClient, endpoint string, config *RequestConfig) *RoleDefinitions {
	return &RoleDefinitions{
		client:   client,
		endpoint: endpoint,
		config:   config,
	}
}

// BasePermissions - Low/High pair of base permissions
type BasePermissions struct {
	High int `json:"High,string"`
	Low  int `json:"Low,string"`
}

// RoleDefInfo - permissions role definition API response payload structure
type RoleDefInfo struct {
	BasePermissions *BasePermissions `json:"BasePermissions"`
	Description     string           `json:"Description"`
	Hidden          bool             `json:"Hidden"`
	ID              int              `json:"Id"`
	Name            string           `json:"Name"`
	Order           int              `json:"Order"`
	RoleTypeKind    int              `json:"RoleTypeKind"`
}

// GetByID gets a role definition by its ID
func (def *RoleDefinitions) GetByID(roleDefID int) (*RoleDefInfo, error) {
	endpoint := fmt.Sprintf("%s/GetById(%d)", def.endpoint, roleDefID)
	return getRoleDef(def, endpoint)
}

// GetByName gets a role definition by its Name
func (def *RoleDefinitions) GetByName(roleDefName string) (*RoleDefInfo, error) {
	endpoint := fmt.Sprintf("%s/GetByName('%s')", def.endpoint, roleDefName)
	return getRoleDef(def, endpoint)
}

// GetByType gets a role definition by its RoleTypeKinds
func (def *RoleDefinitions) GetByType(roleTypeKind int) (*RoleDefInfo, error) {
	endpoint := fmt.Sprintf("%s/GetByType(%d)", def.endpoint, roleTypeKind)
	return getRoleDef(def, endpoint)
}

// Get gets a collection of available role definitions
func (def *RoleDefinitions) Get() ([]*RoleDefInfo, error) {
	client := NewHTTPClient(def.client)
	data, err := client.Get(def.endpoint, def.config)
	if err != nil {
		return nil, err
	}

	res := &struct {
		D struct {
			Results []*RoleDefInfo `json:"results"`
		} `json:"d"`
	}{}

	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}

	return res.D.Results, nil
}

func getRoleDef(def *RoleDefinitions, endpoint string) (*RoleDefInfo, error) {
	client := NewHTTPClient(def.client)

	data, err := client.Post(endpoint, nil, patchConfigHeaders(def.config, HeadersPresets.Verbose.Headers))
	if err != nil {
		return nil, err
	}

	res := &struct {
		RoleDefInfo *RoleDefInfo `json:"d"`
	}{}

	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}

	return res.RoleDefInfo, nil
}
