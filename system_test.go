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
	"errors"
	"io/fs"
	"os"
	osUser "os/user"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/elastic/go-sysinfo/types"
)

type ProcessFeatures struct {
	ProcessInfo          bool
	Environment          bool
	OpenHandleEnumerator bool
	OpenHandleCounter    bool
	Seccomp              bool
	Capabilities         bool
}

var expectedProcessFeatures = map[string]*ProcessFeatures{
	"darwin": {
		ProcessInfo:          true,
		Environment:          true,
		OpenHandleEnumerator: false,
		OpenHandleCounter:    false,
	},
	"linux": {
		ProcessInfo:          true,
		Environment:          true,
		OpenHandleEnumerator: true,
		OpenHandleCounter:    true,
		Seccomp:              true,
		Capabilities:         true,
	},
	"windows": {
		ProcessInfo:          true,
		OpenHandleEnumerator: false,
		OpenHandleCounter:    true,
	},
	"aix": {
		ProcessInfo:          true,
		Environment:          true,
		OpenHandleEnumerator: false,
		OpenHandleCounter:    false,
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
	_, features.OpenHandleEnumerator = process.(types.OpenHandleEnumerator)
	_, features.OpenHandleCounter = process.(types.OpenHandleCounter)
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
	assert.EqualValues(t, os.Getpid(), process.PID())

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
	assert.Equal(t, os.Args, info.Args)

	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, exe, info.Exe)

	if user, err := process.User(); !errors.Is(err, types.ErrNotImplemented) {
		if err != nil {
			t.Fatal(err)
		}
		output["process.user"] = user

		currentUser, err := osUser.Current()
		if err != nil {
			t.Fatal(err)
		}
		assert.EqualValues(t, currentUser.Uid, user.UID)
		assert.EqualValues(t, currentUser.Gid, user.GID)

		if runtime.GOOS != "windows" {
			assert.EqualValues(t, strconv.Itoa(os.Geteuid()), user.EUID)
			assert.EqualValues(t, strconv.Itoa(os.Getegid()), user.EGID)
		}
	}

	if v, ok := process.(types.Environment); ok {
		actualEnv, err := v.Environment()
		if err != nil {
			t.Fatal(err)
		}
		output["process.env"] = actualEnv

		// Format the output to match format from os.Environ().
		keyEqualsValueList := make([]string, 0, len(actualEnv))
		for k, v := range actualEnv {
			keyEqualsValueList = append(keyEqualsValueList, k+"="+v)
		}
		sort.Strings(keyEqualsValueList)

		expectedEnv := os.Environ()
		sort.Strings(expectedEnv)

		assert.Equal(t, expectedEnv, keyEqualsValueList)
	}

	if memInfo, err := process.Memory(); !errors.Is(err, types.ErrNotImplemented) {
		require.NoError(t, err)
		if runtime.GOOS != "windows" {
			// Virtual memory may be reported as
			// zero on some versions of Windows.
			assert.NotZero(t, memInfo.Virtual)
		}
		assert.NotZero(t, memInfo.Resident)
		output["process.mem"] = memInfo
	}

	for {
		cpuTimes, err := process.CPUTime()
		if errors.Is(err, types.ErrNotImplemented) {
			break
		}

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

	if v, ok := process.(types.OpenHandleEnumerator); ok {
		fds, err := v.OpenHandles()
		if assert.NoError(t, err) {
			output["process.fd"] = fds
		}
	}

	if v, ok := process.(types.OpenHandleCounter); ok {
		count, err := v.OpenHandleCount()
		if assert.NoError(t, err) {
			t.Log("open handles count:", count)
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
	} else if err != nil && !strings.Contains(err.Error(), "FQDN") {
		t.Fatal(err)
	}

	info := host.Info()
	assert.NotZero(t, info)

	output := map[string]interface{}{}
	output["host.info"] = info

	if v, ok := host.(types.LoadAverage); ok {
		loadAvg, err := v.LoadAverage()
		if err != nil {
			t.Fatal(err)
		}
		output["host.loadavg"] = loadAvg
	}

	memory, err := host.Memory()
	if err != nil {
		t.Fatal(err)
	}
	output["host.memory"] = memory

	cpu, err := host.CPUTime()
	if errors.Is(err, types.ErrNotImplemented) {
		t.Log("CPU times not implemented")
		return
	}

	if err != nil {
		t.Fatal(err)
	}
	output["host.cpu"] = cpu

	logAsJSON(t, output)
}

func logAsJSON(t testing.TB, v interface{}) {
	if !testing.Verbose() {
		return
	}
	t.Helper()
	j, _ := json.MarshalIndent(v, "", "  ")
	t.Log(string(j))
}

func TestProcesses(t *testing.T) {
	start := time.Now()
	procs, err := Processes()
	t.Log("Processes() took", time.Since(start))
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Found", len(procs), "processes.")
	for _, proc := range procs {
		info, err := proc.Info()
		switch {
		// Ignore processes that no longer exist or that cannot be accessed.
		case errors.Is(err, syscall.ESRCH),
			errors.Is(err, syscall.EPERM),
			errors.Is(err, syscall.EINVAL),
			errors.Is(err, syscall.ENOENT),
			errors.Is(err, fs.ErrPermission):
			continue
		case err != nil:
			t.Fatalf("failed to get process info for PID=%d: %v", proc.PID(), err)
		}

		t.Logf("pid=%v name='%s' exe='%s' args=%+v ppid=%d cwd='%s' start_time=%v",
			info.PID, info.Name, info.Exe, info.Args, info.PPID, info.CWD,
			info.StartTime)
	}
}
