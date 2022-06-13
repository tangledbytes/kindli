package vm

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/utkarsh-pro/kindli/pkg/config"
	"github.com/utkarsh-pro/kindli/pkg/sh"
)

func vmFilePath(vmName string) string {
	return filepath.Join(config.Dir(), fmt.Sprintf("%s.yaml", vmName))
}

//go:embed vm.template
var vmTemplate string

// Start takes VM config overrides and creates a VM using lima
func Start(overrides map[string]interface{}, skipIfExists bool, vmName string) error {
	exists, err := exists(vmName)
	if err != nil {
		return err
	}

	isRunning, err := Running(vmName)
	if err != nil {
		return err
	}

	if isRunning {
		if skipIfExists {
			return nil
		}

		return errors.New("kindli VM already exists")
	}

	if exists || overrides == nil {
		return sh.Run("limactl start --tty=false " + vmName)
	}

	if err := createLimaVMConfig(overrides, vmName); err != nil {
		return err
	}

	return sh.Run("limactl start --tty=false " + vmFilePath(vmName))
}

// Stop stops the currently running VM
func Stop(vmName string) error {
	isRunning, err := Running(vmName)
	if err != nil {
		return err
	}

	if !isRunning {
		return errors.New("kindli VM is not in running state")
	}

	return sh.Run("limactl stop " + vmName)
}

// Delete deltes the stopped VM
func Delete(vmName string) error {
	exist, err := exists(vmName)
	if err != nil {
		return err
	}

	if !exist {
		return errors.New("kindli VM does not exists")
	}

	if err := sh.Run("limactl delete " + vmName); err != nil {
		return err
	}

	return os.Remove(vmFilePath(vmName))
}

func Restart(vmName string) error {
	if err := Stop(vmName); err != nil {
		return err
	}

	return Start(nil, true, vmName)
}

func Status(vmName string) (string, error) {
	resp, err := sh.RunIO(fmt.Sprintf("limactl ls | awk '/NAME/ || /%s/ {print $0}'", vmName))
	if err != nil {
		return "", fmt.Errorf("failed to get status of VM: %s", err)
	}

	return string(resp), nil
}

func Shell(vmName string) error {
	return sh.Run("limactl shell " + vmName)
}

func createLimaVMConfig(overrides map[string]interface{}, vmName string) error {
	u, err := user.Current()
	if err != nil {
		return fmt.Errorf("failed to find username: %s", err)
	}
	overrides["user"] = u.Username

	file, err := os.Create(vmFilePath(vmName))
	if err != nil {
		return err
	}
	defer file.Close()

	parsed, err := template.New("kindli").Parse(vmTemplate)
	if err != nil {
		return err
	}

	return parsed.Execute(file, overrides)
}

func Running(vmName string) (bool, error) {
	out, err := exec.Command("limactl", "ls", "--format={{ .Name }}={{ .Status }}").CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("failed to get VM status: %s", err)
	}

	vms := strings.Split(string(out), "\n")
	for _, vm := range vms {
		prefix := fmt.Sprintf("%s=", vmName)

		if strings.HasPrefix(vm, prefix) {
			status := strings.TrimPrefix(vm, prefix)

			if status != "Stopped" {
				return true, nil
			}

			return false, nil
		}
	}

	return false, nil
}

func exists(vmName string) (bool, error) {
	out, err := exec.Command("limactl", "ls", "--format={{ .Name }}={{ .Status }}").CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("failed to get VM status: %s", err)
	}

	vms := strings.Split(string(out), "\n")
	for _, vm := range vms {
		prefix := fmt.Sprintf("%s=", vmName)

		if strings.HasPrefix(vm, prefix) {
			return true, nil
		}
	}

	return false, nil
}
