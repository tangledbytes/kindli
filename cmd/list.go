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
	"github.com/utkarsh-pro/kindli/pkg/kind"
	"github.com/utkarsh-pro/kindli/pkg/utils"
)

// ListCmd represents list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "list commands lists all of the KinD clusters",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("vm-name")
		all, _ := cmd.Flags().GetBool("all")

		if all {
			utils.ExitIfNotNil(RunList(""))
			return
		}

		utils.ExitIfNotNil(RunList(name))
	},
}

func RunList(vmName string) error {
	return kind.List(vmName)
}

func init() {
	ListCmd.Flags().BoolP("all", "A", false, "Set to list clusters of all the Kindli VMs")
}
