package manual

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/auth/addin"
	"github.com/koltyakov/gosip/auth/adfs"
	"github.com/koltyakov/gosip/auth/basic"
	"github.com/koltyakov/gosip/auth/fba"
	"github.com/koltyakov/gosip/auth/saml"
	"github.com/koltyakov/gosip/auth/tmg"
	u "github.com/koltyakov/gosip/test/utils"
)

// GetAddinAuthTest : Addin auth test scenario
func GetAddinAuthTest() {
	r(&addin.AuthCnfg{}, "./config/private.addin.json")
}

// GetAdfsAuthTest : ADFS auth test scenario
func GetAdfsAuthTest() {
	r(&adfs.AuthCnfg{}, "./config/private.adfs.json")
}

// GetWapAuthTest : WAP -> Basic auth test scenario
func GetWapAuthTest() {
	r(&adfs.AuthCnfg{}, "./config/private.wap.json")
}

// GetWapAdfsAuthTest : WAP -> ADFS auth test scenario
func GetWapAdfsAuthTest() {
	r(&adfs.AuthCnfg{}, "./config/private.wap-adfs.json")
}

// GetBasicAuthTest : NTML auth test scenario
func GetBasicAuthTest() {
	r(&basic.AuthCnfg{}, "./config/private.basic.json")
}

// GetFbaAuthTest : FBA auth test scenario
func GetFbaAuthTest() {
	r(&fba.AuthCnfg{}, "./config/private.fba.json")
}

// GetSamlAuthTest : SAML auth test scenario
func GetSamlAuthTest() {
	r(&saml.AuthCnfg{}, "./config/private.saml.json")
}

// GetTmgAuthTest : TMG auth test scenario
func GetTmgAuthTest() {
	r(&tmg.AuthCnfg{}, "./config/private.tmg.json")
}

func r(auth gosip.AuthCnfg, cnfgPath string) {
	startAt := time.Now()

	configPath := u.ResolveCnfgPath(cnfgPath)
	err := auth.ReadConfig(configPath)
	if err != nil {
		log.Fatalf("unable to get config: %v", err)
	}

	fmt.Printf("siteUrl: %s\n", auth.GetSiteURL())

	client := &gosip.SPClient{AuthCnfg: auth}

	apiEndpoint := auth.GetSiteURL() + "/_api/web?$select=Title"
	req, err := http.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		log.Fatalf("unable to create a request: %v", err)
	}

	req.Header.Set("Accept", "application/json;odata=verbose")

	resp, err := client.Execute(req)
	if err != nil {
		log.Fatalf("unable to request api: %v", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("unable to read a response: %v", err)
	}

	fmt.Printf("response: %s\n", data)
	fmt.Printf("time taken, sec: %f\n", time.Since(startAt).Seconds())
}
