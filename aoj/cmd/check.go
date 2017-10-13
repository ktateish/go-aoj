// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/ktateish/go-aoj"
	"github.com/spf13/cobra"
)

var (
	caseID *int
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check an executable for the problem.",
	Long: `Check an executable for the problem.

problem_id: the problem id of AOJ courses. (e.g. ALDS1_1_A)
path:	    path to the executable implementing the solution.

Example:
    aoj check ALDS1_1_A 0 ./myprog
`,
	Run: checkMain,
}

func checkMain(cmd *cobra.Command, args []string) {
	aoj.Debug(*DEBUG)

	if len(args) < 2 {
		cmd.Usage()
		os.Exit(1)
	}

	prob := args[0]
	path := args[1]

	if strings.Index(path, "/") == -1 {
		dir, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot get cwd: %v\n", err)
			os.Exit(1)
		}
		path = fmt.Sprintf("%s/%s", dir, path)
	}

	t, err := aoj.GetTestcase(prob)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get testcases: %v\n", err)
		os.Exit(1)
	}

	if *caseID >= 0 {
		rs, err := t.CheckCase(*caseID, path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to check: %v\n", err)
			os.Exit(1)
		}
		if !rs {
			os.Exit(1)
		}
		return
	}
	for i := 0; i < t.Length(); i++ {
		rs, err := t.CheckCase(i, path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to check testcase %d: %v\n", i, err)
			os.Exit(1)
		}
		if !rs {
			fmt.Fprintf(os.Stderr, "testcase %d failed\n", i)
			os.Exit(1)
		}
	}
}

func init() {
	RootCmd.AddCommand(checkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	caseID = checkCmd.Flags().IntP("case", "c", -1, "Case id to check")
}
