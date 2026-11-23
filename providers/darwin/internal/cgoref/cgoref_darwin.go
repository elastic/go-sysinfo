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

// Package cgoref calls libproc, libc and the Mach host APIs through cgo. It is
// the reference the pure Go darwin provider is compared against in tests, and
// is not used by the provider itself.
package cgoref

/*
#cgo LDFLAGS:-lproc
#include <errno.h>
#include <stddef.h>
#include <libproc.h>
#include <mach/mach_host.h>
#include <mach/mach_time.h>
#include <sys/sysctl.h>
#include <unistd.h>
#include <uuid/uuid.h>

// Field offsets, as the compiler lays them out in the SDK headers.
enum {
	off_bsd_pbi_flags = offsetof(struct proc_bsdinfo, pbi_flags),
	off_bsd_pbi_status = offsetof(struct proc_bsdinfo, pbi_status),
	off_bsd_pbi_xstatus = offsetof(struct proc_bsdinfo, pbi_xstatus),
	off_bsd_pbi_pid = offsetof(struct proc_bsdinfo, pbi_pid),
	off_bsd_pbi_ppid = offsetof(struct proc_bsdinfo, pbi_ppid),
	off_bsd_pbi_uid = offsetof(struct proc_bsdinfo, pbi_uid),
	off_bsd_pbi_gid = offsetof(struct proc_bsdinfo, pbi_gid),
	off_bsd_pbi_ruid = offsetof(struct proc_bsdinfo, pbi_ruid),
	off_bsd_pbi_rgid = offsetof(struct proc_bsdinfo, pbi_rgid),
	off_bsd_pbi_svuid = offsetof(struct proc_bsdinfo, pbi_svuid),
	off_bsd_pbi_svgid = offsetof(struct proc_bsdinfo, pbi_svgid),
	off_bsd_pbi_comm = offsetof(struct proc_bsdinfo, pbi_comm),
	off_bsd_pbi_name = offsetof(struct proc_bsdinfo, pbi_name),
	off_bsd_pbi_nfiles = offsetof(struct proc_bsdinfo, pbi_nfiles),
	off_bsd_pbi_pgid = offsetof(struct proc_bsdinfo, pbi_pgid),
	off_bsd_pbi_pjobc = offsetof(struct proc_bsdinfo, pbi_pjobc),
	off_bsd_e_tdev = offsetof(struct proc_bsdinfo, e_tdev),
	off_bsd_e_tpgid = offsetof(struct proc_bsdinfo, e_tpgid),
	off_bsd_pbi_nice = offsetof(struct proc_bsdinfo, pbi_nice),
	off_bsd_pbi_start_tvsec = offsetof(struct proc_bsdinfo, pbi_start_tvsec),
	off_bsd_pbi_start_tvusec = offsetof(struct proc_bsdinfo, pbi_start_tvusec),
	off_task_pti_virtual_size = offsetof(struct proc_taskinfo, pti_virtual_size),
	off_task_pti_resident_size = offsetof(struct proc_taskinfo, pti_resident_size),
	off_task_pti_total_user = offsetof(struct proc_taskinfo, pti_total_user),
	off_task_pti_total_system = offsetof(struct proc_taskinfo, pti_total_system),
	off_task_pti_threads_user = offsetof(struct proc_taskinfo, pti_threads_user),
	off_task_pti_threads_system = offsetof(struct proc_taskinfo, pti_threads_system),
	off_task_pti_policy = offsetof(struct proc_taskinfo, pti_policy),
	off_task_pti_faults = offsetof(struct proc_taskinfo, pti_faults),
	off_task_pti_pageins = offsetof(struct proc_taskinfo, pti_pageins),
	off_task_pti_cow_faults = offsetof(struct proc_taskinfo, pti_cow_faults),
	off_task_pti_messages_sent = offsetof(struct proc_taskinfo, pti_messages_sent),
	off_task_pti_messages_received = offsetof(struct proc_taskinfo, pti_messages_received),
	off_task_pti_syscalls_mach = offsetof(struct proc_taskinfo, pti_syscalls_mach),
	off_task_pti_syscalls_unix = offsetof(struct proc_taskinfo, pti_syscalls_unix),
	off_task_pti_csw = offsetof(struct proc_taskinfo, pti_csw),
	off_task_pti_threadnum = offsetof(struct proc_taskinfo, pti_threadnum),
	off_task_pti_numrunning = offsetof(struct proc_taskinfo, pti_numrunning),
	off_task_pti_priority = offsetof(struct proc_taskinfo, pti_priority),
	off_ptinfo = offsetof(struct proc_taskallinfo, ptinfo),
	off_vnode_cdir_path = offsetof(struct vnode_info_path, vip_path),
	off_vnode_cdir = offsetof(struct proc_vnodepathinfo, pvi_cdir),
	off_vnode_rdir_path = offsetof(struct vnode_info_path, vip_path),
	off_vnode_rdir = offsetof(struct proc_vnodepathinfo, pvi_rdir),
	off_vm_free_count = offsetof(vm_statistics64_data_t, free_count),
	off_vm_active_count = offsetof(vm_statistics64_data_t, active_count),
	off_vm_inactive_count = offsetof(vm_statistics64_data_t, inactive_count),
	off_vm_wire_count = offsetof(vm_statistics64_data_t, wire_count),
	off_vm_zero_fill_count = offsetof(vm_statistics64_data_t, zero_fill_count),
	off_vm_reactivations = offsetof(vm_statistics64_data_t, reactivations),
	off_vm_pageins = offsetof(vm_statistics64_data_t, pageins),
	off_vm_pageouts = offsetof(vm_statistics64_data_t, pageouts),
	off_vm_faults = offsetof(vm_statistics64_data_t, faults),
	off_vm_cow_faults = offsetof(vm_statistics64_data_t, cow_faults),
	off_vm_lookups = offsetof(vm_statistics64_data_t, lookups),
	off_vm_hits = offsetof(vm_statistics64_data_t, hits),
	off_vm_purges = offsetof(vm_statistics64_data_t, purges),
	off_vm_purgeable_count = offsetof(vm_statistics64_data_t, purgeable_count),
	off_vm_speculative_count = offsetof(vm_statistics64_data_t, speculative_count),
	off_vm_decompressions = offsetof(vm_statistics64_data_t, decompressions),
	off_vm_compressions = offsetof(vm_statistics64_data_t, compressions),
	off_vm_swapins = offsetof(vm_statistics64_data_t, swapins),
	off_vm_swapouts = offsetof(vm_statistics64_data_t, swapouts),
	off_vm_compressor_page_count = offsetof(vm_statistics64_data_t, compressor_page_count),
	off_vm_throttled_count = offsetof(vm_statistics64_data_t, throttled_count),
	off_vm_external_page_count = offsetof(vm_statistics64_data_t, external_page_count),
	off_vm_internal_page_count = offsetof(vm_statistics64_data_t, internal_page_count),
	off_vm_total_uncompressed_pages_in_compressor = offsetof(vm_statistics64_data_t, total_uncompressed_pages_in_compressor),
	off_cpu_user = offsetof(host_cpu_load_info_data_t, cpu_ticks) + CPU_STATE_USER * sizeof(natural_t),
	off_cpu_system = offsetof(host_cpu_load_info_data_t, cpu_ticks) + CPU_STATE_SYSTEM * sizeof(natural_t),
	off_cpu_idle = offsetof(host_cpu_load_info_data_t, cpu_ticks) + CPU_STATE_IDLE * sizeof(natural_t),
	off_cpu_nice = offsetof(host_cpu_load_info_data_t, cpu_ticks) + CPU_STATE_NICE * sizeof(natural_t),
};
*/
import "C"

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	TaskAllInfoSize   = C.sizeof_struct_proc_taskallinfo
	VnodeInfoSize     = C.sizeof_struct_proc_vnodepathinfo
	VMInfo64Size      = C.sizeof_vm_statistics64_data_t
	CPULoadSize       = C.sizeof_host_cpu_load_info_data_t
	VMInfo64Count     = C.HOST_VM_INFO64_COUNT
	VMInfo64Rev1Count = C.HOST_VM_INFO64_REV1_COUNT
	CPULoadCount      = C.HOST_CPU_LOAD_INFO_COUNT
)

