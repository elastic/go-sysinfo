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
	"fmt"
	"os"
	"path/filepath"

	windows "github.com/elastic/go-windows"
	"golang.org/x/sys/windows/registry"
)

// fallbackSystemRoot is the last-resort default when the registry query and
// both environment variables are unavailable.
const fallbackSystemRoot = `C:\Windows`

// systemRootFromRegistry reads the SystemRoot value from
// HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion, which reflects the
// actual Windows directory regardless of the process environment. Returns ""
// on any error so the caller can fall back gracefully.
func systemRootFromRegistry() string {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows NT\CurrentVersion`,
		registry.READ|registry.WOW64_64KEY)
	if err != nil {
		return ""
	}
	defer k.Close()
	val, _, err := k.GetStringValue("SystemRoot")
	if err != nil {
		return ""
	}
	return val
}

// kernelExePath returns the absolute path to the running kernel image.
// It prefers the registry (immune to a stripped process environment), then
// falls back to %SystemRoot% / %WINDIR%, then to the hardcoded default.
// See #287.
func kernelExePath() string {
	root := systemRootFromRegistry()
	if root == "" {
		root = os.Getenv("SystemRoot")
	}
	if root == "" {
		root = os.Getenv("WINDIR")
	}
	if root == "" {
		root = fallbackSystemRoot
	}
	return filepath.Join(root, "System32", "ntoskrnl.exe")
}

// KernelVersion returns the version of the running Windows build as
// <major>.<minor>.<build>.<ubr>, for example 10.0.22631.7376.
//
// Feature updates delivered as enablement packages change the build number
// without replacing ntoskrnl.exe, so the file version of the kernel image can
// report the previous feature release. The registry reflects the running
// build, and the file version is only used when it is unavailable.
func KernelVersion() (string, error) {
	if version, err := kernelVersionFromRegistry(); err == nil {
		return version, nil
	}
	return kernelFileVersion()
}

// kernelVersionFromRegistry composes the running build from
// HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion. The major and minor
// version numbers exist since Windows 10.
func kernelVersionFromRegistry() (string, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows NT\CurrentVersion`,
		registry.READ|registry.WOW64_64KEY)
	if err != nil {
		return "", err
	}
	defer k.Close()

	major, _, err := k.GetIntegerValue("CurrentMajorVersionNumber")
	if err != nil {
		return "", err
	}
	minor, _, err := k.GetIntegerValue("CurrentMinorVersionNumber")
	if err != nil {
		return "", err
	}
	build, _, err := k.GetStringValue("CurrentBuildNumber")
	if err != nil {
		return "", err
	}
	ubr, _, err := k.GetIntegerValue("UBR")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d.%d.%s.%d", major, minor, build, ubr), nil
}

// kernelFileVersion returns the file version of the kernel image.
func kernelFileVersion() (string, error) {
	versionData, err := windows.GetFileVersionInfo(kernelExePath())
	if err != nil {
		return "", err
	}

	fileVersion, err := versionData.QueryValue("FileVersion")
	if err == nil {
		return fileVersion, nil
	}

	// Make a second attempt through the fixed version info.
	info, err := versionData.FixedFileInfo()
	if err != nil {
		return "", err
	}
	return info.ProductVersion(), nil
}
