package metallb

import (
	_ "embed"
	"fmt"
	"os"
	"text/template"

	"github.com/utkarsh-pro/kindli/pkg/docker"
	"github.com/utkarsh-pro/kindli/pkg/models"
	"github.com/utkarsh-pro/kindli/pkg/sh"
)

var (
	//go:embed metallb.template
	metalLBTemplate string
)

// Install install metallb in the given cluster
func Install(clusterName string) error {
	if err := sh.Run("kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.13.4/manifests/namespace.yaml"); err != nil {
		return fmt.Errorf("failed to install metallb: %w", err)
	}

	if err := sh.Run("kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.13.4/manifests/metallb.yaml"); err != nil {
		return fmt.Errorf("failed to install metallb: %w", err)
	}

	if err := Configure(clusterName); err != nil {
		return fmt.Errorf("failed to configure metallb: %w", err)
	}

	return nil
}

// Configure configures metallb in the given cluster
func Configure(clusterName string) error {
	networkPrefix, err := docker.NetworkInspect("kind", "{{ join (slice (split (index .IPAM.Config 0 \"Subnet\") \".\") 0 2) \".\" }}")
	if err != nil {
		return fmt.Errorf("failed to configure metallb: %w", err)
	}

	c := models.NewCluster(clusterName, "", "")
	if err := c.GetByName(); err != nil {
		return fmt.Errorf("failed to find cluster with name \"%s\": %w", clusterName, err)
	}

	intID := int(c.ID)
	if intID >= 99 {
		return fmt.Errorf("cannot configure more than 99 instances")
	}

	cfg := createNetworkConfig(networkPrefix, intID)

	cfgPath, err := createConfig(cfg)
	if err != nil {
		return fmt.Errorf("failed to generate metallb config: %w", err)
	}

	if err := sh.Run(fmt.Sprintf("kubectl apply -f %s", cfgPath)); err != nil {
		return fmt.Errorf("failed to apply config to kubernetes: %w", err)
	}

	return nil
}

func createNetworkConfig(subnetPrefix string, instanceID int) map[string]interface{} {
	return map[string]interface{}{
		"addresses": fmt.Sprintf("%s.%d.1-%s.1%02d.254", subnetPrefix, instanceID+1, subnetPrefix, instanceID+1),
	}
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
