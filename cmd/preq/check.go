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
package preq

import (
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/utkarsh-pro/kindli/pkg/preq"
)

// CheckCmd represents the check command
var CheckCmd = &cobra.Command{
	Use:     "check",
	Short:   "Checks if prerequisites are avaialable on the system or not - Doesn't install missing prequisites though",
	Aliases: []string{"i"},
	Run: func(cmd *cobra.Command, args []string) {
		if !checkPreqs() {
			log.Error("❌ Failed to satisify prequisites")
			return
		}

		log.Info("🚀 All prequisites satisified")
	},
}

func checkPreqs() bool {
	missings := preq.Missing()

	for _, missing := range missings {
		log.Infof("⚠️ %s missing", missing)
	}

	return len(missings) == 0
}
