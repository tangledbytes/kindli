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

package image

import (
	"github.com/spf13/cobra"
	"github.com/utkarsh-pro/kindli/pkg/image"
	"github.com/utkarsh-pro/kindli/pkg/utils"
)

var LoadCmd = &cobra.Command{
	Use:     "load",
	Short:   "Load OCI images in the VM",
	Example: `kindli image <image-name> <cluster-name>`,
	Args:    cobra.MatchAll(cobra.MinimumNArgs(1), cobra.MaximumNArgs(2)),
	Run: func(cmd *cobra.Command, args []string) {
		image := args[0]
		cluster := ""

		if len(args) == 2 {
			cluster = args[1]
		}

		utils.ExitIfNotNil(RunLoad(image, cluster))
	},
}

func init() {
	ImageCmd.AddCommand()
}

func RunLoad(imageName, clusterName string) error {
	image.Load(imageName, clusterName)
	return nil
}
