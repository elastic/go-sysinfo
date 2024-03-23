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
// +build freebsd,cgo

package freebsd

// #cgo LDFLAGS: -lkvm
//#include <sys/cdefs.h>
//#include <sys/types.h>
//#include <sys/sysctl.h>

//#include <paths.h>
//#include <kvm.h>
//#include <stdlib.h>
import "C"

import (
	"syscall"
	"unsafe"

	"github.com/pkg/errors"
)

const (
	hwPhysmemMIB         = "hw.physmem"
	hwPagesizeMIB        = "hw.pagesize"
	vmVmtotalMIB         = "vm.vmtotal"
	vmSwapmaxpagesMIB    = "vm.swap_maxpages"
	vfsNumfreebuffersMIB = "vfs.numfreebuffers"
	devNull              = "/dev/null"
	kvmOpen              = "kvm_open"
)

func PageSize() (uint32, error) {
	var pageSize uint32
	if err := sysctlByName(hwPagesizeMIB, &pageSize); err != nil {
		return 0, errors.Wrap(err, "failed to get hw.pagesize")
	}

	return pageSize, nil
}

func SwapMaxPages() (uint32, error) {
	var maxPages uint32
	if err := sysctlByName(hwPhysmemMIB, &maxPages); err != nil {
		return 0, errors.Wrap(err, "failed to get vm.swap_maxpages")
	}

	return maxPages, nil
}

func TotalMemory() (uint64, error) {
	var size uint64
	if err := sysctlByName(hwPhysmemMIB, &size); err != nil {
		return 0, errors.Wrap(err, "failed to get hw.physmem")
	}

	return size, nil
}

func VmTotal() (vmTotal, error) {
	var vm vmTotal
	if err := sysctlByName(vmVmtotalMIB, &vm); err != nil {
		return vmTotal{}, errors.Wrap(err, "failed to get vm.vmtotal")
	}

	return vm, nil
}

func NumFreeBuffers() (uint32, error) {
	var numfreebuffers uint32
	if err := sysctlByName(vfsNumfreebuffersMIB, &numfreebuffers); err != nil {
		return 0, errors.Wrap(err, "failed to get vfs.numfreebuffers")
	}

	return numfreebuffers, nil
}

func KvmGetSwapInfo() (kvmSwap, error) {
	var kdC *C.struct_kvm_t

	devNullC := C.CString(devNull)
	defer C.free(unsafe.Pointer(devNullC))
	kvmOpenC := C.CString(kvmOpen)
	defer C.free(unsafe.Pointer(kvmOpenC))

	if kdC, err := C.kvm_open(nil, devNullC, nil, syscall.O_RDONLY, kvmOpenC); kdC == nil {
		return kvmSwap{}, errors.Wrap(err, "failed to open kvm")
	}

	defer C.kvm_close((*C.struct___kvm)(unsafe.Pointer(kdC)))

	var swap kvmSwap
	if n, err := C.kvm_getswapinfo((*C.struct___kvm)(unsafe.Pointer(kdC)), (*C.struct_kvm_swap)(unsafe.Pointer(&swap)), 1, 0); n != 0 {
		return kvmSwap{}, errors.Wrap(err, "failed to get kvm_getswapinfo")
	}

	return swap, nil
}
