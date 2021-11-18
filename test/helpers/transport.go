package helpers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Azure/go-ntlmssp"
	"github.com/koltyakov/gosip"
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

	if _, err := gosip.GetDigest(context.Background(), client); err != nil {
		return fmt.Errorf("unable to get digest: %w", err)
	}

	if _, _, err := client.AuthCnfg.GetAuth(); err != nil {
		return err
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
