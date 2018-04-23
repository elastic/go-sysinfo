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
	"flag"
	"fmt"
	"os"
	"os/exec"
)

const localPkgs = "github.com/elastic/go-sysinfo"

var defaultPaths = []string{"."}

func main() {
	flag.Parse()

	paths := defaultPaths
	if len(flag.Args()) > 0 {
		paths = flag.Args()
	}

	out, err := exec.Command("go", "get", "golang.org/x/tools/cmd/goimports").Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error executing goimports: %v", string(err.(*exec.ExitError).Stderr))
		os.Exit(1)
	}

	args := append([]string{"-l", "-local", localPkgs}, paths...)
	out, err = exec.Command("goimports", args...).Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error", err.(*exec.ExitError))
		os.Exit(1)
	}
	if len(out) > 0 {
		fmt.Fprintln(os.Stderr, "Run goimports on the code.")
		fmt.Printf(string(out))
		os.Exit(1)
	}
}
