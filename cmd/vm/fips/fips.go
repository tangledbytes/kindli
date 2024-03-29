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
	"github.com/spf13/cobra"
)

// FipsCmd represents the fips command
var FipsCmd = &cobra.Command{
	Use:   "fips",
	Short: "Commands for managing FIPS in the VM",
}

func init() {
	FipsCmd.AddCommand(
		CheckCmd,
	)
}
