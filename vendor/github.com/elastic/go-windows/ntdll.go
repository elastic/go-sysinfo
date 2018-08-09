// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

// +build windows

package windows

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	SizeOfProcessBasicInformationStruct = unsafe.Sizeof(ProcessBasicInformationStruct{})
	SizeOfRtlUserProcessParameters      = unsafe.Sizeof(RtlUserProcessParameters{})
)

// NTStatus is an error wrapper for NTSTATUS values, 32bit error-codes returned
// by the NT Kernel.
type NTStatus uint32

// ProcessInformationClass is Go's counterpart for the PROCESSINFOCLASS enumeration
// defined in ntdll.h.
type ProcessInfoClass uint32

const (
	ProcessBasicInformation ProcessInfoClass = iota
	ProcessQuotaLimits
	ProcessIoCounters
	ProcessVmCounters
	ProcessTimes
	ProcessBasePriority
	ProcessRaisePriority
	ProcessDebugPort
	ProcessExceptionPort
	ProcessAccessToken
	ProcessLdtInformation
	ProcessLdtSize
	ProcessDefaultHardErrorMode
	ProcessIoPortHandlers
	ProcessPooledUsageAndLimits
	ProcessWorkingSetWatch
	ProcessUserModeIOPL
	ProcessEnableAlignmentFaultFixup
	ProcessPriorityClass
	ProcessWx86Information
	ProcessHandleCount
	ProcessAffinityMask
	ProcessPriorityBoost
	ProcessDeviceMap
	ProcessSessionInformation
	ProcessForegroundInformation
	ProcessWow64Information
	ProcessImageFileName
	ProcessLUIDDeviceMapsEnabled
	ProcessBreakOnTermination
	ProcessDebugObjectHandle
	ProcessDebugFlags
	ProcessHandleTracing
	ProcessIoPriority
	ProcessExecuteFlags
	ProcessResourceManagement
	ProcessCookie
	ProcessImageInformation
	ProcessCycleTime
	ProcessPagePriority
	ProcessInstrumentationCallback
	ProcessThreadStackAllocation
	ProcessWorkingSetWatchEx
	ProcessImageFileNameWin32
	ProcessImageFileMapping
	ProcessAffinityUpdateMode
	ProcessMemoryAllocationMode
	ProcessGroupInformation
	ProcessTokenVirtualizationEnabled
	ProcessConsoleHostProcess
	ProcessWindowInformation
	ProcessHandleInformation
	ProcessMitigationPolicy
	ProcessDynamicFunctionTableInformation
	ProcessHandleCheckingMode
	ProcessKeepAliveCount
	ProcessRevokeFileHandles
	MaxProcessInfoClass
)

// ProcessBasicInformationStruct is Go's counterpart of the
// PROCESS_BASIC_INFORMATION struct, returned by NtQueryInformationProcess
// when ProcessBasicInformation is requested.
type ProcessBasicInformationStruct struct {
	Reserved1       			 uintptr
	PebBaseAddress  			 uintptr
	Reserved2       			 [2]uintptr
	UniqueProcessId 			 uintptr
	// Undocumented:
	InheritedFromUniqueProcessID uintptr
}

// UnicodeString is Go's equivalent for the _UNICODE_STRING struct.
type UnicodeString struct {
	Size          uint16
	MaximumLength uint16
	Buffer        uintptr
}

// RtlUserProcessParameters is Go's equivalent for the
// _RTL_USER_PROCESS_PARAMETERS struct.
// A few undocumented fields are exposed.
type RtlUserProcessParameters struct {
	Reserved1 [16]byte
	Reserved2 [5]uintptr

	// <undocumented>
	CurrentDirectoryPath   UnicodeString
	CurrentDirectoryHandle uintptr
	DllPath                UnicodeString
	// </undocumented>

	ImagePathName UnicodeString
	CommandLine   UnicodeString
}

// Syscalls
// Warning: NtQueryInformationProcess is an unsupported API that can change
//          in future versions of Windows. Available from XP to Windows 10.
//sys   _NtQueryInformationProcess(handle syscall.Handle, infoClass uint32, info uintptr, infoLen uint32, returnLen *uint32) (ntStatus uint32) = ntdll.NtQueryInformationProcess

// NtQueryInformationProcess is a wrapper for ntdll.NtQueryInformationProcess.
// The handle must have the PROCESS_QUERY_INFORMATION access right.
// Returns an error of type NTStatus.
func NtQueryInformationProcess(handle syscall.Handle, infoClass ProcessInfoClass, info unsafe.Pointer, infoLen uint32) (returnedLen uint32, err error) {
	status := _NtQueryInformationProcess(handle, uint32(infoClass), uintptr(info), infoLen, &returnedLen)
	if status != 0 {
		return returnedLen, NTStatus(status)
	}
	return returnedLen, nil
}

// Error prints the wrapped NTSTATUS in hex form.
func (status NTStatus) Error() string {
	return fmt.Sprintf("ntstatus=%x", uint32(status))
}
