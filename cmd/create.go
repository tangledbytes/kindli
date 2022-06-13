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
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/utkarsh-pro/kindli/pkg/config"
	"github.com/utkarsh-pro/kindli/pkg/docker"
	"github.com/utkarsh-pro/kindli/pkg/kind"
	"github.com/utkarsh-pro/kindli/pkg/utils"
)

var (
	cfg string
)

// CreateCmd represents create command
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new kind cluster",
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("cluster-name")
		utils.ExitIfNotNil(err)

		cname, err := cmd.Flags().GetString("cluster-name")
		utils.ExitIfNotNil(err)

		utils.ExitIfNotNil(RunCreate(cname, name))
	},
}

func init() {
	CreateCmd.Flags().StringVarP(&cfg, "config", "c", "", "kind configuration")
}

func RunCreate(name string, vmName string) error {
	// Create docker context if it doesn't already exists
	ctxExists, err := docker.ExistsContext(vmName)
	if err != nil {
		return err
	}

	if !ctxExists {
		err := docker.CreateContext(vmName, fmt.Sprintf("host=unix://%s", filepath.Join(config.Dir(), "docker.sock")))
		if err != nil {
			return err
		}
	}

	// Switch default docker context to newly created docker context
	err = docker.UseContext(vmName)
	if err != nil {
		return err
	}

	// Create the kind cluster
	err = kind.Create(cfg, kind.CreateConfig{
		Name: name,
	})
	if err != nil {
		return err
	}

	return nil
}
