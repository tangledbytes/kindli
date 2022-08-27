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
	"github.com/utkarsh-pro/kindli/pkg/sh"
	"github.com/utkarsh-pro/kindli/pkg/utils"
)

// EditCmd represents the edit command
var EditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit Kindli VM",
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("vm-name")
		utils.ExitIfNotNil(err)
		utils.ExitIfNotNil(RunEdit(name))
	},
}

func RunEdit(name string) error {
	return sh.Run(fmt.Sprintf("limactl edit %s", name))
}
