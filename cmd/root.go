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
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/utkarsh-pro/kindli/cmd/image"
	"github.com/utkarsh-pro/kindli/cmd/network"
	"github.com/utkarsh-pro/kindli/cmd/preq"
	"github.com/utkarsh-pro/kindli/cmd/vm"
	"github.com/utkarsh-pro/kindli/pkg/kind"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "kindli",
	Short: "Kindli lets users create upto 100 kind clusters in a Linux based virtual machine",
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.AddCommand(
		preq.PreqCmd,
		vm.VMCmd,
		network.NetworkCmd,
		image.ImageCmd,
		CreateCmd,
		DeleteCmd,
		InitCmd,
		PruneCmd,
		DockerEnvCmd,
		ListCmd,
	)

	RootCmd.PersistentFlags().String("vm-name", "kindli", "Name of the VM")
	RootCmd.RegisterFlagCompletionFunc(
		"vm-name",
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			list, err := vm.RunList()
			if err != nil {
				return []string{}, cobra.ShellCompDirectiveNoFileComp
			}

			final := []string{}
			for _, v := range list {
				if strings.HasPrefix(v, toComplete) {
					final = append(final, v)
				}
			}

			return final, cobra.ShellCompDirectiveNoFileComp
		},
	)

	RootCmd.PersistentFlags().String("cluster-name", "kindli", "Name of the cluster")
	RootCmd.RegisterFlagCompletionFunc(
		"cluster-name",
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			vm, _ := cmd.Flags().GetString("vm-name")
			all, _ := cmd.Flags().GetBool("all")

			if all {
				vm = ""
			}

			list, err := kind.PureList(vm)
			if err != nil {
				return []string{}, cobra.ShellCompDirectiveNoFileComp
			}

			final := []string{}
			for _, v := range list {
				if strings.HasPrefix(v, toComplete) {
					final = append(final, v)
				}
			}

			return final, cobra.ShellCompDirectiveNoFileComp
		},
	)
}
