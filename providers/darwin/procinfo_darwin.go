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
	"errors"
	"fmt"
	"unsafe"

	"golang.org/x/sys/unix"
)

//go:generate sh -c "go tool cgo -godefs defs_darwin.go > ztypes_darwin.go"

// proc_info(2) call numbers and flavors from bsd/sys/proc_info.h in XNU.
const (
	procInfoCallPIDInfo = 2 // PROC_INFO_CALL_PIDINFO

	procPIDTaskAllInfo   = 2 // PROC_PIDTASKALLINFO
	procPIDVnodePathInfo = 9 // PROC_PIDVNODEPATHINFO
)

// procInfo invokes the proc_info syscall, which is what libproc's
// proc_pidinfo wraps. Calling it directly avoids cgo while preserving errno.
func procInfo(call, pid, flavor int, arg uint64, buf unsafe.Pointer, size int) (int, error) {
	r, _, errno := unix.Syscall6(unix.SYS_PROC_INFO,
		uintptr(call), uintptr(pid), uintptr(flavor), uintptr(arg),
		uintptr(buf), uintptr(size))
	if errno != 0 {
		return 0, errno
	}
	return int(r), nil
}

func getProcTaskAllInfo(pid int, info *procTaskAllInfo) error {
	size := int(unsafe.Sizeof(*info))
	n, err := procInfo(procInfoCallPIDInfo, pid, procPIDTaskAllInfo, 0, unsafe.Pointer(info), size)
	if err != nil {
		return err
	} else if n != size {
		return errors.New("failed to read process info with proc_pidinfo")
	}
	return nil
}

func getProcVnodePathInfo(pid int, info *procVnodePathInfo) error {
	size := int(unsafe.Sizeof(*info))
	n, err := procInfo(procInfoCallPIDInfo, pid, procPIDVnodePathInfo, 0, unsafe.Pointer(info), size)
	if err != nil {
		return fmt.Errorf("failed to read vnode info with proc_pidinfo: %w", err)
	} else if n != size {
		return errors.New("failed to read vnode info with proc_pidinfo")
	}
	return nil
}
