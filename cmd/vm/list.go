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
	"strings"

	"github.com/spf13/cobra"
	"github.com/utkarsh-pro/kindli/pkg/utils"
	"github.com/utkarsh-pro/kindli/pkg/vm"
)

// ListCmd represents the shell command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "Prints the list of VMs",
	Run: func(cmd *cobra.Command, args []string) {
		vms, err := RunList()
		utils.ExitIfNotNil(err)

		fmt.Println(strings.Join(vms, "\n"))
	},
}

func RunList() ([]string, error) {
	return vm.List()
}
