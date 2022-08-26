package metallb

import (
	_ "embed"
	"fmt"
	"os"
	"text/template"

	"github.com/utkarsh-pro/kindli/pkg/models"
	"github.com/utkarsh-pro/kindli/pkg/networking"
	"github.com/utkarsh-pro/kindli/pkg/sh"
)

var (
	//go:embed metallb.template
	metalLBTemplate string
)

// Install install metallb in the given cluster
func Install(clusterName string) error {
	if err := sh.Run("kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.12.1/manifests/namespace.yaml"); err != nil {
		return fmt.Errorf("failed to install metallb: %w", err)
	}

	if err := sh.Run("kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.12.1/manifests/metallb.yaml"); err != nil {
		return fmt.Errorf("failed to install metallb: %w", err)
	}

	if err := Configure(clusterName); err != nil {
		return fmt.Errorf("failed to configure metallb: %w", err)
	}

	return nil
}

// Configure configures metallb in the given cluster
func Configure(clusterName string) error {
	c := models.NewCluster(clusterName, "", "")
	if err := c.GetByName(); err != nil {
		return fmt.Errorf("failed to find cluster with name \"%s\": %w", clusterName, err)
	}

	intID := int(c.ID)
	if intID >= 99 {
		return fmt.Errorf("cannot configure more than 99 instances")
	}

	ipv4SubnetPrefix, ipv6SubnetPrefix, err := getSubnetPrefix("kind")
	if err != nil {
		return fmt.Errorf("failed to configure metallb: %w", err)
	}

	cfg := createNetworkConfig(ipv4SubnetPrefix, ipv6SubnetPrefix, intID)

	cfgPath, err := createConfig(cfg)
	if err != nil {
		return fmt.Errorf("failed to generate metallb config: %w", err)
	}

	if err := sh.Run(fmt.Sprintf("kubectl apply -f %s", cfgPath)); err != nil {
		return fmt.Errorf("failed to apply config to kubernetes: %w", err)
	}

	return nil
}

func getSubnetPrefix(network string) (string, string, error) {
	ipv4, err := networking.GetIPv4SubnetPrefix(network)
	if err != nil {
		return "", "", fmt.Errorf("failed to get ipv4 subnet prefix: %w", err)
	}
	ipv6, err := networking.GetIPv6SubnetPrefix(network)
	if err != nil {
		return "", "", fmt.Errorf("failed to get ipv6 subnet prefix: %w", err)
	}

	return ipv4, ipv6, nil
}

func createNetworkConfig(ipv4SubnetPrefix, ipv6SubnetPrefix string, instanceID int) map[string]interface{} {
	return map[string]interface{}{
		"ipv4Range": createIPv4NetworkConfig(ipv4SubnetPrefix, instanceID),
		"ipv6Range": createIPv6NetworkConfig(ipv6SubnetPrefix, instanceID),
	}
}

func createIPv4NetworkConfig(subnetPrefix string, instanceID int) string {
	return fmt.Sprintf("%s.%x.0/24", subnetPrefix, instanceID+1)
}

func createIPv6NetworkConfig(subnetPrefix string, instanceID int) string {
	return fmt.Sprintf("%s:%x::", subnetPrefix, instanceID+1)
}

func createConfig(cfg map[string]interface{}) (string, error) {
	parsed, err := template.New("metallb.template").Parse(metalLBTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to create metallb config: %s", err)
	}

	file, err := os.CreateTemp("", "*.yaml")
	if err != nil {
		return "", fmt.Errorf("failed to create metallb config: %s", err)
	}
	defer file.Close()

	if err := parsed.Execute(file, cfg); err != nil {
		return "", fmt.Errorf("failed to create metallb config: %s", err)
	}

	return file.Name(), nil
}
