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

//go:build (amd64 && cgo) || (arm64 && cgo)
// +build amd64,cgo arm64,cgo

package darwin

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/elastic/go-sysinfo/internal/registry"
)

var (
	_ registry.HostProvider    = darwinSystem{}
	_ registry.ProcessProvider = darwinSystem{}
)

func TestKernProcInfo(t *testing.T) {
	var p process
	if err := kern_procargs(os.Getpid(), &p); err != nil {
		t.Fatal(err)
	}

	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, exe, p.exe)
	assert.Equal(t, os.Args, p.args)
}

const (
	noValueEnvVar    = "_GO_SYSINFO_NO_VALUE"
	emptyValueEnvVar = "_GO_SYSINFO_EMPTY_VALUE"
	fooValueEnvVar   = "_GO_SYSINFO_FOO_VALUE"
)

func TestProcessEnvironment(t *testing.T) {
	cmd := exec.Command("go", "test", "-v", "-run", "^TestProcessEnvironmentInternal$")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env,
		// Activate the test case.
		"GO_SYSINFO_ENV_TESTING=1",
		// Set specific values that the test asserts.
		noValueEnvVar,
		emptyValueEnvVar+"=",
		fooValueEnvVar+"=FOO",
	)

	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "TestProcessEnvironmentInternal failed:\n"+string(out))
}

func TestProcessEnvironmentInternal(t *testing.T) {
	// This test case is executes in its own process space with a specific
	// environment set by TestProcessEnvironment.
	if os.Getenv("GO_SYSINFO_ENV_TESTING") != "1" {
		t.Skip()
	}

	var p process
	if err := kern_procargs(os.Getpid(), &p); err != nil {
		t.Fatal(err)
	}

	value, exists := p.env[noValueEnvVar]
	assert.True(t, exists, "Missing "+noValueEnvVar)
	assert.Equal(t, "", value)

	value, exists = p.env[emptyValueEnvVar]
	assert.True(t, exists, "Missing "+emptyValueEnvVar)
	assert.Equal(t, "", value)

	assert.Equal(t, "FOO", p.env[fooValueEnvVar])
}
