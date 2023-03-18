package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/auth/addin"
	"github.com/koltyakov/gosip/auth/adfs"
	"github.com/koltyakov/gosip/auth/azurecert"
	"github.com/koltyakov/gosip/auth/azurecreds"
	"github.com/koltyakov/gosip/auth/device"
	"github.com/koltyakov/gosip/auth/fba"
	"github.com/koltyakov/gosip/auth/ntlm"
	"github.com/koltyakov/gosip/auth/saml"
	"github.com/koltyakov/gosip/auth/tmg"
)

// NewAuthByStrategy resolves AuthCnfg object based on strategy name
func NewAuthByStrategy(strategy string) (gosip.AuthCnfg, error) {
	var auth gosip.AuthCnfg

	switch strategy {
	case "azurecert":
		auth = &azurecert.AuthCnfg{}
		break
	case "azurecreds":
		auth = &azurecreds.AuthCnfg{}
		break
	case "device":
		auth = &device.AuthCnfg{}
		break
	case "addin":
		auth = &addin.AuthCnfg{}
		break
	case "adfs":
		auth = &adfs.AuthCnfg{}
		break
	case "fba":
		auth = &fba.AuthCnfg{}
		break
	case "ntlm":
		auth = &ntlm.AuthCnfg{}
		break
	case "saml":
		auth = &saml.AuthCnfg{}
		break
	case "tmg":
		auth = &tmg.AuthCnfg{}
		break
	default:
		return nil, fmt.Errorf("can't resolve the strategy: %s", strategy)
	}

	return auth, nil
}

// NewAuthFromFile resolves AuthCnfg object based on private file
// private.json must contain "strategy" property along with strategy-specific properties
func NewAuthFromFile(privateFile string) (gosip.AuthCnfg, error) {
	jsonFile, err := os.Open(privateFile)
	if err != nil {
		return nil, err
	}
	defer func() { _ = jsonFile.Close() }()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var cnfg struct {
		Strategy string `json:"strategy"`
	}
	if err := json.Unmarshal(byteValue, &cnfg); err != nil {
		return nil, err
	}

	auth, err := NewAuthByStrategy(cnfg.Strategy)
	if err != nil {
		return nil, err
	}

	if err := auth.ParseConfig(byteValue); err != nil {
		return nil, err
	}

	return auth, nil
}
