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
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/ktateish/go-aoj"
	"github.com/spf13/cobra"
)

// testcaseCmd represents the testcase command
var testcaseCmd = &cobra.Command{
	Use:   "testcase",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: testcaseMain,
}

var (
	lengthFlag *bool
	input      *int
	output     *int
)

func testcaseMain(cmd *cobra.Command, args []string) {
	aoj.Debug(*DEBUG)

	if len(args) < 1 || args[0] == "" {
		cmd.Usage()
		os.Exit(1)
	}

	prob := args[0]

	t, err := aoj.GetTestcase(prob)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get testcases: %v\n", err)
		os.Exit(1)
	}

	if *input == -1 && *output == -1 {
		*lengthFlag = true
	}
	if *input >= 0 && *output >= 0 {
		cmd.Usage()
		os.Exit(1)
	}

	if *lengthFlag == true {
		fmt.Printf("%d\n", t.Length())
		return
	}

	cid := *input
	getReader := t.CaseInput

	if *output >= 0 {
		cid = *output
		getReader = t.CaseOutput
	}

	r, err := getReader(cid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get test case %d of %s\n", cid, prob)
		os.Exit(1)
	}
	defer r.Close()

	bo := bufio.NewWriter(os.Stdout)

	io.Copy(bo, r)
	bo.Flush()
}

func init() {
	RootCmd.AddCommand(testcaseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testcaseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	lengthFlag = testcaseCmd.Flags().BoolP("length", "l", false, "Show the number of test cases.")
	input = testcaseCmd.Flags().IntP("input", "i", -1, "Show input of the specified test case.")
	output = testcaseCmd.Flags().IntP("output", "o", -1, "Show output of the specified test case.")
}
