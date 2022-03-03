/*
Copyright © 2022 Utkarsh Srivastava <utkarsh@sagacious.dev>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package preq

import (
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	ppreq "github.com/utkarsh-pro/kindli/pkg/preq"
)

// InstallCmd represents the install command
var InstallCmd = &cobra.Command{
	Use:     "install",
	Short:   "Install missing prequisites",
	Aliases: []string{"i"},
	Run: func(cmd *cobra.Command, args []string) {
		RunInstall()
	},
}

func RunInstall() error {
	if err := install(); err != nil {
		log.Error("❌ Failed to satisify prequisites")
		return err
	}

	log.Info("🚀 All prequisites satisified")

	return nil
}

func install() error {
	missings := ppreq.Missing()

	for _, missing := range missings {
		log.Infof("⚠️ %s missing: Attempting to install...", missing)

		switch missing {
		case "brew":
			if err := ppreq.InstallBrew(); err != nil {
				return err
			}
		case "git":
			if err := ppreq.InstallGit(); err != nil {
				return err
			}
		case "make":
			if err := ppreq.InstallMake(); err != nil {
				return err
			}
		case "automake":
			if err := ppreq.InstallAutoMake(); err != nil {
				return err
			}
		case "autoconf":
			if err := ppreq.InstallAutoConf(); err != nil {
				return err
			}
		case "limactl":
			if err := ppreq.InstallLima(); err != nil {
				return err
			}
		case "vde_switch":
			if err := ppreq.InstallSwitch(); err != nil {
				return err
			}
		case "vde_vmnet":
			if err := ppreq.InstallVMNet(); err != nil {
				return err
			}
		}
		log.Infof("✅ Installed %s", missing)
	}

	return nil
}
