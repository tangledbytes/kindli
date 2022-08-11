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

package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/utkarsh-pro/kindli/cmd/vm"
	"github.com/utkarsh-pro/kindli/pkg/config"
	"github.com/utkarsh-pro/kindli/pkg/sh"
	"github.com/utkarsh-pro/kindli/pkg/utils"
)

var PruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "prune kindli will cleanup kindli",
	Long: `prune will prune kindli

Prune process will perform the following operations:
1. Stop all of the running VMs
2. Delete all of the VMs
3. Cleanup the ~/.kindli directory
4. Optionally clean up the lima cache`,
	Run: func(cmd *cobra.Command, args []string) {
		all, err := cmd.Flags().GetBool("all")
		utils.ExitIfNotNil(err)

		if !all {
			name, err := cmd.Flags().GetString("vm-name")
			utils.ExitIfNotNil(err)

			// Stop the running VM
			if err := vm.RunStop(name); err != nil {
				logrus.Warn("Failed to stop VM: ", err)
			}

			// Delete the VM
			if err := vm.RunDelete(name); err != nil {
				logrus.Error("Failed to delete VM: ", err)
			}

			return
		}

		vms, err := vm.RunList()
		utils.ExitIfNotNil(err)

		for _, name := range vms {
			// Stop the running VM
			if err := vm.RunStop(name); err != nil {
				logrus.Warn("Failed to stop VM: ", err)
			}

			// Delete the VM
			if err := vm.RunDelete(name); err != nil {
				logrus.Error("Failed to delete VM: ", err)
			}
		}

		// Cleanup dirs
		if err := config.CleanupDir(); err != nil {
			logrus.Error("Failed to cleanup ~/.kindli: ", err)
			logrus.Info("You can remove ~/.kindli manually")
		}

		// Cleanup lima config
		cleanLima, err := cmd.Flags().GetBool("clean-lima")
		utils.ExitIfNotNil(err)
		if cleanLima {
			if err := sh.Run("limactl prune"); err != nil {
				logrus.Error("Failed to remove lima cache: ", err)
			}
		}
	},
}

func init() {
	PruneCmd.Flags().BoolP("all", "a", true, "If true, prune will delete all of the VMs")
	PruneCmd.Flags().Bool("clean-lima", false, "If true, prune will clear lima cache")
}
