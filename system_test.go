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

package sysinfo

import (
	"encoding/json"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/elastic/go-sysinfo/types"
)

type ProcessFeatures struct {
	ProcessInfo    bool
	Environment    bool
	FileDescriptor bool
	Seccomp        bool
	Capabilities   bool
}

var expectedProcessFeatures = map[string]*ProcessFeatures{
	"darwin": &ProcessFeatures{
		ProcessInfo:    true,
		Environment:    true,
		FileDescriptor: false,
	},
	"linux": &ProcessFeatures{
		ProcessInfo:    true,
		Environment:    true,
		FileDescriptor: true,
		Seccomp:        true,
		Capabilities:   true,
	},
	"windows": &ProcessFeatures{
		ProcessInfo: true,
	},
}

func TestProcessFeaturesMatrix(t *testing.T) {
	const GOOS = runtime.GOOS
	var features ProcessFeatures

	process, err := Self()
	if err == types.ErrNotImplemented {
		assert.Nil(t, expectedProcessFeatures[GOOS], "unexpected ErrNotImplemented for %v", GOOS)
		return
	} else if err != nil {
		t.Fatal(err)
	}
	features.ProcessInfo = true

	_, features.Environment = process.(types.Environment)
	_, features.FileDescriptor = process.(types.FileDescriptor)
	_, features.Seccomp = process.(types.Seccomp)
	_, features.Capabilities = process.(types.Capabilities)

	assert.Equal(t, expectedProcessFeatures[GOOS], &features)
	logAsJSON(t, map[string]interface{}{
		"features": features,
	})
}

func TestSelf(t *testing.T) {
	process, err := Self()
	if err == types.ErrNotImplemented {
		t.Skip("process provider not implemented on", runtime.GOOS)
	} else if err != nil {
		t.Fatal(err)
	}

	if runtime.GOOS == "linux" {
		// Do some dummy work to spend user CPU time.
		var v int
		for i := 0; i < 999999999; i++ {
			v += i * i
		}
	}

	output := map[string]interface{}{}
	info, err := process.Info()
	if err != nil {
		t.Fatal(err)
	}
	output["process.info"] = info
	assert.EqualValues(t, os.Getpid(), info.PID)
	assert.EqualValues(t, os.Getppid(), info.PPID)
	assert.Equal(t, os.Args, info.Args)

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	assert.EqualValues(t, wd, info.CWD)

	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, exe, info.Exe)

	if err != nil {
		t.Fatal(err)
	}
	assert.WithinDuration(t, info.StartTime, time.Now(), 10*time.Second)

	if v, ok := process.(types.Environment); ok {
		expectedEnv := map[string]string{}
		for _, keyValue := range os.Environ() {
			parts := strings.SplitN(keyValue, "=", 2)
			if len(parts) != 2 {
				t.Fatal("failed to parse os.Environ()")
			}
			expectedEnv[parts[0]] = parts[1]
		}
		actualEnv, err := v.Environment()
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, expectedEnv, actualEnv)
		output["process.env"] = actualEnv
	}

	memInfo, err := process.Memory()
	require.NoError(t, err)
	if runtime.GOOS != "windows" {
		// Virtual memory may be reported as
		// zero on some versions of Windows.
		assert.NotZero(t, memInfo.Virtual)
	}
	assert.NotZero(t, memInfo.Resident)
	output["process.mem"] = memInfo

	for {
		cpuTimes, err := process.CPUTime()
		require.NoError(t, err)
		if cpuTimes.Total() != 0 {
			output["process.cpu"] = cpuTimes
			break
		}
		// Spin until CPU times are non-zero.
		// Some operating systems have a very
		// low resolution on process CPU
		// measurement.
	}

	if v, ok := process.(types.FileDescriptor); ok {
		count, err := v.FileDescriptorCount()
		if assert.NoError(t, err) {
			t.Log("file descriptor count:", count)
		}
		fds, err := v.FileDescriptors()
		if assert.NoError(t, err) {
			output["process.fd"] = fds
		}
	}

	if v, ok := process.(types.Seccomp); ok {
		seccompInfo, err := v.Seccomp()
		if assert.NoError(t, err) {
			assert.NotZero(t, seccompInfo)
			output["process.seccomp"] = seccompInfo
		}
	}

	if v, ok := process.(types.Capabilities); ok {
		capInfo, err := v.Capabilities()
		if assert.NoError(t, err) {
			assert.NotZero(t, capInfo)
			output["process.capabilities"] = capInfo
		}
	}

	logAsJSON(t, output)
}

func TestHost(t *testing.T) {
	host, err := Host()
	if err == types.ErrNotImplemented {
		t.Skip("host provider not implemented on", runtime.GOOS)
	} else if err != nil {
		t.Fatal(err)
	}

	info := host.Info()
	assert.NotZero(t, info)

	memory, err := host.Memory()
	if err != nil {
		t.Fatal(err)
	}

	cpu, err := host.CPUTime()
	if err != nil {
		t.Fatal(err)
	}

	logAsJSON(t, map[string]interface{}{
		"host.info":   info,
		"host.memory": memory,
		"host.cpu":    cpu,
	})
}

func logAsJSON(t testing.TB, v interface{}) {
	if !testing.Verbose() {
		return
	}
	t.Helper()
	j, _ := json.MarshalIndent(v, "", "  ")
	t.Log(string(j))
}
