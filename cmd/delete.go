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
	"github.com/utkarsh-pro/kindli/pkg/docker"
	"github.com/utkarsh-pro/kindli/pkg/kind"
	"github.com/utkarsh-pro/kindli/pkg/utils"
)

// DeleteCmd represents create command
var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete given kind cluster",
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("vm-name")
		utils.ExitIfNotNil(err)

		cname, err := cmd.Flags().GetString("cluster-name")
		utils.ExitIfNotNil(err)

		utils.ExitIfNotNil(docker.UseContext(name))

		utils.ExitIfNotNil(kind.Delete(cname))
	},
}
