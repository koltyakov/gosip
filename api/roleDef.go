package api

import (
	"encoding/json"
	"fmt"

	"github.com/koltyakov/gosip"
)

// RoleDefinitions represents SharePoint permissions Role Definitions API queryable object struct
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

// RoleDefInfo - permissions role definition API response payload structure
type RoleDefInfo struct {
	BasePermissions struct {
		High int `json:"High,string"`
		Low  int `json:"Low,string"`
	} `json:"BasePermissions"`
	Description  string `json:"Description"`
	Hidden       bool   `json:"Hidden"`
	ID           int    `json:"Id"`
	Name         string `json:"Name"`
	Order        int    `json:"Order"`
	RoleTypeKind int    `json:"RoleTypeKind"`
}

func getRoleDef(def *RoleDefinitions, endpoint string) (*RoleDefInfo, error) {
	sp := NewHTTPClient(def.client)

	data, err := sp.Post(endpoint, nil, HeadersPresets.Verbose.Headers)
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

// GetByID ...
func (def *RoleDefinitions) GetByID(roleDefID int) (*RoleDefInfo, error) {
	endpoint := fmt.Sprintf("%s/GetById(%d)", def.endpoint, roleDefID)
	return getRoleDef(def, endpoint)
}

// GetByName ...
func (def *RoleDefinitions) GetByName(roleDefName string) (*RoleDefInfo, error) {
	endpoint := fmt.Sprintf("%s/GetByName('%s')", def.endpoint, roleDefName)
	return getRoleDef(def, endpoint)
}

// GetByType ...
func (def *RoleDefinitions) GetByType(roleTypeKind int) (*RoleDefInfo, error) {
	endpoint := fmt.Sprintf("%s/GetByType(%d)", def.endpoint, roleTypeKind)
	return getRoleDef(def, endpoint)
}

// Get ...
func (def *RoleDefinitions) Get() ([]*RoleDefInfo, error) {
	sp := NewHTTPClient(def.client)
	data, err := sp.Get(def.endpoint, getConfHeaders(def.config))
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
