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
	"fmt"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
)

// The Mach host_statistics calls have no syscall equivalent, so they are
// resolved from libSystem at runtime with purego. This works with
// CGO_ENABLED=0.

const (
	libSystemPath = "/usr/lib/libSystem.B.dylib"

	kernSuccess = 0

	hostCPULoadInfoFlavor = 3  // HOST_CPU_LOAD_INFO
	hostVMInfo64Flavor    = 4  // HOST_VM_INFO64
	hostCPULoadInfoCount  = 4  // sizeof(host_cpu_load_info_data_t) / sizeof(integer_t)
	hostVMInfo64Count     = 38 // sizeof(vm_statistics64_data_t) / sizeof(integer_t)

	// clockTicks is sysconf(_SC_CLK_TCK) on macOS, which is the constant CLK_TCK.
	clockTicks = 100
)

type machFuncs struct {
	hostSelf         func() uint32
	hostStatistics   func(host uint32, flavor int32, out unsafe.Pointer, count *uint32) int32
	hostStatistics64 func(host uint32, flavor int32, out unsafe.Pointer, count *uint32) int32

	host uint32
}

var (
	machOnce sync.Once
	mach     machFuncs
	machErr  error
)

func loadMach() (*machFuncs, error) {
	machOnce.Do(func() {
		lib, err := purego.Dlopen(libSystemPath, purego.RTLD_LAZY|purego.RTLD_GLOBAL)
		if err != nil {
			machErr = fmt.Errorf("loading %s: %w", libSystemPath, err)
			return
		}
		purego.RegisterLibFunc(&mach.hostSelf, lib, "mach_host_self")
		purego.RegisterLibFunc(&mach.hostStatistics, lib, "host_statistics")
		purego.RegisterLibFunc(&mach.hostStatistics64, lib, "host_statistics64")
		// mach_host_self takes a new port reference on every call; fetch it once.
		mach.host = mach.hostSelf()
	})
	return &mach, machErr
}

func getHostCPULoadInfo() (*cpuUsage, error) {
	m, err := loadMach()
	if err != nil {
		return nil, err
	}

	var cpu cpuUsage
	count := uint32(hostCPULoadInfoCount)
	if status := m.hostStatistics(m.host, hostCPULoadInfoFlavor, unsafe.Pointer(&cpu), &count); status != kernSuccess {
		return nil, fmt.Errorf("host_statistics returned status %d", status)
	}

	return &cpu, nil
}

// getClockTicks returns the number of click ticks in one jiffie.
func getClockTicks() int {
	return clockTicks
}

func getHostVMInfo64() (*vmStatistics64Data, error) {
	m, err := loadMach()
	if err != nil {
		return nil, err
	}

	var vmStat vmStatistics64Data
	count := uint32(hostVMInfo64Count)
	if status := m.hostStatistics64(m.host, hostVMInfo64Flavor, unsafe.Pointer(&vmStat), &count); status != kernSuccess {
		return nil, fmt.Errorf("host_statistics64 returned status %d", status)
	}

	return &vmStat, nil
}
