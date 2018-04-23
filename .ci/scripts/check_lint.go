// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var ignoreWarnings = []string{
	`don't use underscores in Go names`,
}

var ignoreWarningsRe = regexp.MustCompile(strings.Join(ignoreWarnings, "|"))

func main() {
	flag.Parse()

	golintArgs := flag.Args()
	if len(golintArgs) == 0 {
		out, err := exec.Command("go", "list", "./...").Output()
		if err != nil {
			fmt.Fprintln(os.Stderr, "error", err)
			os.Exit(1)
		}
		golintArgs = strings.Split(string(out), "\n")
	}

	out, err := exec.Command("go", "get", "github.com/golang/lint/golint").Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error", err)
		os.Exit(1)
	}

	golint := exec.Command("golint", golintArgs...)
	out, err = golint.Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error", err)
		os.Exit(1)
	}

	out, err = filterIgnores(out)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error", err)
		os.Exit(1)
	}

	if len(out) > 0 {
		fmt.Printf(string(out))
		os.Exit(1)
	}
}

func filterIgnores(out []byte) ([]byte, error) {
	var buffer bytes.Buffer
	s := bufio.NewScanner(bytes.NewReader(out))
	for s.Scan() {
		if !ignoreWarningsRe.Match(s.Bytes()) {
			buffer.Write(s.Bytes())
			buffer.WriteByte('\n')
		}
	}

	return buffer.Bytes(), s.Err()
}
