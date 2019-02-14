package scenarios

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/koltyakov/gosip"
)

func spRequest(auth gosip.AuthCnfg, cnfgPath string) {
	startAt := time.Now()

	configPath := resolveCnfgPath(cnfgPath)
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
