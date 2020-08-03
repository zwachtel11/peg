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
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/zwachtel11/peg/pkg/peg"
)

type pullOptions struct {
	manifest string
	username string
	password string
	outfile  string
}

var pullOpts = &pullOptions{}

var pullCmd = &cobra.Command{
	Use:     "pull",
	Short:   "pulls kubernetes deployment manifest from the registry",
	Long:    "pulls kubernetes deployment manifest from the registry",
	Example: "peg pull --manifest=cr.io/test-image:latest",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runPull()
	},
}

func init() {
	pullCmd.Flags().StringVarP(&pullOpts.username, "username", "u", "", "registry username")
	pullCmd.Flags().StringVarP(&pullOpts.password, "password", "p", "", "registry password or identity token")
	pullCmd.Flags().StringVar(&pullOpts.manifest, "manifest", "", "manifest name")
	pullCmd.Flags().StringVar(&pullOpts.outfile, "outfile", "", "location to write the manifest")
	pullCmd.MarkFlagRequired("manifest")
	RootCmd.AddCommand(pullCmd)
}

func runPull() error {
	resolver := newResolver(pullOpts.username, pullOpts.password)
	ctx := context.Background()
	bytes, err := peg.Pull(ctx, resolver, pullOpts.manifest)
	if err != nil {
		fmt.Printf("err: %s", err)
	}

	if pullOpts.outfile != "" {
		err := ioutil.WriteFile(pullOpts.outfile, bytes, 0644)
		if err != nil {
			fmt.Printf("err: %s", err)
		}
	} else {
		fmt.Printf(string(bytes))
	}

	return nil
}
