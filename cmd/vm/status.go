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
package vm

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/utkarsh-pro/kindli/pkg/utils"
	"github.com/utkarsh-pro/kindli/pkg/vm"
)

// StatusCmd represents the status command
var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Status of Kindli VM",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("vm-name")
		all, _ := cmd.Flags().GetBool("all")

		if all {
			utils.ExitIfNotNil(RunStatus(""))
			return
		}

		utils.ExitIfNotNil(RunStatus(name))
	},
}

func init() {
	StatusCmd.Flags().BoolP("all", "A", false, "Show status of all VMs")
}

func RunStatus(name string) error {
	status, err := vm.Status(name)
	if err != nil {
		return err
	}

	fmt.Println(status)

	return nil
}
