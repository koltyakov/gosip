package helpers

import (
	"fmt"

	"github.com/koltyakov/gosip"
	u "github.com/koltyakov/gosip/test/utils"
)

// CheckDigest : check getting form digest
func CheckDigest(auth gosip.AuthCnfg, cnfgPath string) error {
	err := auth.ReadConfig(u.ResolveCnfgPath(cnfgPath))
	if err != nil {
		return err
	}

	client := &gosip.SPClient{
		AuthCnfg: auth,
	}

	digest, err := gosip.GetDigest(client)
	if err != nil {
		return fmt.Errorf("unable to get digest: %v", err)
	}

	if digest == "" {
		return fmt.Errorf("got empty digest")
	}

	cachedDigest, err := gosip.GetDigest(client)
	if err != nil {
		return fmt.Errorf("unable to get cached digest: %v", err)
	}

	if digest != cachedDigest {
		return fmt.Errorf("digest cache is broken")
	}

	return nil
}
