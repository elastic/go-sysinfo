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

//go:build amd64 || arm64
// +build amd64 arm64

package darwin

import (
	"errors"
	"os"
	"os/exec"
	"syscall"
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

func TestProcesses(t *testing.T) {
	var s darwinSystem
	processes, err := s.Processes()
	if err != nil {
		t.Fatal(err)
	}

	var count int
	for _, proc := range processes {
		processInfo, err := proc.Info()
		switch {
		// Ignore processes that no longer exist or that cannot be accessed.
		case errors.Is(err, syscall.ESRCH),
			errors.Is(err, syscall.EPERM),
			errors.Is(err, syscall.EINVAL):
			continue
		case err != nil:
			t.Fatalf("failed to get process info for PID=%d: %v", proc.PID(), err)
		default:
			count++
		}

		if processInfo.PID == 0 {
			t.Fatalf("empty pid in %#v", processInfo)
		}

		if processInfo.Exe == "" {
			t.Fatalf("empty exec in %#v", processInfo)
		}

		u, err := proc.User()
		require.NoError(t, err)

		require.NotEmpty(t, u.UID)
		require.NotEmpty(t, u.EUID)
		require.NotEmpty(t, u.SUID)
		require.NotEmpty(t, u.GID)
		require.NotEmpty(t, u.EGID)
		require.NotEmpty(t, u.SGID)
	}

	assert.NotZero(t, count, "failed to get process info for any processes")
}
