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
package fips

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/utkarsh-pro/kindli/pkg/utils"
	"github.com/utkarsh-pro/kindli/pkg/vm"
)

// CheckCmd represents the check command
var CheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check if FIPS is enabled in the VM",
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("vm-name")
		utils.ExitIfNotNil(err)
		utils.ExitIfNotNil(RunCheck(name))
	},
}

func RunCheck(name string) error {
	status, err := vm.FipsCheck(name)
	if err != nil {
		return err
	}

	if status {
		fmt.Println("FIPS is enabled in the VM")
		return nil
	}

	fmt.Println("FIPS is not enabled in the VM")
	return nil
}