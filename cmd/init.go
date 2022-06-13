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
	"github.com/spf13/cobra"
	"github.com/utkarsh-pro/kindli/cmd/network"
	"github.com/utkarsh-pro/kindli/cmd/preq"
	"github.com/utkarsh-pro/kindli/cmd/vm"
	"github.com/utkarsh-pro/kindli/pkg/utils"
)

var (
	skipPreqInstall = false
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize kindli",
	Long: `init will initialize kindli

Initialization process will perform the following operations:
1. Install all the prerequisites (can be skipped, see the flags)
2. Start kindli default VM
3. Create a default kind cluster in the VM
4. Setup e2e networking to the KinD network`,
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Install preqs
		if !skipPreqInstall {
			utils.ExitIfNotNil(preq.RunInstall())
		}

		name, err := cmd.Flags().GetString("vm-name")
		utils.ExitIfNotNil(err)

		// 2. Start the VM
		utils.ExitIfNotNil(vm.RunStart(name))
		utils.ExitIfNotNil(vm.RunRestart(name))

		cname, err := cmd.Flags().GetString("cluster-name")
		utils.ExitIfNotNil(err)

		// 3. Create default cluster
		utils.ExitIfNotNil(RunCreate(cname, name))

		// 4. Setup e2e networking
		utils.ExitIfNotNil(network.RunSetup(name))
	},
}

func init() {
	// Support flags in the VM start command
	InitCmd.Flags().AddFlagSet(vm.StartCmd.LocalFlags())

	InitCmd.Flags().BoolVar(&skipPreqInstall, "skip-preq-install", false, "if set to true, prerequisite install will be skipped")
}
