/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"io/ioutil"
	"kubeconf/pkg/merger"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// mergeCmd represents the merge command
var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merges two kubeconfig file into one",
	Long:  `Kubeconf is a cli tool which merges kubeconfig files.`,
	Run: func(cmd *cobra.Command, args []string) {

		kubeconfigPath := cmd.Flag("kubeconfig").Value.String()

		currentKubeConfig := merger.NewKubeConfig(kubeconfigPath)

		newConfigPath := cmd.Flag("new-config").Value.String()
		dry, _ := cmd.LocalFlags().GetBool("dry")
		showChanges, _ := cmd.LocalFlags().GetBool("show-changes")
		outputFile, _ := cmd.LocalFlags().GetString("output")

		if newConfigPath != "" {
			newConfig := merger.NewKubeConfig(newConfigPath)
			currentKubeConfig.MergeNewConfig(*newConfig)

			if currentKubeConfig.IsChanged {
				if showChanges {
					if len(currentKubeConfig.ToAddClusters) > 0 {
						fmt.Println(len(currentKubeConfig.ToAddClusters))
						fmt.Println("Clusters will be added to kubeconfig:")
						toShow, _ := yaml.Marshal(currentKubeConfig.ToAddClusters)
						fmt.Println(string(toShow))
					}

					if len(currentKubeConfig.ToAddContexts) > 0 {
						fmt.Println("Contexts will be added to kubeconfig:")
						toShow, _ := yaml.Marshal(currentKubeConfig.ToAddContexts)
						fmt.Println(string(toShow))
					}

					if len(currentKubeConfig.ToAddUsers) > 0 {
						fmt.Println("Users will be added to kubeconfig:")
						toShow, _ := yaml.Marshal(currentKubeConfig.ToAddUsers)
						fmt.Println(string(toShow))
					}
				}

				fmt.Fprintf(os.Stdout, "Changes will be applied... Do you accept? yes-no (y/n)\n")

				var input string
				fmt.Scanln(&input)

				if input == "yes" || input == "y" {
					configBytes, err := yaml.Marshal(currentKubeConfig)
					if dry {
						if err != nil {
							fmt.Println(err.Error())
						} else {
							fmt.Println(string(configBytes))
						}
					} else {
						writePath := kubeconfigPath
						if outputFile != "" {
							writePath = outputFile
						}
						err := ioutil.WriteFile(writePath, configBytes, 0644)
						if err != nil {
							fmt.Println(err.Error())
						}
					}
				} else {
					fmt.Println("Changes discarded.")
				}
			}
		}
	},
}

func init() {
	userHomeDir, _ := os.UserHomeDir()

	mergeCmd.Flags().String("new-config", "", "New config file path.")
	mergeCmd.Flags().String("kubeconfig", userHomeDir+"/.kube/config", "The kubeconfig file path.")
	mergeCmd.Flags().Bool("dry", false, "Shows generated output and does not persists config changes.")
	mergeCmd.Flags().Bool("show-changes", false, "Shows resources which will be added to kubeconfig.")
	mergeCmd.Flags().String("output", "", "Output file path.")
}

func Execute() {
	if err := mergeCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
