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
	"go/build"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const localPkgs = "github.com/elastic/go-sysinfo"

var defaultPaths = []string{"."}

func main() {
	log.SetFlags(0)
	flag.Parse()

	paths := defaultPaths
	if len(flag.Args()) > 0 {
		paths = flag.Args()
	}

	goGet := exec.Command("go", "get", "-u", "golang.org/x/tools/cmd/goimports")
	goGet.Env = os.Environ()
	goGet.Env = append(goGet.Env, "GO111MODULE=off")
	out, err := goGet.Output()
	if err != nil {
		log.Fatalf("failed to %v: %v", strings.Join(goGet.Args, " "), err)
	}

	goimports := exec.Command(filepath.Join(build.Default.GOPATH, "bin", "goimports"),
		append([]string{"-l", "-local", localPkgs}, paths...)...)
	out, err = goimports.Output()
	if err != nil {
		log.Fatalf("failed to %v: %v", strings.Join(goimports.Args, " "), err)
	}
	if len(out) > 0 {
		fmt.Fprintln(os.Stderr, "Run goimports on the code.")
		fmt.Printf(string(out))
		os.Exit(1)
	}
}
