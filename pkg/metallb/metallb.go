package metallb

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/utkarsh-pro/kindli/pkg/config"
	"github.com/utkarsh-pro/kindli/pkg/models"
	"github.com/utkarsh-pro/kindli/pkg/networking"
	"github.com/utkarsh-pro/kindli/pkg/sh"
	"github.com/utkarsh-pro/kindli/pkg/utils"
)

var (
	instanceDirName = "metallb"
	instanceDirPath = ""

	//go:embed metallb.template
	metalLBTemplate string

	//go:embed metallb.l2ad.yaml
	metalLBAdv string
)

func init() {
	instanceDirPath = filepath.Join(config.Dir(), instanceDirName)
	utils.ExitIfNotNil(os.MkdirAll(instanceDirPath, 0777))
}

// Install install metallb in the given cluster
func Install(clusterName string) error {
	if err := sh.Run("kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.13.5/config/manifests/metallb-native.yaml"); err != nil {
		return fmt.Errorf("failed to install metallb: %w", err)
	}

	// if err := sh.Run("kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.12.1/manifests/metallb.yaml"); err != nil {
	// 	return fmt.Errorf("failed to install metallb: %w", err)
	// }

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

	cfgPath, err := createConfig(cfg, clusterName)
	if err != nil {
		return fmt.Errorf("failed to generate metallb config: %w", err)
	}

	if err := sh.Run("kubectl wait --for=condition=available --timeout=600s deployment -n metallb-system controller"); err != nil {
		return fmt.Errorf("failed to wait for metallb controller to be available: %w", err)
	}

	if err := sh.Run(fmt.Sprintf("kubectl apply -f %s", cfgPath)); err != nil {
		return fmt.Errorf("failed to apply IP Address Pool config to kubernetes: %w", err)
	}

	if err := sh.Run(fmt.Sprintf("cat <<EOF | kubectl apply -f -\n%s\nEOF", metalLBAdv)); err != nil {
		return fmt.Errorf("failed to apply L2 Advertisement config to kubernetes: %w", err)
	}

	return nil
}

func LoadConfigFromDisk(clusterName string) (map[string]interface{}, error) {
	path := filepath.Join(instanceDirPath, fmt.Sprintf("%s.yaml", clusterName))

	yaml, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load metallb config: %w", err)
	}

	cfgMap, err := utils.MapFromYAML(yaml)
	if err != nil {
		return nil, fmt.Errorf("failed to load metallb config: %w", err)
	}

	return cfgMap, nil
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

func createConfig(cfg map[string]interface{}, clusterName string) (string, error) {
	parsed, err := template.New("metallb.template").Parse(metalLBTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to create metallb config: %s", err)
	}

	path := filepath.Join(instanceDirPath, fmt.Sprintf("%s.yaml", clusterName))
	file, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("failed to create metallb config: %s", err)
	}
	defer file.Close()

	if err := parsed.Execute(file, cfg); err != nil {
		return "", fmt.Errorf("failed to create metallb config: %s", err)
	}

	return file.Name(), nil
}
