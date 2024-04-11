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

//go:build freebsd && cgo

package freebsd

/*
#cgo LDFLAGS: -lkvm
#include <sys/cdefs.h>
#include <sys/types.h>
#include <sys/sysctl.h>

#include <paths.h>
#include <kvm.h>
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/unix"
)

// kvmGetSwapInfo returns swap summary statistics for the system. It accesses
// the kernel virtual memory (kvm) images by using libkvm.
func kvmGetSwapInfo() (*kvmSwap, error) {
	// Obtain a KVM file descriptor.
	var errstr *C.char
	kd := C.kvm_open(nil, nil, nil, unix.O_RDONLY, errstr)
	if errstr != nil {
		return nil, fmt.Errorf("failed calling kvm_open: %s", C.GoString(errstr))
	}
	defer C.kvm_close(kd)

	var swap kvmSwap
	if n, err := C.kvm_getswapinfo(kd, (*C.struct_kvm_swap)(unsafe.Pointer(&swap)), 1, 0); n != 0 {
		return nil, fmt.Errorf("failed to get kvm_getswapinfo: %w", err)
	}

	return &swap, nil
}
