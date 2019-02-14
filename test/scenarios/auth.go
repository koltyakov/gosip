package scenarios

import (
	"github.com/koltyakov/gosip/auth/addin"
	"github.com/koltyakov/gosip/auth/adfs"
	"github.com/koltyakov/gosip/auth/basic"
	"github.com/koltyakov/gosip/auth/fba"
	"github.com/koltyakov/gosip/auth/saml"
	"github.com/koltyakov/gosip/auth/tmg"
)

// GetAddinAuthTest : Addin auth test scenario
func GetAddinAuthTest() {
	spRequest(&addin.AuthCnfg{}, "./config/private.addin.json")
}

// GetAdfsAuthTest : ADFS auth test scenario
func GetAdfsAuthTest() {
	spRequest(&adfs.AuthCnfg{}, "./config/private.adfs.json")
}

// GetBasicAuthTest : NTML auth test scenario
func GetBasicAuthTest() {
	spRequest(&basic.AuthCnfg{}, "./config/private.basic.json")
}

// GetFbaAuthTest : FBA auth test scenario
func GetFbaAuthTest() {
	spRequest(&fba.AuthCnfg{}, "./config/private.fba.json")
}

// GetSamlAuthTest : SAML auth test scenario
func GetSamlAuthTest() {
	spRequest(&saml.AuthCnfg{}, "./config/private.saml.json")
}

// GetTmgAuthTest : TMG auth test scenario
func GetTmgAuthTest() {
	spRequest(&tmg.AuthCnfg{}, "./config/private.tmg.json")
}
