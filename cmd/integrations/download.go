// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package integrations

import (
	"github.com/srinandan/integrationcli/apiclient"
	"github.com/srinandan/integrationcli/client/integrations"

	"github.com/spf13/cobra"
)

//DownloadVerCmd to publish an integration flow version
var DownloadVerCmd = &cobra.Command{
	Use:   "download",
	Short: "Download an integration flow version",
	Long:  "Download an integration flow version",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if err = apiclient.SetRegion(region); err != nil {
			return err
		}
		if err = validate(); err != nil {
			return err
		}
		return apiclient.SetProjectID(project)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if version != "" {
			_, err = integrations.Download(name, version)
		} else if userLabel != "" {
			_, err = integrations.DownloadUserLabel(name, userLabel)
		} else if snapshot != "" {
			_, err = integrations.DownloadSnapshot(name, snapshot)
		}
		return

	},
}

func init() {
	DownloadVerCmd.Flags().StringVarP(&name, "name", "n",
		"", "Integration flow name")
	DownloadVerCmd.Flags().StringVarP(&version, "ver", "v",
		"", "Integration flow version")
	DownloadVerCmd.Flags().StringVarP(&userLabel, "user-label", "u",
		"", "Integration flow user label")
	DownloadVerCmd.Flags().StringVarP(&snapshot, "snapshot", "s",
		"", "Integration flow snapshot number")

	_ = DownloadVerCmd.MarkFlagRequired("name")
}
