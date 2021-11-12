package helpers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Azure/go-ntlmssp"
	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/api"
	u "github.com/koltyakov/gosip/test/utils"
)

func CheckTransport(auth gosip.AuthCnfg, cnfgPath string) error {
	err := auth.ReadConfig(u.ResolveCnfgPath(cnfgPath))
	if err != nil {
		return err
	}

	client := &gosip.SPClient{
		AuthCnfg: auth,
		Client: http.Client{
			Transport: &http.Transport{TLSHandshakeTimeout: 25 * time.Second},
		},
	}

	sp := api.NewSP(client)
	if _, err := sp.ContextInfo(); err != nil {
		return fmt.Errorf("can't get SP context: %s", err)
	}

	if auth.GetStrategy() == "ntlm" {
		n, ok := client.Transport.(ntlmssp.Negotiator)
		if !ok {
			return fmt.Errorf("transport configuration leak")
		}

		tr, ok := n.RoundTripper.(*http.Transport)
		if !ok {
			return fmt.Errorf("transport configuration leak")
		}

		if tr.TLSHandshakeTimeout != 25*time.Second {
			return fmt.Errorf("transport configuration leak")
		}

		return nil
	}

	// None NTLM strategies

	tr, ok := client.Transport.(*http.Transport)
	if !ok {
		return fmt.Errorf("transport configuration leak")
	}

	if tr.TLSHandshakeTimeout != 25*time.Second {
		return fmt.Errorf("transport configuration leak")
	}

	return nil
}
