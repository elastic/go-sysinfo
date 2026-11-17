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
	"unsafe"

	"golang.org/x/sys/unix"
)

// MachineID returns the Hardware UUID also accessible via
// About this Mac -> System Report and as the field
// IOPlatformUUID in the output of "ioreg -d2 -c IOPlatformExpertDevice".
func MachineID() (string, error) {
	return getHostUUID()
}

// getHostUUID calls the gethostuuid syscall (what libc's gethostuuid wraps)
// and formats the result like uuid_unparse_upper.
func getHostUUID() (string, error) {
	var id [16]byte
	wait := unix.Timespec{Sec: 5} // 5 seconds

	_, _, errno := unix.Syscall(unix.SYS_GETHOSTUUID,
		uintptr(unsafe.Pointer(&id[0])), uintptr(unsafe.Pointer(&wait)), 0)
	if errno != 0 {
		return "", fmt.Errorf("gethostuuid failed: %w", errno)
	}

	return fmt.Sprintf("%X-%X-%X-%X-%X", id[0:4], id[4:6], id[6:8], id[8:10], id[10:16]), nil
}
