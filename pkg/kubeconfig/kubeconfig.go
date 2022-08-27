package kubeconfig

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/utkarsh-pro/kindli/pkg/utils"
)

// SetCurrentContext takes a context name and sets it as the current context
func SetCurrentContext(name string) error {
	return DoXOnKubeconfig(func(mp map[string]interface{}) (bool, error) {
		val, ok := utils.MapGet(mp, "current-context")
		if !ok {
			utils.MapSet(mp, name, "current-context")
			return true, nil
		}

		ctx, ok := val.(string)
		if !ok {
			utils.MapSet(mp, name, "current-context")
			return true, nil
		}

		if ctx == name {
			return false, nil
		}

		utils.MapSet(mp, name, "current-context")
		return true, nil
	})
}

// DoXOnKubeconfig executes a function on the kubeconfig file
func DoXOnKubeconfig(x func(map[string]interface{}) (bool, error)) error {
	byt, err := getKubeconfig()
	if err != nil {
		return fmt.Errorf("failed to perform action on kubeconfig: %w", err)
	}

	mp, err := parseKubeConfig(byt)
	if err != nil {
		return fmt.Errorf("failed to perform action on kubeconfig: %w", err)
	}

	if proceed, err := x(mp); err != nil {
		return fmt.Errorf("failed to perform action on kubeconfig: %w", err)
	} else if !proceed {
		return nil
	}

	byt, err = serializeKubeconfig(mp)
	if err != nil {
		return fmt.Errorf("failed to perform action on kubeconfig: %w", err)
	}

	if err := writeKubeconfig(byt); err != nil {
		return fmt.Errorf("failed to perform action on kubeconfig: %w", err)
	}

	return nil
}

// KubeconfigPath returns the path to the kubeconfig file
func KubeconfigPath() string {
	env := os.Getenv("KUBECONFIG")
	if env != "" {
		return env
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "~/.kube/config"
	}

	return filepath.Join(home, ".kube", "config")
}

func parseKubeConfig(byt []byte) (map[string]interface{}, error) {
	mp, err := utils.MapFromYAML(byt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse kubeconfig: %w", err)
	}

	return mp, nil
}

func getKubeconfig() ([]byte, error) {
	byt, err := os.ReadFile(KubeconfigPath())
	if err != nil {
		return nil, fmt.Errorf("failed to read kubeconfig: %w", err)
	}

	return byt, nil
}

func serializeKubeconfig(mp map[string]interface{}) ([]byte, error) {
	return utils.MapToYAML(mp)
}

func writeKubeconfig(byt []byte) error {
	path := KubeconfigPath()
	if err := os.WriteFile(path, byt, 0644); err != nil {
		return fmt.Errorf("failed to write kubeconfig: %w", err)
	}

	return nil
}
