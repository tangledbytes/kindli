package kind

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/utkarsh-pro/kindli/pkg/config"
	"github.com/utkarsh-pro/kindli/pkg/registry"
	"github.com/utkarsh-pro/kindli/pkg/sh"
	"github.com/utkarsh-pro/kindli/pkg/store"
	"github.com/utkarsh-pro/kindli/pkg/utils"
	"gopkg.in/yaml.v2"
)

var (
	instancesDirName = "kind"
	instanceDirPath  = ""

	//go:embed kind.template
	kindTemplate string
)

type CreateConfig struct {
	DockerRegistry bool
	QuayRegistry   bool
	GCRRegistry    bool
}

func DefaultCreateConfig() CreateConfig {
	return CreateConfig{
		DockerRegistry: true,
		QuayRegistry:   true,
		GCRRegistry:    true,
	}
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
		return err
	}

	// Get the name of the user's kind config => Also a sanity test for the config
	name, err := getUserKindCfgName(userKindCfg)
	if err != nil {
		return err
	}

	// Check if the instance with same name exists or not
	if Exists(name) {
		return fmt.Errorf("instance with name \"%s\" already exists", name)
	}

	// Load kindli store
	lstore, ok := store.Get()
	if !ok {
		return fmt.Errorf("failed to get data from store")
	}

	instance := len(lstore.(map[string]interface{}))

	if instance >= 99 {
		return fmt.Errorf("cannot create more than 99 instances")
	}

	// Setup networking info
	createNetworking(instance, userKindCfg)

	// Setup registry info
	createRegistries(cfg, userKindCfg)

	// Custom Kind Config
	customConfig, err := createCustomKindConfig(userKindCfg)
	if err != nil {
		return fmt.Errorf("failed to create kind config with overrides: %s", err)
	}

	// Persist the altered user config
	newPath, err := persistAlteredConfig(name, customConfig)
	if err != nil {
		return fmt.Errorf("failed to persist kind config locally")
	}

	// Create kind cluster
	if err := sh.Run(fmt.Sprintf("kind create cluster --config %s", newPath)); err != nil {
		return fmt.Errorf("failed to create kind cluster: %s", err)
	}

	// Save the new instance in the store
	store.Set(
		map[string]interface{}{
			"path": newPath,
		},
		name,
	)

	return nil
}

func List() []string {
	resp := []string{}

	instances, ok := store.Get()
	if !ok {
		return resp
	}

	casted, ok := instances.(map[string]interface{})
	if !ok {
		return resp
	}

	for k := range casted {
		resp = append(resp, k)
	}

	return resp
}

func Delete(name string) error {
	if name == "" {
		name = "kindli"
	}

	_, ok := store.Get(name)
	if !ok {
		return fmt.Errorf("instance with name \"%s\" does not exists", name)
	}

	if err := sh.Run(fmt.Sprintf("kind delete cluster --name=%s", name)); err != nil {
		return fmt.Errorf("failed to delete kind instance: %s", err)
	}

	if err := os.Remove(filepath.Join(instanceDirPath, fmt.Sprintf("%s.yaml", name))); err != nil {
		return fmt.Errorf("failed to delete instance: %s", err)
	}

	store.DeleteTop(name)

	return nil
}

func Get(name string) (interface{}, bool) {
	return store.Get(name)
}

func Exists(name string) bool {
	_, exists := store.Get(name)

	return exists
}

func getUserKindCfgName(userKindCfg map[string]interface{}) (string, error) {
	name, ok := utils.MapGet(userKindCfg, "name")
	if !ok {
		return "kindli", nil
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

func createRegistries(cfg CreateConfig, userKindCfg map[string]interface{}) {
	registries := map[string]string{}

	for _, registry := range registry.Knowns() {
		switch registry.Name {
		case "dockerio-registry":
			if cfg.DockerRegistry {
				registries["docker"] = fmt.Sprintf("%s:%s", registry.Name, registry.Port)
			}
		case "quayio-registry":
			if cfg.DockerRegistry {
				registries["quay"] = fmt.Sprintf("%s:%s", registry.Name, registry.Port)
			}
		case "gcrio-registry":
			if cfg.DockerRegistry {
				registries["gcr"] = fmt.Sprintf("%s:%s", registry.Name, registry.Port)
			}
		}
	}

	utils.MapSet(userKindCfg, registries, "registries")
}

func loadUserKindConfig(path string) (map[string]interface{}, error) {
	mp := make(map[string]interface{})

	if path != "" {
		byt, err := ioutil.ReadFile(path)
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
