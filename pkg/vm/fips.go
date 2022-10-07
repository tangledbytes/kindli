package vm

import (
	"fmt"
	"strings"

	"github.com/utkarsh-pro/kindli/pkg/sh"
)

// FipsCheck returns true if FIPS is enabled for the given VM
func FipsCheck(vmName string) (bool, error) {
	if err := FipsCheckOSSupport(vmName); err != nil {
		return false, err
	}
	resp, err := sh.RunIO("limactl shell " + vmName + " -- cat /proc/sys/crypto/fips_enabled")
	if err != nil {
		return false, fmt.Errorf("failed to check FIPS status: %w", err)
	}

	return strings.TrimSpace(string(resp)) == "1", nil
}

// FipsCheckOSSupport returns true if Kindli supports FIPS in the VMs OS
func FipsCheckOSSupport(vmName string) error {
	resp, err := sh.RunIO("limactl shell " + vmName + " -- sh -c \"cat /etc/os-release | grep '^ID=' | tr -d 'ID='\"")
	if err != nil {
		return fmt.Errorf("failed to get VM OS: %w", err)
	}

	os := strings.TrimSpace(string(resp))
	if os == "debian" {
		return nil
	}

	return fmt.Errorf("FIPS is not supported by Kindli for %s - try creating a new VM", os)
}
