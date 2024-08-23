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

//go:build integration

package linux

import (
	"bytes"
	"go/build"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"golang.org/x/sys/execabs"
)

const (
	wantHostname = "hostname"
	wantDomain   = "some.domain"
	wantFQDN     = wantHostname + "." + wantDomain
)

func TestHost_FQDN(t *testing.T) {
	if _, err := execabs.LookPath("docker"); err != nil {
		t.Skipf("Skipping because docker was not found: %v", err)
	}

	tcs := []struct {
		name       string
		Hostname   string
		Domainname string
		Cmd        []string
	}{
		{
			name:       "TestHost_FQDN_set_hostname+domainname",
			Hostname:   wantHostname,
			Domainname: wantDomain,
			Cmd: []string{
				"go", "test", "-v",
				"-tags", "integration,docker",
				"-run", "^TestHost_FQDN_set$",
				"./providers/linux",
			},
		},
		{
			name:     "TestHost_FQDN_set_hostname_only",
			Hostname: wantFQDN,
			Cmd: []string{
				"go", "test", "-v",
				"-tags", "integration,docker",
				"-run", "^TestHost_FQDN_set$",
				"./providers/linux",
			},
		},
		{
			name: "TestHost_FQDN_not_set",
			Cmd: []string{
				"go", "test", "-v", "-count", "1",
				"-tags", "integration,docker",
				"-run", "^TestHost_FQDN_not_set$",
				"./providers/linux",
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			dockerRun(t, tc.Hostname, tc.Domainname, tc.Cmd)
		})
	}
}

// dockerRun executes the given command inside the golang Docker container.
// It will set the container's hostname and domain according to the given
// values.
//
// The container's stdout and stderr are passed through. The test will fail
// if docker run returns a non-zero exit code.
func dockerRun(t *testing.T, hostname, domain string, command []string) {
	t.Helper()

	// Determine the repository root.
	_, filename, _, _ := runtime.Caller(0)
	repoRoot, err := filepath.Abs(filepath.Join(filepath.Dir(filename), "../.."))
	if err != nil {
		t.Fatal(err)
	}

	// Use the same version of Go inside the container.
	goVersion := strings.TrimPrefix(runtime.Version(), "go")

	args := []string{
		"run",
		"--rm",
		"-v", build.Default.GOPATH + ":/go", // Mount GOPATH for caching.
		"-v", repoRoot + ":/go-sysinfo",
		"-w=/go-sysinfo",
	}
	if hostname != "" {
		args = append(args, "--hostname="+hostname)
	}
	if domain != "" {
		args = append(args, "--domainname="+domain)
	}
	args = append(args, "golang:"+goVersion)
	args = append(args, command...)

	buf := new(bytes.Buffer)
	cmd := execabs.Command("docker", args...)
	cmd.Stdout = buf
	cmd.Stderr = buf

	t.Logf("Running docker container using %q", args)
	defer t.Log("Exiting container")

	err = cmd.Run()
	t.Logf("Container output:\n%s", buf.String())
	if err != nil {
		t.Fatal(err)
	}
}
