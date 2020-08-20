/*

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
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zwachtel11/peg/pkg/peg"
)

type deployOptions struct {
	manifest   string
	kubeconfig string
	username   string
	password   string
}

var deployOpts = &deployOptions{}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploys kubernetes deployment manifest from the registry",
	Long:  "deploys kubernetes deployment manifest from the registry",
	Example: "	peg deploy --manifest=cr.io/test-image:latest",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDeploy()
	},
}

func init() {
	deployCmd.Flags().StringVar(&deployOpts.manifest, "manifest", "", "manifest name")
	deployCmd.MarkFlagRequired("manifest")
	deployCmd.Flags().StringVar(&deployOpts.kubeconfig, "kubeconfig", "", "kubeconfig")

	RootCmd.AddCommand(deployCmd)
}

func runDeploy() error {

	resolver := newResolver(deployOpts.username, deployOpts.password)
	ctx := context.Background()

	err := peg.Deploy(ctx, resolver, deployOpts.manifest, deployOpts.kubeconfig)
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	return nil
}