func pidinfo(pid int, flavor C.int, size int) ([]byte, error) {
	buf := make([]byte, size)
	n, err := C.proc_pidinfo(C.int(pid), flavor, 0, unsafe.Pointer(&buf[0]), C.int(size))
	if n != C.int(size) {
		if errno, ok := err.(syscall.Errno); ok && errno != 0 {
			return nil, errno
		}
		return nil, fmt.Errorf("proc_pidinfo returned %d", n)
	}
	return buf, nil
}

// ProcTaskAllInfo returns the raw struct proc_taskallinfo for pid.
func ProcTaskAllInfo(pid int) ([]byte, error) {
	return pidinfo(pid, C.PROC_PIDTASKALLINFO, TaskAllInfoSize)
}

// ProcVnodePathInfo returns the raw struct proc_vnodepathinfo for pid.
func ProcVnodePathInfo(pid int) ([]byte, error) {
	return pidinfo(pid, C.PROC_PIDVNODEPATHINFO, VnodeInfoSize)
}

// HostUUID is the cgo gethostuuid + uuid_unparse_upper implementation.
func HostUUID() (string, error) {
	var id [16]C.uchar
	wait := C.struct_timespec{5, 0}
	if ret, err := C.gethostuuid(&id[0], &wait); ret != 0 {
		return "", fmt.Errorf("gethostuuid failed with %v: %v", ret, err)
	}
	var s [37]C.char
	C.uuid_unparse_upper(&id[0], &s[0])
	return C.GoString(&s[0]), nil
}

