/*
Copyright © 2022 Marcel Blijleven

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
	"github.com/marcelblijleven/jira-helper/pkg"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

// createReleaseCmd represents the createRelease command
var createReleaseCmd = &cobra.Command{
	Use:   "createRelease",
	Short: "Create a fix version in Jira",
	Long: `Create a fix version in Jira for the project with the provided name.

The release state of the fix version will be set to "released" and the day will be set to 
today.`,
	Run: func(cmd *cobra.Command, args []string) {
		httpClient := http.DefaultClient
		httpClient.Timeout = time.Second * 15
		client, err := pkg.NewJiraClient(host, user, token, httpClient)
		cobra.CheckErr(err)
		cobra.CheckErr(client.CreateFixVersion(version, project, released))
	},
}

func init() {
	rootCmd.AddCommand(createReleaseCmd)
	createReleaseCmd.Aliases = []string{"createVersion"}
}
