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

type pushOptions struct {
	file     string
	manifest string
}

var pushOpts = &pushOptions{}

var pushCmd = &cobra.Command{
	Use:     "push",
	Short:   "pushes kubernetes deployment manifest to registry",
	Long:    "pushes kubernetes deployment manifest to registry",
	Example: "peg push --file=~/path/to/config ",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runPush()
	},
}

func init() {
	pushCmd.Flags().StringVar(&pushOpts.file, "file", "", "Path to input configuration")
	pushCmd.MarkFlagRequired("file")
	pushCmd.Flags().StringVar(&pushOpts.manifest, "manifest", "", "Manifest Name")
	RootCmd.AddCommand(pushCmd)
}

func runPush() error {
	resolver := newResolver(loginOpts.username, loginOpts.password)
	ctx := context.Background()
	err := peg.Push(ctx, resolver, pushOpts.manifest, pushOpts.file)
	if err != nil {
		fmt.Printf("err: %s", err)
	}

	return nil
}
