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
#include <libproc.h>
#include <mach/mach_host.h>
#include <mach/mach_time.h>
#include <sys/sysctl.h>
#include <unistd.h>
#include <uuid/uuid.h>
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
