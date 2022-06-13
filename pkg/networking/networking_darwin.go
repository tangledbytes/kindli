package networking

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/utkarsh-pro/kindli/pkg/sh"
)

func Setup(vmName string) error {
	srcIP, err := sh.RunIO("ifconfig bridge100 | grep 'inet ' | cut -d' ' -f2")
	if err != nil {
		return fmt.Errorf("failed to get bridge100 interface IPv4 address: %s", err)
	}
	logrus.Debug("Bridge100 IP:", srcIP)

	limaIPAddr, err := getLimaVMIPAddress(vmName)
	if err != nil {
		return fmt.Errorf("failed to get VM internal interace IP address: %s", err)
	}
	logrus.Debug("Lima IP Address:", limaIPAddr)

	if err := sh.Run(fmt.Sprintf("sudo route -nv add -net 172.18 %s", trim([]byte(limaIPAddr)))); err != nil {
		return fmt.Errorf("failed to setup route from system to VM: %s", err)
	}

	kindIf, err := sh.RunIO(fmt.Sprintf("limactl shell %s -- ip -o link show | awk -F': ' '{print $2}' | grep 'br-'", vmName))
	if err != nil {
		return fmt.Errorf("failed to get kind network interface name: %s", err)
	}

	destNet := "172.18.0.0/16"
	hostIf := "lima0"

	if err := sh.Run(fmt.Sprintf("limactl shell %s -- sudo iptables -t filter -A FORWARD -4 -p tcp -s %s -d %s -j ACCEPT -i %s -o %s", vmName, trim(srcIP), destNet, hostIf, trim(kindIf))); err != nil {
		return fmt.Errorf("failed to setup route from VM network interface to kind network interface")
	}

	return nil
}

func getLimaVMIPAddress(vmName string) (string, error) {
	// Try on lima0 interface
	limaIPAddr, err := sh.RunIO(fmt.Sprintf("limactl shell %s -- ip -o -4 a s | grep lima0 | grep -E -o 'inet [0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}' | cut -d' ' -f2", vmName))
	if err == nil && string(limaIPAddr) != "" {
		return string(limaIPAddr), nil
	}

	// Try on the eth0 interface
	limaIPAddr, err = sh.RunIO(fmt.Sprintf("limactl shell %s -- ip -o -4 a s | grep eth0 | grep -E -o 'inet [0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}' | cut -d' ' -f2", vmName))
	if err == nil && string(limaIPAddr) != "" {
		return string(limaIPAddr), nil
	}

	return "", fmt.Errorf("failed to get VM internal interace")
}

func trim(data []byte) string {
	return strings.Trim(string(data), " \n")
}
