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
	"github.com/koltyakov/gosip/auth/fba"
	"github.com/koltyakov/gosip/auth/ntlm"
	"github.com/koltyakov/gosip/auth/saml"
	"github.com/koltyakov/gosip/auth/tmg"
	u "github.com/koltyakov/gosip/test/utils"
)

// GetAddinAuthTest : Addin auth test scenario
func GetAddinAuthTest() *gosip.SPClient {
	return r(&addin.AuthCnfg{}, "./config/private.spo-addin.json")
}

// GetAdfsAuthTest : ADFS auth test scenario
func GetAdfsAuthTest() *gosip.SPClient {
	return r(&adfs.AuthCnfg{}, "./config/private.onprem-wap.json")
	// return r(&adfs.AuthCnfg{}, "./config/private.onprem-wap-adfs.json")
	// return r(&adfs.AuthCnfg{}, "./config/private.onprem-adfs.json")
}

// GetWapAuthTest : WAP -> Basic auth test scenario
func GetWapAuthTest() *gosip.SPClient {
	return r(&adfs.AuthCnfg{}, "./config/private.onprem-wap.json")
}

// GetWapAdfsAuthTest : WAP -> ADFS auth test scenario
func GetWapAdfsAuthTest() *gosip.SPClient {
	return r(&adfs.AuthCnfg{}, "./config/private.onprem-wap-adfs.json")
}

// GetNtlmAuthTest : NTML auth test scenario
func GetNtlmAuthTest() *gosip.SPClient {
	return r(&ntlm.AuthCnfg{}, "./config/private.onprem-ntlm.json")
}

// GetFbaAuthTest : FBA auth test scenario
func GetFbaAuthTest() *gosip.SPClient {
	return r(&fba.AuthCnfg{}, "./config/private.onprem-fba.json")
}

// GetSamlAuthTest : SAML auth test scenario
func GetSamlAuthTest() *gosip.SPClient {
	return r(&saml.AuthCnfg{}, "./config/private.spo-user.json")
}

// GetTmgAuthTest : TMG auth test scenario
func GetTmgAuthTest() *gosip.SPClient {
	return r(&tmg.AuthCnfg{}, "./config/private.onprem-tmg.json")
}

// GetOnlineADFSTest : SPO ADFS auth test scenario
func GetOnlineADFSTest() *gosip.SPClient {
	return r(&saml.AuthCnfg{}, "./config/private.spo-adfs.json")
}

func r(auth gosip.AuthCnfg, cnfgPath string) *gosip.SPClient {
	startAt := time.Now()

	configPath := u.ResolveCnfgPath(cnfgPath)
	err := auth.ReadConfig(configPath)
	if err != nil {
		log.Fatalf("unable to get config: %v", err)
	}

	fmt.Printf("siteUrl: %s\n", auth.GetSiteURL())

	client := &gosip.SPClient{AuthCnfg: auth}

	endpoint := auth.GetSiteURL() + "/_api/web?$select=Title"
	req, err := http.NewRequest("GET", endpoint, nil)
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

	return client
}
