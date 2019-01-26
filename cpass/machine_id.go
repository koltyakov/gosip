package cpass

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/denisbrodbeck/machineid" // port cpass's implementation
)

func getMachineID(original bool) string {
	guid, err := machineid.ID()
	if err != nil {
		panic(err)
	}
	guid = strings.ToLower(guid)
	if !original {
		hasher := sha256.New()
		hasher.Write([]byte(guid))
		// guid = fmt.Sprintf("%x", hasher.Sum(nil))
		guid = hex.EncodeToString(hasher.Sum(nil))
	}
	return guid
}
