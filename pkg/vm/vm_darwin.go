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

	"github.com/sirupsen/logrus"
	"github.com/utkarsh-pro/kindli/pkg/config"
	"github.com/utkarsh-pro/kindli/pkg/models"
	"github.com/utkarsh-pro/kindli/pkg/sh"
)

func vmFilePath(vmName string) string {
	return filepath.Join(config.Dir(), fmt.Sprintf("%s.yaml", vmName))
}

func defaultDockerPort() int {
	return 2375
}

//go:embed vm.template
var vmTemplate string

// Start takes VM config overrides and creates a VM using lima
func Start(overrides map[string]interface{}, skipIfExists bool, vmName string) error {
	logrus.Infoln("Starting VM:", vmName)

	exists, err := exists(vmName)
	if err != nil {
		return fmt.Errorf("failed to check if VM exists: %w", err)
	}

	isRunning, err := Running(vmName)
	if err != nil {
		return fmt.Errorf("failed to check if VM is running: %w", err)
	}

	if isRunning {
		if skipIfExists {
			return nil
		}

		return fmt.Errorf("vm is already running")
	}

	if exists || overrides == nil {
		err := sh.Run("limactl start --tty=false " + vmName)
		if err != nil {
			return fmt.Errorf("failed to start VM: %w", err)
		}

		return nil
	}

	// Create a new VM
	port, err := models.GetMaxVMDockerPort()
	if err != nil {
		return fmt.Errorf("failed to get max docker port for the VM: %w", err)
	}
	if port == 0 {
		port = defaultDockerPort()
	} else {
		port++
	}

	vm := models.NewVM(vmName, vmFilePath(vmName), port)

	if err := createLimaVMConfig(overrides, vm); err != nil {
		return fmt.Errorf("failed to create lima VM config: %w", err)
	}
	if err := vm.Save(); err != nil {
		return fmt.Errorf("failed to save vm instance: %w", err)
	}

	if err := sh.Run("limactl start --tty=false " + vm.LimaConfigPath); err != nil {
		return fmt.Errorf("failed to start VM: %w", err)
	}

	return nil
}

// Stop stops the currently running VM
func Stop(vmName string) error {
	isRunning, err := Running(vmName)
	if err != nil {
		return err
	}

	if !isRunning {
		return errors.New("VM is not in running state")
	}

	return sh.Run("limactl stop " + vmName)
}

// Delete deletes the stopped VM
func Delete(vmName string) error {
	exist, err := exists(vmName)
	if err != nil {
		return err
	}

	if !exist {
		return errors.New("VM does not exists")
	}

	vm := models.NewVM(vmName, "", 0)
	if err := vm.GetByName(); err != nil {
		return fmt.Errorf("failed to get VM by name: %w", err)
	}

	if err := sh.Run("limactl delete " + vm.Name); err != nil {
		return err
	}

	if err := os.Remove(vm.LimaConfigPath); err != nil {
		return err
	}

	return vm.Delete()
}

// Restart restarts the VM
func Restart(vmName string) error {
	if err := Stop(vmName); err != nil {
		return err
	}

	return Start(nil, true, vmName)
}

// Status shows the status of the VM
func Status(vmName string) (string, error) {
	resp, err := sh.RunIO(fmt.Sprintf("limactl ls | awk '/NAME/ || /%s/ {print $0}'", vmName))
	if err != nil {
		return "", fmt.Errorf("failed to get status of VM: %s", err)
	}

	return string(resp), nil
}

// Shell opens the shell of the VM
func Shell(vmName string, args ...string) error {
	if len(args) == 0 {
		return sh.Run("limactl shell " + vmName)
	}

	return sh.Run("limactl shell " + vmName + " -- " + strings.Join(args, " "))
}

// Running checks if the VM is running
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

// List returns a list of all the VMs
func List() ([]string, error) {
	vms, err := models.ListVM()
	if err != nil {
		return nil, err
	}

	var vmNames []string
	for _, vm := range vms {
		vmNames = append(vmNames, vm.Name)
	}

	return vmNames, nil
}

func createLimaVMConfig(overrides map[string]interface{}, vm *models.VM) error {
	logrus.Debug("Creating lima VM config at:", vm.LimaConfigPath)
	u, err := user.Current()
	if err != nil {
		return fmt.Errorf("failed to find username: %w", err)
	}
	overrides["user"] = u.Username
	overrides["vmName"] = vm.Name
	overrides["dockerPort"] = vm.DockerPort

	file, err := os.Create(vm.LimaConfigPath)
	if err != nil {
		return fmt.Errorf("failed to create VM config file: %w", err)
	}
	defer file.Close()

	parsed, err := template.New(vm.Name).Parse(vmTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse VM config template: %w", err)
	}

	if err := parsed.Execute(file, overrides); err != nil {
		return fmt.Errorf("failed to execute VM config template: %w", err)
	}

	return nil
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
