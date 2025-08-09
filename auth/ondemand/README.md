# On-Demand auth flow

During the development, it's common to face a situation when production-level auth (AddIn Onli, Azure AD application) can't be configured in the desired timeframes and no auth strategies work. A simple example might be 2FA (multi-factor authentication) or custom ADFS provider. As a quick workaround, the On-Demand auth can help.

On-Demand means that an interactive browser session is started where a user can provide the credentials as if he/she opens the SharePoint site and follows the same flow as reaching the site in a browser.

In that strategy, the application actually opens the browser and communicates via debug protocol for the auth cookies when uses them in the requests.

On-Demand auth is based on [Lorca](https://github.com/zserge/lorca) project, however, a vital part of the [functionality](https://github.com/zserge/lorca/issues/46) is not exposed as a public API in Lorca, so the dependency is imported from a [fork](https://github.com/koltyakov/lorca) with only that small change in exposing one additional method.

Lorca masters Chrome Debug Protocol, therefore, the Chrome/Chromium browser must be installed in the system where On-Demand auth is intended to be called.

## Auth configuration and usage

```golang
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/api"
	strategy "github.com/koltyakov/gosip-sandbox/strategies/ondemand"
)

func main() {

	authCnfg := &strategy.AuthCnfg{
		SiteURL: os.Getenv("SPAUTH_SITEURL"),
	}

	client := &gosip.SPClient{AuthCnfg: authCnfg}
	sp := api.NewSP(client)

	res, err := sp.Web().Select("Title").Get()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Site title: %s\n", res.Data().Title)

}
```