// ClockTicks is sysconf(_SC_CLK_TCK).
func ClockTicks() int { return int(C.sysconf(C._SC_CLK_TCK)) }

// HostCPULoad returns host_cpu_load_info_data_t as user, system, idle, nice ticks.
func HostCPULoad() ([4]uint32, error) {
	var cpu [4]uint32
	count := C.mach_msg_type_number_t(C.HOST_CPU_LOAD_INFO_COUNT)
	status := C.host_statistics(C.host_t(C.mach_host_self()), C.HOST_CPU_LOAD_INFO,
		C.host_info_t(unsafe.Pointer(&cpu[0])), &count)
	if status != C.KERN_SUCCESS {
		return cpu, fmt.Errorf("host_statistics returned %d", status)
	}
	return cpu, nil
}

// HostVMInfo64 returns the raw vm_statistics64_data_t.
func HostVMInfo64() ([]byte, error) {
	buf := make([]byte, VMInfo64Size)
	count := C.mach_msg_type_number_t(C.HOST_VM_INFO64_COUNT)
	status := C.host_statistics64(C.host_t(C.mach_host_self()), C.HOST_VM_INFO64,
		C.host_info_t(unsafe.Pointer(&buf[0])), &count)
	if status != C.KERN_SUCCESS {
		return nil, fmt.Errorf("host_statistics64 returned %d", status)
	}
	return buf, nil
}

