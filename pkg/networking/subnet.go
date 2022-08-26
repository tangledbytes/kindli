package networking

import (
	"fmt"
	"net"
	"strings"

	"github.com/utkarsh-pro/kindli/pkg/docker"
)

// GetIPv4SubnetPrefix returns the subnet prefix for the given docker network
func GetIPv4SubnetPrefix(network string) (string, error) {
	subnet, err := docker.NetworkInspect(network, "{{ index .IPAM.Config 0 \"Subnet\"}}")
	if err != nil {
		return "", fmt.Errorf("failed to get subnet prefix: %w", err)
	}

	ip, sub, err := net.ParseCIDR(subnet)
	if err != nil {
		return "", fmt.Errorf("failed to get subnet prefix: %w", err)
	}

	size, _ := sub.Mask.Size()
	if size != 16 {
		return "", fmt.Errorf("ipv4 subnet of only size 16 is supported")
	}

	return strings.Join(strings.Split(ip.String(), ".")[:2], "."), nil
}

// GetIPv6SubnetPrefix returns the subnet prefix for the given docker network
func GetIPv6SubnetPrefix(network string) (string, error) {
	subnet, err := docker.NetworkInspect(network, "{{ index .IPAM.Config 1 \"Subnet\"}}")
	if err != nil {
		return "", fmt.Errorf("failed to get subnet prefix: %w", err)
	}

	ip, sub, err := net.ParseCIDR(subnet)
	if err != nil {
		return "", fmt.Errorf("failed to get subnet prefix: %w", err)
	}

	size, _ := sub.Mask.Size()
	if size != 64 {
		return "", fmt.Errorf("ipv6 subnet of only size 64 is supported")
	}

	return strings.TrimSuffix(ip.String(), "::"), nil
}
