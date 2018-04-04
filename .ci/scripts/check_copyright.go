// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
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
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

// Directories to ignore.
var ignores = []string{
	`vendor/`,
	`^testing/`,
}

var ignoreRe *regexp.Regexp

// File extensions to check
var checkExts = map[string]bool{
	".c":  true,
	".go": true,
}

// Valid copyright headers, searched for in the top five lines in each file.
var copyrightRegexps = []string{
	`Licensed to Elasticsearch B.V.`,
	`Created by cgo -godefs - DO NOT EDIT`,
	`MACHINE GENERATED`,
}

var copyrightRe = regexp.MustCompile(strings.Join(copyrightRegexps, "|"))

func init() {
	ignorePattern := strings.Join(ignores, "|")

	if runtime.GOOS == "windows" {
		// Modify file separators for Windows.
		ignorePattern = strings.Replace(ignorePattern, "/", `\\`, -1)
	}

	ignoreRe = regexp.MustCompile(ignorePattern)
}

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		args = []string{"."}
	}

	for _, dir := range args {
		err := filepath.Walk(dir, checkCopyright)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

func checkCopyright(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return nil
	}
	if ignoreRe.MatchString(path) {
		return nil
	}
	if !checkExts[filepath.Ext(path)] {
		return nil
	}

	fd, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fd.Close()

	scanner := bufio.NewScanner(fd)
	for i := 0; scanner.Scan() && i < 5; i++ {
		if copyrightRe.MatchString(scanner.Text()) {
			return nil
		}
	}

	return fmt.Errorf("Missing copyright in %s?", path)
}
