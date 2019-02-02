package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	s "github.com/koltyakov/gosip/test/scenarios"
)

func main() {

	// s.CpassDummyTest()
	// s.CpassAutoModeTest()
	// s.ConfigReaderTest()

	// httpRequestsTest()
	// restRequestTest()

	// s.ConfigReader2Test()

	s.GetAuthTest()

}

func httpRequestsTest() {
	resp, err := http.Get("https://www.arvosys.com")
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))
}

func restRequestTest() {
	endpoint := "https://spnode.sharepoint.com/sites/surveys/_api/web?$select=Title"

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Accept", "application/json;odata=minimalmetadata")
	req.Header.Set("Cookie", "rtFa=o+Ly6hUF4Hawt7xwt8D61SLv9y3D0KMxaON0Jud4AZQmNjdDRjIwOTItNzhFQS00OEM0LTlBRkItMkUzQkU3NkNCMDFDUq4/cQzo6sYVfOC022EXIZEUWPkQJFH59/6QHECiVf/RVx5J7jG5ry6VrQ80Zie1ozRd8grU8xhOZP6XU0Goeu1b0SSIb7Mu2P/VWZlWsG6aHXYoCZLAQhXzCyUebvbT0dQtH8fmHmyh+w/efOBCHDjtdJZTfp8votypBBTyLDFt8EtY2oVE3InOG5+A8i3qc0/H/dLtNqgVaRn9XqqoNZX+dwtQksCJYvNxXEnXc713DKfedY7vViz1EnAhQvoS9pPZPIdy3DcdpUiDChJiRQmC/6cSCrREi+eIqJe3kGhel2mfYTykq9HC0O8eiYM88eQ4KX2pbz7aii3Qhr1XMEUAAAA=; FedAuth=77u/PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0idXRmLTgiPz48U1A+VjUsMGguZnxtZW1iZXJzaGlwfDEwMDNiZmZkYTQ2ODg5MDBAbGl2ZS5jb20sMCMuZnxtZW1iZXJzaGlwfGFuZHJldy5rb2x0eWFrb3ZAc3Bub2RlLm9ubWljcm9zb2Z0LmNvbSwxMzE5MDY1MjcyODAwMDAwMDAsMTMxODY5NDAzMjgwMDAwMDAwLDEzMTkzNDM4NjkxNTExMjQzMywwLjAuMC4wLDMsNjdjZjIwOTItNzhlYS00OGM0LTlhZmItMmUzYmU3NmNiMDFjLCxWMiExMDAzQkZGREE0Njg4OTAwITEzMTkwNjUyNzI4LGI1MTRiNTllLWIwZGUtNzAwMC02MDM0LTFkNDJlNjMxODFiZCwwYWRmYjk5ZS03MDQ5LTcwMDAtNjAzNC0xMzdlZWUyM2MxMTcsLDAsMTMxOTE3MjQzOTE1NjI5MTIwLDEzMTkxOTc5OTkxNTYyOTEyMCwsQlJwRmdYaFRmcllLbktJVFlBTlV3ckZvUzlEZWszdndidWRQOXo1clhXNklYclFnL1UreXo5bkFzc25MWllMYU5xc3l1SnBMeXVycnFoYUNFc0ZrNGRXQmJNRkx0TFJVR1hmaHZrTW1jd0Yvc2RWU3BxTGsxdmNtRVYrUmpIQkpHTUVqNVBrVUhFV0ZkYUFrRFNyY2ZTckd3WnBPN3hHVUNjWTlxUkRSYTBHRENJNG1RZEJ5Zk5oSnc1VWVncWtHd05hcm5FZWVRc1RaQWFTOVVsajN6cXpLai9JUXR0WHdDaDduVmVTTWE0azNoVXo2VHk5RllaSDU2enVoblFBcm1RTGR2Rmd3N1dFQWllQTZZZEEwNDhCRXN5NzZQc051encxTXVzYWtCZ0ptY3ZDYXJUN0lnZlpBWTM5MHh6bCt0elhnZGtDaWswM3BHTWhGcVB4Uk9nPT08L1NQPg==;")

	// Do the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	// log.Println(result)
	log.Println(result["Title"])

}
