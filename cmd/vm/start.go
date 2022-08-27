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

var (
	cpu    int
	mem    string
	disk   string
	arch   string
	mounts []string
)

// StartCmd represents the start command
var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a new VM for Kindli",
	Long: `Start a new VM for Kindli.

NOTE: VM will be created using lima`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if arch != "" && arch != "x86_64" && arch != "aarch64" {
			return fmt.Errorf("invalid --arch value, can be only \"x86_64\" or \"aarch64\"")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("vm-name")
		utils.ExitIfNotNil(err)
		utils.ExitIfNotNil(RunStart(name))
	},
}

func init() {
	StartCmd.Flags().IntVar(&cpu, "cpu", 4, "specify number of cpu assigned to VM")
	StartCmd.Flags().StringVar(&mem, "mem", "16GiB", "specify memory to be assigned to VM")
	StartCmd.Flags().StringVar(&disk, "disk", "100GiB", "specify disk space assigned to the VM")
	StartCmd.Flags().StringVar(&arch, "arch", "", "VM architecture")
	StartCmd.Flags().StringSliceVar(&mounts, "mount", nil, "specify mounts in form of <PATH>:rw to make the mount available for read/write or in form of <PATH>:ro to make the mount available only for reading")
}

func RunStart(name string) error {
	return vm.Start(createOverrides(), true, name)
}

func createOverrides() map[string]interface{} {
	overrides := map[string]interface{}{
		"CPU":    cpu,
		"Memory": mem,
		"Disk":   disk,
		"Arch":   arch,
	}

	parsedMounts, err := parseMounts(mounts)
	utils.ExitIfNotNil(err)

	overrides["Mounts"] = parsedMounts

	return overrides
}

func parseMounts(mounts []string) ([]map[string]interface{}, error) {
	mapped := []map[string]interface{}{}

	for _, mount := range mounts {
		splitted := strings.Split(mount, ":")
		if len(splitted) != 2 {
			return nil, fmt.Errorf("failed to parse mount: %s", mount)
		}

		writable := false

		if splitted[1] == "rw" {
			writable = true
		}

		mapped = append(mapped, map[string]interface{}{
			"location": splitted[0],
			"writable": writable,
		})
	}

	return mapped, nil
}
