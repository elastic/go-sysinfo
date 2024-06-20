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

//go:build freebsd

package freebsd

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os/exec"
	"testing"
	"time"
)

func TestArchitecture(t *testing.T) {
	arch, err := architecture()
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, arch)
	assert.Regexp(t, `(amd64|i386|powerpc(64(le)?|spe)?|armv(6|7)|aarch64|riscv64*|mips(n)?(32|64)?(el)?|sparc64)`, arch)
}

func TestBootTime(t *testing.T) {
	bootTime, err := bootTime()
	if err != nil {
		t.Fatal(err)
	}

	bootDiff := time.Since(bootTime)
	// t.Logf("bootTime in seconds: %#v", int64(bootDiff.Seconds()))

	cmd := exec.Command("/usr/bin/uptime", "--libxo=json")
	upcmd, err := cmd.Output()

	if err != nil {
		t.Fatal(err)
	}

	t.Logf(string(upcmd))

	type UptimeOutput struct {
		UptimeInformation struct {
			Uptime int64 `json:"uptime"`
		} `json:"uptime-information"`
	}

	var upInfo UptimeOutput
	err = json.Unmarshal(upcmd, &upInfo)

	if err != nil {
		t.Fatal(err)
	}

	upsec := upInfo.UptimeInformation.Uptime
	uptime := time.Duration(upsec * int64(time.Second))
	// t.Logf("uptime in seconds: %#v", int64(uptime.Seconds()))

	assert.InDelta(t, uptime, bootDiff, float64(5*time.Second))
}

func TestCPUStateTimes(t *testing.T) {
	times, err := cpuStateTimes()
	if err != nil {
		t.Fatal(err)
	}

	require.NotNil(t, times)
	assert.NotZero(t, *times)
}

func TestKernelVersion(t *testing.T) {
	kernel, err := kernelVersion()
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, kernel)
}

func TestMachineID(t *testing.T) {
	machineID, err := machineID()
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, machineID)
}

func TestOperatingSystem(t *testing.T) {
	t.Run("freebsd", func(t *testing.T) {
		os, err := operatingSystem()
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "freebsd", os.Type)
		assert.Equal(t, "freebsd", os.Family)
		assert.Equal(t, "freebsd", os.Platform)
		assert.Equal(t, "FreeBSD", os.Name)
		assert.Regexp(t, `\d{1,2}\.\d{1,2}-(RELEASE|STABLE|CURRENT|RC[0-9]|ALPHA(\d{0,2})|BETA(\d{0,2}))(-p\d)?`, os.Version)
		assert.Regexp(t, `\d{1,2}`, os.Major)
		assert.Regexp(t, `\d{1,2}`, os.Minor)
		assert.Regexp(t, `\d{1,2}`, os.Patch)
		assert.Regexp(t, `(RELEASE|STABLE|CURRENT|RC[0-9]|ALPHA([0-9]{0,2})|BETA([0-9]{0,2}))`, os.Build)
		t.Logf("%#v", os)
	})
}
