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
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "jira-helper",
	Short: "Helper tool to create and assign version in Jira from the CLI and CI/CD",
	Long: `Helper tool to interact with Jira from CI/CD scripts. Its main purpose is to create and assign version
based on Github releases to Jira tickets.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize()
	rootCmd.PersistentFlags().StringVarP(&user, userFlagName, userShorthand, "", userUsage)
	rootCmd.PersistentFlags().StringVarP(&host, hostFlagName, hostShorthand, "", hostUsage)
	rootCmd.PersistentFlags().StringVarP(&project, projectFlagName, projectShorthand, "", projectUsage)
	rootCmd.PersistentFlags().StringVarP(&released, releasedFlagName, releasedShorthand, "", releasedUsage)
	rootCmd.PersistentFlags().StringVarP(&token, tokenFlagName, tokenShorthand, "", tokenUsage)
	rootCmd.PersistentFlags().StringVarP(&version, versionFlagName, versionShorthand, "", versionUsage)

	cobra.CheckErr(rootCmd.MarkPersistentFlagRequired(userFlagName))
	cobra.CheckErr(rootCmd.MarkPersistentFlagRequired(hostFlagName))
	cobra.CheckErr(rootCmd.MarkPersistentFlagRequired(projectFlagName))
	cobra.CheckErr(rootCmd.MarkPersistentFlagRequired(releasedFlagName))
	cobra.CheckErr(rootCmd.MarkPersistentFlagRequired(tokenFlagName))
	cobra.CheckErr(rootCmd.MarkPersistentFlagRequired(versionFlagName))
}
