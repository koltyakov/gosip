package cpass

import (
	"testing"
)

func TestGetMachineID(t *testing.T) {
	machineID, err := getMachineID(true)
	if err != nil {
		t.Error(err)
	}
	if machineID == "" {
		t.Error("Got empty machine id")
	}
}
