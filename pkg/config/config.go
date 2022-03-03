package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
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

func CleanupDir() error {
	if err := os.RemoveAll(configDir); err != nil {
		return fmt.Errorf("failed to cleanup data directory: %s", err)
	}

	return nil
}

// Dir returns path to the config directory
func Dir() string {
	return configDir
}

func Logger() {
	llevel := os.Getenv("LOG_LEVEL")
	switch strings.ToLower(llevel) {
	case "trace":
		logrus.SetLevel(logrus.TraceLevel)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	case "panic":
		logrus.SetLevel(logrus.PanicLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func CleanEnv() {
	os.Unsetenv("DOCKER_HOST")
	os.Unsetenv("DOCKER_CONTEXT")
	os.Unsetenv("DOCKER_MACHINE_NAME")
}
