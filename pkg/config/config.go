package config

import (
	"os"
	"path/filepath"

	"github.com/utkarsh-pro/kindli/pkg/utils"
)

var (
	configDir = ""
)

func init() {
	SetupDir()
}

// SetupDir sets up the config dir required for kindli to function properly
//
// Note: If SetupDir encounters an error then the program will exit with an error
func SetupDir() {
	home, err := os.UserHomeDir()
	utils.ExitIfNotNil(err)

	configDir = filepath.Join(home, ".kindli")

	utils.ExitIfNotNil(os.MkdirAll(configDir, 0777))
}

// Dir returns path to the config directory
func Dir() string {
	return configDir
}
