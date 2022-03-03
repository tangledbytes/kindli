package preq

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/utkarsh-pro/kindli/pkg/sh"
)

var (
	corePreqs = []string{
		"brew",
		"git",
		"make",
		"automake",
		"autoconf",
		"limactl",
	}

	preqs = []string{
		"vde_switch",
		"vde_vmnet",
	}

	preqpath = "/opt/vde/bin"
)

func Missing() []string {
	missing := []string{}

	for _, corePreq := range corePreqs {

		if _, err := exec.LookPath(corePreq); err != nil {
			missing = append(missing, corePreq)
		} else if corePreq == "limactl" {
			if _, err := os.Stat("/etc/sudoers.d/lima"); errors.Is(err, os.ErrNotExist) {
				missing = append(missing, corePreq)
			}
		}
	}

	for _, preq := range preqs {
		if _, err := os.Stat(fmt.Sprintf("%s/%s", preqpath, preq)); errors.Is(err, os.ErrNotExist) {
			missing = append(missing, preq)
		}
	}

	return missing
}

func InstallBrew() error {
	return sh.Run("/bin/bash -c '$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)'")
}

func InstallGit() error {
	return sh.Run("brew install git")
}

func InstallMake() error {
	return sh.Run("brew install make")
}

func InstallAutoMake() error {
	return sh.Run("brew install automake")
}

func InstallAutoConf() error {
	return sh.Run("brew install autoconf")
}

func InstallLima() error {
	return sh.RunMany([]string{
		"brew install lima",
		"limactl sudoers | sudo tee /etc/sudoers.d/lima",
	})
}

func InstallSwitch() error {
	dir := fmt.Sprint(time.Now().UnixNano())

	return sh.RunMany([]string{
		"cd $TMPDIR",
		"mkdir " + dir,
		"cd " + dir,
		"git clone https://github.com/virtualsquare/vde-2.git",
		"cd vde-2",
		"autoreconf -fis",
		"./configure --prefix=/opt/vde",
		"make",
		"sudo make install",
	})
}

func InstallVMNet() error {
	dir := fmt.Sprint(time.Now().UnixNano())

	return sh.RunMany([]string{
		"cd $TMPDIR",
		"mkdir " + dir,
		"cd " + dir,
		"git clone https://github.com/lima-vm/vde_vmnet",
		"cd vde_vmnet",
		"make PREFIX=/opt/vde",
		"sudo make PREFIX=/opt/vde install.bin",
	})
}
