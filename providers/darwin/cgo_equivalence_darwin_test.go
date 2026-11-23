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

// TestCgoEquivalentLayout checks that every field of the Go mirror structs sits
// at exactly the offset the SDK headers give it. Unlike comparing live values,
// this is deterministic, and it would catch swapped or shifted fields.
func TestCgoEquivalentLayout(t *testing.T) {
	goOffsets := map[string]uintptr{
		"taskallinfo.Pbsd.Pbi_flags":                            unsafe.Offsetof(procTaskAllInfo{}.Pbsd.Pbi_flags),
		"taskallinfo.Pbsd.Pbi_status":                           unsafe.Offsetof(procTaskAllInfo{}.Pbsd.Pbi_status),
		"taskallinfo.Pbsd.Pbi_xstatus":                          unsafe.Offsetof(procTaskAllInfo{}.Pbsd.Pbi_xstatus),
		"taskallinfo.Pbsd.Pbi_pid":                              unsafe.Offsetof(procTaskAllInfo{}.Pbsd.Pbi_pid),
		"taskallinfo.Pbsd.Pbi_ppid":                             unsafe.Offsetof(procTaskAllInfo{}.Pbsd.Pbi_ppid),
		"taskallinfo.Pbsd.Pbi_uid":                              unsafe.Offsetof(procTaskAllInfo{}.Pbsd.Pbi_uid),
		"taskallinfo.Pbsd.Pbi_gid":                              unsafe.Offsetof(procTaskAllInfo{}.Pbsd.Pbi_gid),
		"taskallinfo.Pbsd.Pbi_ruid":                             unsafe.Offsetof(procTaskAllInfo{}.Pbsd.Pbi_ruid),
		"taskallinfo.Pbsd.Pbi_rgid":                             unsafe.Offsetof(procTaskAllInfo{}.Pbsd.Pbi_rgid),
		"taskallinfo.Pbsd.Pbi_svuid":                            unsafe.Offsetof(procTaskAllInfo{}.Pbsd.Pbi_svuid),
		"taskallinfo.Pbsd.Pbi_svgid":                            unsafe.Offsetof(procTaskAllInfo{}.Pbsd.Pbi_svgid),
		"taskallinfo.Pbsd.Pbi_comm":                             unsafe.Offsetof(procTaskAllInfo{}.Pbsd.Pbi_comm),
		"taskallinfo.Pbsd.Pbi_name":                             unsafe.Offsetof(procTaskAllInfo{}.Pbsd.Pbi_name),
		"taskallinfo.Pbsd.Pbi_nfiles":                           unsafe.Offsetof(procTaskAllInfo{}.Pbsd.Pbi_nfiles),
		"taskallinfo.Pbsd.Pbi_pgid":                             unsafe.Offsetof(procTaskAllInfo{}.Pbsd.Pbi_pgid),
		"taskallinfo.Pbsd.Pbi_pjobc":                            unsafe.Offsetof(procTaskAllInfo{}.Pbsd.Pbi_pjobc),
		"taskallinfo.Pbsd.E_tdev":                               unsafe.Offsetof(procTaskAllInfo{}.Pbsd.E_tdev),
		"taskallinfo.Pbsd.E_tpgid":                              unsafe.Offsetof(procTaskAllInfo{}.Pbsd.E_tpgid),
		"taskallinfo.Pbsd.Pbi_nice":                             unsafe.Offsetof(procTaskAllInfo{}.Pbsd.Pbi_nice),
		"taskallinfo.Pbsd.Pbi_start_tvsec":                      unsafe.Offsetof(procTaskAllInfo{}.Pbsd.Pbi_start_tvsec),
		"taskallinfo.Pbsd.Pbi_start_tvusec":                     unsafe.Offsetof(procTaskAllInfo{}.Pbsd.Pbi_start_tvusec),
		"taskallinfo.Ptinfo.Virtual_size":                       unsafe.Offsetof(procTaskAllInfo{}.Ptinfo.Virtual_size),
		"taskallinfo.Ptinfo.Resident_size":                      unsafe.Offsetof(procTaskAllInfo{}.Ptinfo.Resident_size),
		"taskallinfo.Ptinfo.Total_user":                         unsafe.Offsetof(procTaskAllInfo{}.Ptinfo.Total_user),
		"taskallinfo.Ptinfo.Total_system":                       unsafe.Offsetof(procTaskAllInfo{}.Ptinfo.Total_system),
		"taskallinfo.Ptinfo.Threads_user":                       unsafe.Offsetof(procTaskAllInfo{}.Ptinfo.Threads_user),
		"taskallinfo.Ptinfo.Threads_system":                     unsafe.Offsetof(procTaskAllInfo{}.Ptinfo.Threads_system),
		"taskallinfo.Ptinfo.Policy":                             unsafe.Offsetof(procTaskAllInfo{}.Ptinfo.Policy),
		"taskallinfo.Ptinfo.Faults":                             unsafe.Offsetof(procTaskAllInfo{}.Ptinfo.Faults),
		"taskallinfo.Ptinfo.Pageins":                            unsafe.Offsetof(procTaskAllInfo{}.Ptinfo.Pageins),
		"taskallinfo.Ptinfo.Cow_faults":                         unsafe.Offsetof(procTaskAllInfo{}.Ptinfo.Cow_faults),
		"taskallinfo.Ptinfo.Messages_sent":                      unsafe.Offsetof(procTaskAllInfo{}.Ptinfo.Messages_sent),
		"taskallinfo.Ptinfo.Messages_received":                  unsafe.Offsetof(procTaskAllInfo{}.Ptinfo.Messages_received),
		"taskallinfo.Ptinfo.Syscalls_mach":                      unsafe.Offsetof(procTaskAllInfo{}.Ptinfo.Syscalls_mach),
		"taskallinfo.Ptinfo.Syscalls_unix":                      unsafe.Offsetof(procTaskAllInfo{}.Ptinfo.Syscalls_unix),
		"taskallinfo.Ptinfo.Csw":                                unsafe.Offsetof(procTaskAllInfo{}.Ptinfo.Csw),
		"taskallinfo.Ptinfo.Threadnum":                          unsafe.Offsetof(procTaskAllInfo{}.Ptinfo.Threadnum),
		"taskallinfo.Ptinfo.Numrunning":                         unsafe.Offsetof(procTaskAllInfo{}.Ptinfo.Numrunning),
		"taskallinfo.Ptinfo.Priority":                           unsafe.Offsetof(procTaskAllInfo{}.Ptinfo.Priority),
		"taskallinfo.Ptinfo":                                    unsafe.Offsetof(procTaskAllInfo{}.Ptinfo),
		"vnodepathinfo.Cdir.Path":                               unsafe.Offsetof(procVnodePathInfo{}.Cdir.Path),
		"vnodepathinfo.Cdir":                                    unsafe.Offsetof(procVnodePathInfo{}.Cdir),
		"vnodepathinfo.Rdir.Path":                               unsafe.Offsetof(procVnodePathInfo{}.Rdir.Path),
		"vnodepathinfo.Rdir":                                    unsafe.Offsetof(procVnodePathInfo{}.Rdir),
		"vmstatistics64.Free_count":                             unsafe.Offsetof(vmStatistics64Data{}.Free_count),
		"vmstatistics64.Active_count":                           unsafe.Offsetof(vmStatistics64Data{}.Active_count),
		"vmstatistics64.Inactive_count":                         unsafe.Offsetof(vmStatistics64Data{}.Inactive_count),
		"vmstatistics64.Wire_count":                             unsafe.Offsetof(vmStatistics64Data{}.Wire_count),
		"vmstatistics64.Zero_fill_count":                        unsafe.Offsetof(vmStatistics64Data{}.Zero_fill_count),
		"vmstatistics64.Reactivations":                          unsafe.Offsetof(vmStatistics64Data{}.Reactivations),
		"vmstatistics64.Pageins":                                unsafe.Offsetof(vmStatistics64Data{}.Pageins),
		"vmstatistics64.Pageouts":                               unsafe.Offsetof(vmStatistics64Data{}.Pageouts),
		"vmstatistics64.Faults":                                 unsafe.Offsetof(vmStatistics64Data{}.Faults),
		"vmstatistics64.Cow_faults":                             unsafe.Offsetof(vmStatistics64Data{}.Cow_faults),
		"vmstatistics64.Lookups":                                unsafe.Offsetof(vmStatistics64Data{}.Lookups),
		"vmstatistics64.Hits":                                   unsafe.Offsetof(vmStatistics64Data{}.Hits),
		"vmstatistics64.Purges":                                 unsafe.Offsetof(vmStatistics64Data{}.Purges),
		"vmstatistics64.Purgeable_count":                        unsafe.Offsetof(vmStatistics64Data{}.Purgeable_count),
		"vmstatistics64.Speculative_count":                      unsafe.Offsetof(vmStatistics64Data{}.Speculative_count),
		"vmstatistics64.Decompressions":                         unsafe.Offsetof(vmStatistics64Data{}.Decompressions),
		"vmstatistics64.Compressions":                           unsafe.Offsetof(vmStatistics64Data{}.Compressions),
		"vmstatistics64.Swapins":                                unsafe.Offsetof(vmStatistics64Data{}.Swapins),
		"vmstatistics64.Swapouts":                               unsafe.Offsetof(vmStatistics64Data{}.Swapouts),
		"vmstatistics64.Compressor_page_count":                  unsafe.Offsetof(vmStatistics64Data{}.Compressor_page_count),
		"vmstatistics64.Throttled_count":                        unsafe.Offsetof(vmStatistics64Data{}.Throttled_count),
		"vmstatistics64.External_page_count":                    unsafe.Offsetof(vmStatistics64Data{}.External_page_count),
		"vmstatistics64.Internal_page_count":                    unsafe.Offsetof(vmStatistics64Data{}.Internal_page_count),
		"vmstatistics64.Total_uncompressed_pages_in_compressor": unsafe.Offsetof(vmStatistics64Data{}.Total_uncompressed_pages_in_compressor),
		"cpuload.User":                                          unsafe.Offsetof(cpuUsage{}.User),
		"cpuload.System":                                        unsafe.Offsetof(cpuUsage{}.System),
		"cpuload.Idle":                                          unsafe.Offsetof(cpuUsage{}.Idle),
		"cpuload.Nice":                                          unsafe.Offsetof(cpuUsage{}.Nice),
	}

	require.Len(t, goOffsets, len(cgoref.Offsets), "Go and C field lists differ")
	for name, want := range cgoref.Offsets {
		got, ok := goOffsets[name]
		if !assert.True(t, ok, "%s missing from Go offsets", name) {
			continue
		}
		assert.Equal(t, want, got, "offset of %s", name)
	}
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
			// Only compare identity fields. Pbi_flags, Pbi_status, Pbi_nfiles
			// and the task counters are runtime state and can change between
			// the two reads. Field placement is covered by TestCgoEquivalentLayout.
			want := *(*procTaskAllInfo)(unsafe.Pointer(&wantTask[0]))
			assert.Equal(t, want.Pbsd.Pbi_pid, task.Pbsd.Pbi_pid, "pid %d pbi_pid", pid)
			assert.Equal(t, want.Pbsd.Pbi_ppid, task.Pbsd.Pbi_ppid, "pid %d pbi_ppid", pid)
			assert.Equal(t, want.Pbsd.Pbi_uid, task.Pbsd.Pbi_uid, "pid %d pbi_uid", pid)
			assert.Equal(t, want.Pbsd.Pbi_gid, task.Pbsd.Pbi_gid, "pid %d pbi_gid", pid)
			assert.Equal(t, want.Pbsd.Pbi_ruid, task.Pbsd.Pbi_ruid, "pid %d pbi_ruid", pid)
			assert.Equal(t, want.Pbsd.Pbi_rgid, task.Pbsd.Pbi_rgid, "pid %d pbi_rgid", pid)
			assert.Equal(t, want.Pbsd.Pbi_svuid, task.Pbsd.Pbi_svuid, "pid %d pbi_svuid", pid)
			assert.Equal(t, want.Pbsd.Pbi_svgid, task.Pbsd.Pbi_svgid, "pid %d pbi_svgid", pid)
			assert.Equal(t, want.Pbsd.Pbi_comm, task.Pbsd.Pbi_comm, "pid %d pbi_comm", pid)
			assert.Equal(t, want.Pbsd.Pbi_name, task.Pbsd.Pbi_name, "pid %d pbi_name", pid)
			assert.Equal(t, want.Pbsd.Pbi_pgid, task.Pbsd.Pbi_pgid, "pid %d pbi_pgid", pid)
			assert.Equal(t, want.Pbsd.Pbi_start_tvsec, task.Pbsd.Pbi_start_tvsec, "pid %d pbi_start_tvsec", pid)
			assert.Equal(t, want.Pbsd.Pbi_start_tvusec, task.Pbsd.Pbi_start_tvusec, "pid %d pbi_start_tvusec", pid)
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
	// Sanity check only: page counts drift between reads, so this cannot prove
	// field placement. TestCgoEquivalentLayout does that exactly.
	total := float64(want.Free_count + want.Active_count + want.Inactive_count + want.Wire_count)
	tol := total * 0.05
	assert.InDelta(t, want.Free_count, gotVM.Free_count, tol)
	assert.InDelta(t, want.Active_count, gotVM.Active_count, tol)
	assert.InDelta(t, want.Inactive_count, gotVM.Inactive_count, tol)
	assert.InDelta(t, want.Wire_count, gotVM.Wire_count, tol)
	assert.InDelta(t, float64(want.Total_uncompressed_pages_in_compressor), float64(gotVM.Total_uncompressed_pages_in_compressor), tol)
	assert.GreaterOrEqual(t, gotVM.Faults, want.Faults)
}
