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

package darwin

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
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

func TestParseKernProcargs2(t *testing.T) {
	testCases := []struct {
		data    []byte
		process process
		err     error
	}{
		{data: nil, err: io.EOF},
		{data: []byte{}, err: io.EOF},
		{data: []byte{0, 0, 0, 0}, process: process{env: map[string]string{}}},
		{data: []byte{5, 0, 0, 0}, process: process{env: map[string]string{}}},
		{
			data: buildKernProcargs2Data(3, "./example", []string{"/Users/test/example", "--one", "--two"}, []string{"TZ=UTC"}),
			process: process{
				exe:  "./example",
				args: []string{"/Users/test/example", "--one", "--two"},
				env: map[string]string{
					"TZ": "UTC",
				},
			},
		},
	}

	for i, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var p process
			err := parseKernProcargs2(tc.data, &p)
			if tc.err != nil {
				assert.ErrorIs(t, err, tc.err)
			} else {
				assert.EqualValues(t, tc.process, p)
			}
		})
	}
}

func FuzzParseKernProcargs2(f *testing.F) {
	f.Add([]byte(nil))
	f.Add([]byte{0, 0, 0, 0})
	f.Add([]byte{10, 0, 0, 0})
	f.Add(buildKernProcargs2Data(2, "./foo", []string{"/Users/john/foo", "-c"}, []string{"TZ=UTC"}))
	f.Add(buildKernProcargs2Data(100, "./foo", []string{"/Users/john/foo", "-c"}, []string{"TZ=UTC"}))

	f.Fuzz(func(t *testing.T, b []byte) {
		p := &process{}
		_ = parseKernProcargs2(b, p)
	})
}

// buildKernProcargs2Data builds a response that is similar to what
// sysctl kern.procargs2 returns.
func buildKernProcargs2Data(argc int32, exe string, args, envs []string) []byte {
	// argc
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, uint32(argc))

	// exe
	data = append(data, []byte(exe)...)
	data = append(data, nullTerminator...)

	// args
	for _, arg := range args {
		data = append(data, []byte(arg)...)
		data = append(data, nullTerminator...)
	}

	// env
	for _, env := range envs {
		data = append(data, []byte(env)...)
		data = append(data, nullTerminator...)
	}

	// The returned buffer from the real kern.procargs2 contains more data than
	// what go-sysinfo parses. This is a rough simulation of that extra data.
	data = append(data, bytes.Repeat(nullTerminator, 100)...)
	data = append(data, []byte("ptr_munge=")...)
	data = append(data, bytes.Repeat(nullTerminator, 18)...)
	data = append(data, []byte("main_stack==")...)
	data = append(data, bytes.Repeat(nullTerminator, 43)...)
	data = append(data, []byte("executable_file=0x1a01000010,0x36713a1")...)
	data = append(data, []byte("dyld_file=0x1a01000010,0xfffffff0008839c")...)
	data = append(data, []byte("executable_cdhash=5ca6024f9cdaa3a9fe515bfad77e1acf0f6b15b6")...)
	data = append(data, []byte("executable_boothash=a4a5613c07091ef0a221ee75a924341406eab85e")...)
	data = append(data, []byte("arm64e_abi=os")...)
	data = append(data, []byte("th_port=")...)
	data = append(data, bytes.Repeat(nullTerminator, 11)...)

	return data
}
