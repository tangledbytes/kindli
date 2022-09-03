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

package network

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/utkarsh-pro/kindli/pkg/networking"
	"github.com/utkarsh-pro/kindli/pkg/utils"
)

var SetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "setup e2e networking with cluster",
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("vm-name")
		utils.ExitIfNotNil(err)

		utils.ExitIfNotNil(RunSetup(name))
	},
}

func RunSetup(name string) error {
	if !warnUser() {
		return nil
	}

	return networking.Setup(name)
}

func warnUser() bool {
	var input string

	fmt.Print("⚠️  Warning: Managing up routes requires privilege escalation and would require root password. Do you want to continue? [y/n]: ")
	fmt.Scanln(&input)

	return strings.ToLower(input) == "y"
}
