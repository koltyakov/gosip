package cpass

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/denisbrodbeck/machineid" // port cpass's implementation
)

func getMachineID(original bool) (machineID string, err error) {
	machineID, err = machineid.ID()
	if err != nil {
		return
	}
	machineID = strings.ToLower(machineID)
	if !original {
		hasher := sha256.New()
		hasher.Write([]byte(machineID))
		machineID = hex.EncodeToString(hasher.Sum(nil))
	}
	return
}
