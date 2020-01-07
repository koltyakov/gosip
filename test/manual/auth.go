package manual

import (
	"fmt"
	"io/ioutil"
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

// GetTestClient gets a client for a strategy
func GetTestClient(strategy string) (*gosip.SPClient, error) {
	var client *gosip.SPClient
	var err error
	switch strategy {
	case "addin":
		client, err = GetAddinAuthTest()
		break
	case "adfs":
		client, err = GetAdfsAuthTest()
		break
	case "fba":
		client, err = GetFbaAuthTest()
		break
	case "ntlm":
		client, err = GetNtlmAuthTest()
		break
	case "saml":
		client, err = GetSamlAuthTest()
		break
	case "tmg":
		client, err = GetTmgAuthTest()
		break
	default:
		return nil, fmt.Errorf("can't resolve the strategy: %s", strategy)
	}
	return client, err
}

// GetAddinAuthTest : Addin auth test scenario
func GetAddinAuthTest() (*gosip.SPClient, error) {
	return r(&addin.AuthCnfg{}, "./config/private.spo-addin.json")
}

// GetAdfsAuthTest : ADFS auth test scenario
func GetAdfsAuthTest() (*gosip.SPClient, error) {
	return r(&adfs.AuthCnfg{}, "./config/private.onprem-wap.json")
	// return r(&adfs.AuthCnfg{}, "./config/private.onprem-wap-adfs.json")
	// return r(&adfs.AuthCnfg{}, "./config/private.onprem-adfs.json")
}

// GetWapAuthTest : WAP -> Basic auth test scenario
func GetWapAuthTest() (*gosip.SPClient, error) {
	return r(&adfs.AuthCnfg{}, "./config/private.onprem-wap.json")
}

// GetWapAdfsAuthTest : WAP -> ADFS auth test scenario
func GetWapAdfsAuthTest() (*gosip.SPClient, error) {
	return r(&adfs.AuthCnfg{}, "./config/private.onprem-wap-adfs.json")
}

// GetNtlmAuthTest : NTML auth test scenario
func GetNtlmAuthTest() (*gosip.SPClient, error) {
	return r(&ntlm.AuthCnfg{}, "./config/private.onprem-ntlm.json")
}

// GetFbaAuthTest : FBA auth test scenario
func GetFbaAuthTest() (*gosip.SPClient, error) {
	return r(&fba.AuthCnfg{}, "./config/private.onprem-fba.json")
}

// GetSamlAuthTest : SAML auth test scenario
func GetSamlAuthTest() (*gosip.SPClient, error) {
	return r(&saml.AuthCnfg{}, "./config/private.spo-user.json")
}

// GetTmgAuthTest : TMG auth test scenario
func GetTmgAuthTest() (*gosip.SPClient, error) {
	return r(&tmg.AuthCnfg{}, "./config/private.onprem-tmg.json")
}

// GetOnlineADFSTest : SPO ADFS auth test scenario
func GetOnlineADFSTest() (*gosip.SPClient, error) {
	return r(&saml.AuthCnfg{}, "./config/private.spo-adfs.json")
}

func r(auth gosip.AuthCnfg, cnfgPath string) (*gosip.SPClient, error) {
	startAt := time.Now()

	configPath := u.ResolveCnfgPath(cnfgPath)
	err := auth.ReadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("unable to get config: %v", err)
	}

	fmt.Printf("Site Url: %s\n", auth.GetSiteURL())

	client := &gosip.SPClient{AuthCnfg: auth}

	endpoint := auth.GetSiteURL() + "/_api/web?$select=Title"
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create a request: %v", err)
	}

	req.Header.Set("Accept", "application/json;odata=verbose")

	resp, err := client.Execute(req)
	if err != nil {
		return nil, fmt.Errorf("unable to request the api: %v", err)
	}
	defer resp.Body.Close()

	if _, err := ioutil.ReadAll(resp.Body); err != nil {
		return nil, fmt.Errorf("unable to read api response: %v", err)
	}

	// fmt.Printf("response: %s\n", data)
	fmt.Printf("connection established in %f seconds\n", time.Since(startAt).Seconds())
	fmt.Printf("below is the results of manual tests\n\n")

	return client, nil
}
