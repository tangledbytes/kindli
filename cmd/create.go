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
	cfg                string
	name               string
	skipDockerRegistry bool
	skipQuayRegistry   bool
	skipGCRRegistry    bool
)

// CreateCmd represents create command
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new kind cluster",
	Run: func(cmd *cobra.Command, args []string) {
		utils.ExitIfNotNil(RunCreate())
	},
}

func init() {
	CreateCmd.Flags().StringVarP(&cfg, "config", "c", "", "kind configuration")
	CreateCmd.Flags().StringVar(&name, "name", "", "kind cluster name")
	CreateCmd.Flags().BoolVar(&skipDockerRegistry, "skip-registry-docker", false, "skip installing docker registry")
	CreateCmd.Flags().BoolVar(&skipGCRRegistry, "skip-registry-gcr", false, "skip installing GCR registry")
	CreateCmd.Flags().BoolVar(&skipQuayRegistry, "skip-registry-quay", false, "skip installing Quay registry")
}

func RunCreate() error {
	// Create docker context if it doesn't already exists
	ctxExists, err := docker.ExistsContext("kindli")
	if err != nil {
		return err
	}

	if !ctxExists {
		err := docker.CreateContext("kindli", fmt.Sprintf("host=unix://%s", filepath.Join(config.Dir(), "docker.sock")))
		if err != nil {
			return err
		}
	}

	// Switch default docker context to newly created docker context
	err = docker.UseContext("kindli")
	if err != nil {
		return err
	}

	// Create the kind cluster
	err = kind.Create(cfg, kind.CreateConfig{
		DockerRegistry: !skipDockerRegistry,
		QuayRegistry:   !skipQuayRegistry,
		GCRRegistry:    !skipGCRRegistry,
		Name:           name,
	})
	if err != nil {
		return err
	}

	return nil
}
