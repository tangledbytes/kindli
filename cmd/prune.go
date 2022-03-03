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
	"github.com/utkarsh-pro/kindli/cmd/vm"
	"github.com/utkarsh-pro/kindli/pkg/config"
	"github.com/utkarsh-pro/kindli/pkg/utils"
)

var PruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "prune kindli will cleanup kindli",
	Long: `prune will prune kindli

Prune process will perform the following operations:
1. Stop the running VM
2. Delete the VM
3. Cleanup the ~/.kindli directory`,
	Run: func(cmd *cobra.Command, args []string) {
		// Stop the running VM
		utils.ExitIfNotNil(vm.RunStop())

		// Delete the VM
		utils.ExitIfNotNil(vm.RunDelete())

		// Cleanup dirs
		utils.ExitIfNotNil(config.CleanupDir())
	},
}
