// Copyright 2018 The Terraformer Authors.
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

package cmd

import (
	"log"
	"os"

	rundeck_terraforming "github.com/GoogleCloudPlatform/terraformer/providers/rundeck"

	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	"github.com/spf13/cobra"
)

const (
	defaultRundeckAPIVersion = "38"
)

func newCmdRundeckImporter(options ImportOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rundeck",
		Short: "Import current state to Terraform configuration from Rundeck",
		Long:  "Import current state to Terraform configuration from Rundeck",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get configuration from environment variables
			url := os.Getenv("RUNDECK_URL")
			if len(url) == 0 {
				log.Fatal("RUNDECK_URL environment variable is required")
			}

			token := os.Getenv("RUNDECK_TOKEN")
			username := os.Getenv("RUNDECK_USERNAME")
			password := os.Getenv("RUNDECK_PASSWORD")

			// Require either token or username/password
			if token == "" && (username == "" || password == "") {
				log.Fatal("Either RUNDECK_TOKEN or both RUNDECK_USERNAME and RUNDECK_PASSWORD must be set")
			}

			apiVersion := os.Getenv("RUNDECK_API_VERSION")
			if apiVersion == "" {
				apiVersion = defaultRundeckAPIVersion
			}

			insecureSSL := os.Getenv("RUNDECK_INSECURE_SSL")
			if insecureSSL == "" {
				insecureSSL = "false"
			}

			provider := newRundeckProvider()
			log.Println(provider.GetName() + " importing")
			
			err := Import(provider, options, []string{
				url,
				token,
				username,
				password,
				apiVersion,
				insecureSSL,
			})
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.AddCommand(listCmd(newRundeckProvider()))
	baseProviderFlags(cmd.PersistentFlags(), &options, "jobs,projects", "job=id1:id2:id4")
	return cmd
}

func newRundeckProvider() terraformutils.ProviderGenerator {
	return &rundeck_terraforming.RundeckProvider{}
}