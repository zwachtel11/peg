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
	"flag"
	"os"

	"github.com/spf13/cobra"
)

type Options struct {
}

var opts = &Options{}

var RootCmd = &cobra.Command{
	Use:          "peg",
	SilenceUsage: true,
	Short:        "\n",
	Long:         "",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		// TODO: print error stack if log v>0
		// TODO: print cmd help if validation error
		os.Exit(1)
	}
}

func init() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}