// Offsets maps "<struct>.<field path>" to the C offsetof for every field the
// Go mirror structs define.
var Offsets = map[string]uintptr{
	"taskallinfo.Pbsd.Pbi_flags":                            C.off_bsd_pbi_flags,
	"taskallinfo.Pbsd.Pbi_status":                           C.off_bsd_pbi_status,
	"taskallinfo.Pbsd.Pbi_xstatus":                          C.off_bsd_pbi_xstatus,
	"taskallinfo.Pbsd.Pbi_pid":                              C.off_bsd_pbi_pid,
	"taskallinfo.Pbsd.Pbi_ppid":                             C.off_bsd_pbi_ppid,
	"taskallinfo.Pbsd.Pbi_uid":                              C.off_bsd_pbi_uid,
	"taskallinfo.Pbsd.Pbi_gid":                              C.off_bsd_pbi_gid,
	"taskallinfo.Pbsd.Pbi_ruid":                             C.off_bsd_pbi_ruid,
	"taskallinfo.Pbsd.Pbi_rgid":                             C.off_bsd_pbi_rgid,
	"taskallinfo.Pbsd.Pbi_svuid":                            C.off_bsd_pbi_svuid,
	"taskallinfo.Pbsd.Pbi_svgid":                            C.off_bsd_pbi_svgid,
	"taskallinfo.Pbsd.Pbi_comm":                             C.off_bsd_pbi_comm,
	"taskallinfo.Pbsd.Pbi_name":                             C.off_bsd_pbi_name,
	"taskallinfo.Pbsd.Pbi_nfiles":                           C.off_bsd_pbi_nfiles,
	"taskallinfo.Pbsd.Pbi_pgid":                             C.off_bsd_pbi_pgid,
	"taskallinfo.Pbsd.Pbi_pjobc":                            C.off_bsd_pbi_pjobc,
	"taskallinfo.Pbsd.E_tdev":                               C.off_bsd_e_tdev,
	"taskallinfo.Pbsd.E_tpgid":                              C.off_bsd_e_tpgid,
	"taskallinfo.Pbsd.Pbi_nice":                             C.off_bsd_pbi_nice,
	"taskallinfo.Pbsd.Pbi_start_tvsec":                      C.off_bsd_pbi_start_tvsec,
	"taskallinfo.Pbsd.Pbi_start_tvusec":                     C.off_bsd_pbi_start_tvusec,
	"taskallinfo.Ptinfo.Virtual_size":                       C.off_task_pti_virtual_size,
	"taskallinfo.Ptinfo.Resident_size":                      C.off_task_pti_resident_size,
	"taskallinfo.Ptinfo.Total_user":                         C.off_task_pti_total_user,
	"taskallinfo.Ptinfo.Total_system":                       C.off_task_pti_total_system,
	"taskallinfo.Ptinfo.Threads_user":                       C.off_task_pti_threads_user,
	"taskallinfo.Ptinfo.Threads_system":                     C.off_task_pti_threads_system,
	"taskallinfo.Ptinfo.Policy":                             C.off_task_pti_policy,
	"taskallinfo.Ptinfo.Faults":                             C.off_task_pti_faults,
	"taskallinfo.Ptinfo.Pageins":                            C.off_task_pti_pageins,
	"taskallinfo.Ptinfo.Cow_faults":                         C.off_task_pti_cow_faults,
	"taskallinfo.Ptinfo.Messages_sent":                      C.off_task_pti_messages_sent,
	"taskallinfo.Ptinfo.Messages_received":                  C.off_task_pti_messages_received,
	"taskallinfo.Ptinfo.Syscalls_mach":                      C.off_task_pti_syscalls_mach,
	"taskallinfo.Ptinfo.Syscalls_unix":                      C.off_task_pti_syscalls_unix,
	"taskallinfo.Ptinfo.Csw":                                C.off_task_pti_csw,
	"taskallinfo.Ptinfo.Threadnum":                          C.off_task_pti_threadnum,
	"taskallinfo.Ptinfo.Numrunning":                         C.off_task_pti_numrunning,
	"taskallinfo.Ptinfo.Priority":                           C.off_task_pti_priority,
	"taskallinfo.Ptinfo":                                    C.off_ptinfo,
	"vnodepathinfo.Cdir.Path":                               C.off_vnode_cdir_path,
	"vnodepathinfo.Cdir":                                    C.off_vnode_cdir,
	"vnodepathinfo.Rdir.Path":                               C.off_vnode_rdir_path,
	"vnodepathinfo.Rdir":                                    C.off_vnode_rdir,
	"vmstatistics64.Free_count":                             C.off_vm_free_count,
	"vmstatistics64.Active_count":                           C.off_vm_active_count,
	"vmstatistics64.Inactive_count":                         C.off_vm_inactive_count,
	"vmstatistics64.Wire_count":                             C.off_vm_wire_count,
	"vmstatistics64.Zero_fill_count":                        C.off_vm_zero_fill_count,
	"vmstatistics64.Reactivations":                          C.off_vm_reactivations,
	"vmstatistics64.Pageins":                                C.off_vm_pageins,
	"vmstatistics64.Pageouts":                               C.off_vm_pageouts,
	"vmstatistics64.Faults":                                 C.off_vm_faults,
	"vmstatistics64.Cow_faults":                             C.off_vm_cow_faults,
	"vmstatistics64.Lookups":                                C.off_vm_lookups,
	"vmstatistics64.Hits":                                   C.off_vm_hits,
	"vmstatistics64.Purges":                                 C.off_vm_purges,
	"vmstatistics64.Purgeable_count":                        C.off_vm_purgeable_count,
	"vmstatistics64.Speculative_count":                      C.off_vm_speculative_count,
	"vmstatistics64.Decompressions":                         C.off_vm_decompressions,
	"vmstatistics64.Compressions":                           C.off_vm_compressions,
	"vmstatistics64.Swapins":                                C.off_vm_swapins,
	"vmstatistics64.Swapouts":                               C.off_vm_swapouts,
	"vmstatistics64.Compressor_page_count":                  C.off_vm_compressor_page_count,
	"vmstatistics64.Throttled_count":                        C.off_vm_throttled_count,
	"vmstatistics64.External_page_count":                    C.off_vm_external_page_count,
	"vmstatistics64.Internal_page_count":                    C.off_vm_internal_page_count,
	"vmstatistics64.Total_uncompressed_pages_in_compressor": C.off_vm_total_uncompressed_pages_in_compressor,
	"cpuload.User":                                          C.off_cpu_user,
	"cpuload.System":                                        C.off_cpu_system,
	"cpuload.Idle":                                          C.off_cpu_idle,
	"cpuload.Nice":                                          C.off_cpu_nice,
}
