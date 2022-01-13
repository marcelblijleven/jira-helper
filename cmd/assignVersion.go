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
	"jira-helper/pkg"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

// assignVersionCmd represents the assignVersion command
var assignVersionCmd = &cobra.Command{
	Use:   "assignVersion",
	Short: "Assigns a version to all provided issues in the release body",
	Long: `Assigns a version to all provided issues. The issue numbers are retrieved from
the provided release body.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		httpClient := http.DefaultClient
		httpClient.Timeout = time.Second * 15
		client, err := pkg.NewJiraClient(host, user, token, httpClient)

		if err != nil {
			return err
		}

		err = pkg.AssignVersions(body, version, client)

		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(assignVersionCmd)
	assignVersionCmd.Flags().StringVarP(&user, "user", "u", "", "user used for authenticating against the Jira API")
	assignVersionCmd.Flags().StringVarP(&host, "host", "s", "", "host of the Jira API")
	assignVersionCmd.Flags().StringVarP(&project, "project", "p", "", "Abbreviation of the Jira project, e.g. GGWM")
	assignVersionCmd.Flags().StringVarP(&token, "token", "t", "", "Token used to authenticate against the Jira API")
	assignVersionCmd.Flags().StringVarP(&version, "version", "v", "", "Version name")
	assignVersionCmd.Flags().StringVarP(&body, "releaseBody", "b", "", "The body of the Github release")
}
