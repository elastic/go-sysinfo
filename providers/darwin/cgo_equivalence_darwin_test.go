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

//go:build darwin && cgo && (amd64 || arm64)

package darwin

import (
	"os"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/elastic/go-sysinfo/providers/darwin/internal/cgoref"
)

// These tests compare the pure Go implementation (which is what is used with
// and without cgo) against libproc/libc/Mach called through cgo.

func TestCgoEquivalentStructSizes(t *testing.T) {
	assert.EqualValues(t, cgoref.TaskAllInfoSize, unsafe.Sizeof(procTaskAllInfo{}))
	assert.EqualValues(t, cgoref.VnodeInfoSize, unsafe.Sizeof(procVnodePathInfo{}))
	// vmStatistics64Data is the REV1 layout. Newer SDKs append
	// fields, so only a prefix is requested from the kernel (count must match
	// our struct size, never the SDK's larger HOST_VM_INFO64_COUNT).
	assert.LessOrEqual(t, uintptr(unsafe.Sizeof(vmStatistics64Data{})), uintptr(cgoref.VMInfo64Size))
	assert.EqualValues(t, cgoref.CPULoadSize, unsafe.Sizeof(cpuUsage{}))
	assert.EqualValues(t, cgoref.VMInfo64Rev1Count, unsafe.Sizeof(vmStatistics64Data{})/4)
	assert.EqualValues(t, cgoref.CPULoadCount, unsafe.Sizeof(cpuUsage{})/4)
	assert.EqualValues(t, hostVMInfo64Count, unsafe.Sizeof(vmStatistics64Data{})/4)
	assert.EqualValues(t, hostCPULoadInfoCount, cgoref.CPULoadCount)
}

func TestCgoEquivalentMachineID(t *testing.T) {
	want, err := cgoref.HostUUID()
	require.NoError(t, err)
	got, err := MachineID()
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestCgoEquivalentClockTicks(t *testing.T) {
	assert.Equal(t, cgoref.ClockTicks(), getClockTicks())
}

func TestCgoEquivalentProcInfo(t *testing.T) {
	for _, pid := range []int{os.Getpid(), os.Getppid(), 1} {
		wantTask, wantTaskErr := cgoref.ProcTaskAllInfo(pid)
		var task procTaskAllInfo
		gotTaskErr := getProcTaskAllInfo(pid, &task)
		require.Equal(t, wantTaskErr == nil, gotTaskErr == nil, "pid %d: cgo=%v go=%v", pid, wantTaskErr, gotTaskErr)
		if wantTaskErr == nil {
			// The bsdinfo half is stable for a running process.
			got := unsafe.Slice((*byte)(unsafe.Pointer(&task)), unsafe.Sizeof(task))
			bsd := unsafe.Sizeof(task.Pbsd)
			assert.Equal(t, wantTask[:bsd], got[:bsd], "pid %d bsdinfo", pid)
		} else {
			assert.Equal(t, wantTaskErr, gotTaskErr, "pid %d errno", pid)
		}

		wantVnode, wantVnodeErr := cgoref.ProcVnodePathInfo(pid)
		var vnode procVnodePathInfo
		gotVnodeErr := getProcVnodePathInfo(pid, &vnode)
		require.Equal(t, wantVnodeErr == nil, gotVnodeErr == nil, "pid %d: cgo=%v go=%v", pid, wantVnodeErr, gotVnodeErr)
		if wantVnodeErr == nil {
			got := unsafe.Slice((*byte)(unsafe.Pointer(&vnode)), unsafe.Sizeof(vnode))
			assert.Equal(t, wantVnode, got, "pid %d vnodepathinfo", pid)
		}
	}
}

func TestCgoEquivalentHostStatistics(t *testing.T) {
	wantCPU, err := cgoref.HostCPULoad()
	require.NoError(t, err)
	gotCPU, err := getHostCPULoadInfo()
	require.NoError(t, err)
	got := [4]uint32{gotCPU.User, gotCPU.System, gotCPU.Idle, gotCPU.Nice}
	for i := range got {
		// Counters only increase; they must be close between two reads.
		assert.GreaterOrEqual(t, got[i], wantCPU[i])
		assert.InDelta(t, wantCPU[i], got[i], 1000)
	}

	wantVM, err := cgoref.HostVMInfo64()
	require.NoError(t, err)
	gotVM, err := getHostVMInfo64()
	require.NoError(t, err)
	want := *(*vmStatistics64Data)(unsafe.Pointer(&wantVM[0]))
	// Page counts drift slightly between reads; allow 1% of the total.
	total := float64(want.Free_count + want.Active_count + want.Inactive_count + want.Wire_count)
	tol := total * 0.01
	assert.InDelta(t, want.Free_count, gotVM.Free_count, tol)
	assert.InDelta(t, want.Active_count, gotVM.Active_count, tol)
	assert.InDelta(t, want.Inactive_count, gotVM.Inactive_count, tol)
	assert.InDelta(t, want.Wire_count, gotVM.Wire_count, tol)
	assert.InDelta(t, float64(want.Total_uncompressed_pages_in_compressor), float64(gotVM.Total_uncompressed_pages_in_compressor), tol)
	assert.GreaterOrEqual(t, gotVM.Faults, want.Faults)
}
