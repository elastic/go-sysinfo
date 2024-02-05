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

package windows

import (
	"golang.org/x/sys/windows"

	go_windows "github.com/elastic/go-windows"
)

const (
	imageFileMachineAmd64 = 0x8664
	imageFileMachineArm64 = 0xAA64
	archIntel             = "x86_64"
	archArm64             = "arm64"
)

func Architecture() (string, error) {
	systemInfo, err := go_windows.GetNativeSystemInfo()
	if err != nil {
		return "", err
	}

	return systemInfo.ProcessorArchitecture.String(), nil
}

func NativeArchitecture() (string, error) {
	var processMachine, nativeMachine uint16
	// the pseudo handle doesn't need to be closed
	var currentProcessHandle = windows.CurrentProcess()

	err := windows.IsWow64Process2(currentProcessHandle, &processMachine, &nativeMachine)
	if err != nil {
		return "", err
	}

	nativeArch := ""

	switch nativeMachine {
	case imageFileMachineAmd64:
		// for parity with Architecture() as amd64 and x86_64 are used interchangeably
		nativeArch = archIntel
	case imageFileMachineArm64:
		nativeArch = archArm64
	default:
	}

	return nativeArch, nil
}
