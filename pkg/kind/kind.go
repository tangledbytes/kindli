package kind

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/sirupsen/logrus"
	"github.com/utkarsh-pro/kindli/pkg/config"
	"github.com/utkarsh-pro/kindli/pkg/kubeconfig"
	"github.com/utkarsh-pro/kindli/pkg/metallb"
	"github.com/utkarsh-pro/kindli/pkg/models"
	"github.com/utkarsh-pro/kindli/pkg/sh"
	"github.com/utkarsh-pro/kindli/pkg/utils"
	"gopkg.in/yaml.v2"
)

var (
	instancesDirName    = "kind"
	instanceDirPath     = ""
	defaultInstanceName = "kindli"

	//go:embed kind.template
	kindTemplate string
)

type CreateConfig struct {
	Name        string
	VMName      string
	SkipMetalLB bool
}

func init() {
	instanceDirPath = filepath.Join(config.Dir(), instancesDirName)
	utils.ExitIfNotNil(os.MkdirAll(instanceDirPath, 0777))
}

// Create takes path to a kind configuration file and creates
// a new kind instance in the VM based on the config file passed
func Create(cfgPath string, cfg CreateConfig) error {
	// Load user's kind config
	userKindCfg, err := loadUserKindConfig(cfgPath)
	if err != nil {
		return fmt.Errorf("failed to read user config: %s", err)
	}

	// Get the name of the user's kind config => Also a sanity test for the config
	name, err := getUserKindCfgName(userKindCfg, cfg.Name)
	if err != nil {
		return err
	}

	// Check if the instance with same name exists or not
	if Exists(name, cfg.VMName) {
		logrus.Warn("instance already exists: skipping cluster creation")
		logrus.Warn("skipped cluster creation - proceed with metallb configuration")
		if err := kubeconfig.SetCurrentContext(kindifyClusterName(name)); err != nil {
			return fmt.Errorf("failed to set kubeconfig context: %w", err)
		}

		if err := metallb.Install(name); err != nil {
			return fmt.Errorf("failed to create metallb config for the kind cluster: %w", err)
		}

		return nil
	}

	// Create kind cluster
	if err := createKindCluster(name, cfg.VMName, userKindCfg); err != nil {
		return fmt.Errorf("failed to create kind cluster: %s", err)
	}

	// Create metallb for the kind cluster
	if !cfg.SkipMetalLB {
		if err := metallb.Install(name); err != nil {
			return fmt.Errorf("failed to create metallb config for the kind cluster: %w", err)
		}
	}

	return nil
}

func Delete(name string) error {
	c := models.NewCluster(name, "", "")
	if err := c.GetByName(); err != nil {
		return fmt.Errorf("instance with name \"%s\" does not exists", name)
	}

	if err := sh.Run(fmt.Sprintf("kind delete cluster --name=%s", c.Name)); err != nil {
		return fmt.Errorf("failed to delete kind instance: %s", err)
	}

	if err := os.Remove(c.KindConfigPath); err != nil {
		return fmt.Errorf("failed to delete instance: %w", err)
	}

	if err := c.Delete(); err != nil {
		return fmt.Errorf("failed to delete cluster: %w", err)
	}

	return nil
}

func Exists(name, vmName string) bool {
	ok, err := models.NewCluster(name, "", vmName).Exists()
	if err != nil {
		logrus.Warn("failed to find instance in store: %w", err)
		return false
	}

	return ok
}

func List(vmName string) error {
	clusters, err := models.ListCluster()
	if err != nil {
		return fmt.Errorf("failed to list clusters: %w", err)
	}

	if len(clusters) == 0 {
		logrus.Warn("No clusters found - create a cluster with `kindli create`")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 4, 8, 4, ' ', 0)
	fmt.Fprintln(w, "NAME\tVMNAME")

	for _, c := range clusters {
		if c.VM == vmName || vmName == "" {
			fmt.Fprintf(w, "%s\t%s\n", c.Name, c.VM)
		}
	}

	return w.Flush()
}

func createKindCluster(name, vmName string, userKindCfg map[string]interface{}) error {
	// Save the new instance in the store
	cluster := models.NewCluster(name, "", vmName)
	if err := cluster.AssignID(); err != nil {
		return fmt.Errorf("failed to assign ID to cluster: %w", err)
	}

	// Setup networking info
	createNetworking(int(cluster.ID), userKindCfg)

	// Custom Kind Config
	customConfig, err := createCustomKindConfig(userKindCfg)
	if err != nil {
		return fmt.Errorf("failed to create kind config with overrides: %w", err)
	}

	// Persist the altered user config
	cluster.KindConfigPath, err = persistAlteredConfig(name, customConfig)
	if err != nil {
		return fmt.Errorf("failed to persist kind config locally: %w", err)
	}

	// Create kind cluster
	if err := sh.Run(fmt.Sprintf("kind create cluster --config %s", cluster.KindConfigPath)); err != nil {
		return fmt.Errorf("failed to create kind cluster: %w", err)
	}
	if err := cluster.Save(); err != nil {
		return fmt.Errorf("failed to save cluster: %w", err)
	}

	return nil
}

func getUserKindCfgName(userKindCfg map[string]interface{}, customName string) (string, error) {
	name, ok := utils.MapGet(userKindCfg, "name")
	if !ok {
		if customName == "" {
			customName = defaultInstanceName
		}

		userKindCfg["name"] = customName
		return customName, nil
	}

	nameStr, ok := name.(string)
	if !ok {
		return "", fmt.Errorf("invalid name")
	}

	return nameStr, nil
}

func persistAlteredConfig(name string, userKindCfg map[string]interface{}) (string, error) {
	path := filepath.Join(instanceDirPath, fmt.Sprintf("%s.yaml", name))
	file, err := os.Create(path)
	if err != nil {
		return path, fmt.Errorf("failed to create config file: %s", err)
	}

	return path, yaml.NewEncoder(file).Encode(userKindCfg)
}

func createNetworking(instance int, userKindCfg map[string]interface{}) {
	svcSubnet := fmt.Sprintf("10.%d.0.0/16", instance)
	podSubnet := fmt.Sprintf("10.1%02d.0.0/16", instance)

	utils.MapSet(userKindCfg, svcSubnet, "networking", "serviceSubnet")
	utils.MapSet(userKindCfg, podSubnet, "networking", "podSubnet")
}

func loadUserKindConfig(path string) (map[string]interface{}, error) {
	mp := make(map[string]interface{})

	if path != "" {
		byt, err := os.ReadFile(path)
		if err != nil {
			return mp, err
		}

		mp, err = utils.MapFromYAML(byt)
		if err != nil {
			return nil, err
		}
	}

	return mp, nil
}

func createCustomKindConfig(userKindCfg map[string]interface{}) (map[string]interface{}, error) {
	parsed, err := template.New("kind.template").Parse(kindTemplate)
	if err != nil {
		return nil, err
	}

	// Convert "nodes" field to YAML
	nodes, ok := userKindCfg["nodes"]
	if ok {
		nodeByt, err := yaml.Marshal(nodes)
		if err != nil {
			return nil, err
		}

		userKindCfg["nodes"] = string(nodeByt)
	}

	buf := &bytes.Buffer{}
	if err := parsed.Execute(buf, userKindCfg); err != nil {
		return nil, err
	}

	return utils.MapFromYAML(buf.Bytes())
}

func kindifyClusterName(name string) string {
	return "kind-" + name
}
