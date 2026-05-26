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
	"path/filepath"
	"testing"
)

// resolveKernelExePath is like kernelExePath but accepts an injected
// systemRoot string so tests can exercise the env-var fallback chain
// without touching the registry.
func resolveKernelExePath(regRoot, systemRoot, windir string) string {
	root := regRoot
	if root == "" {
		root = systemRoot
	}
	if root == "" {
		root = windir
	}
	if root == "" {
		root = fallbackSystemRoot
	}
	return filepath.Join(root, "System32", "ntoskrnl.exe")
}

func TestKernelExePath(t *testing.T) {
	tests := []struct {
		name    string
		regRoot string // simulated registry value (empty = registry miss)
		envRoot string // %SystemRoot%
		windir  string // %WINDIR%
		want    string
	}{
		{
			// Registry wins even when env vars differ -- the primary path.
			name:    "registry value used when present",
			regRoot: `W:\Windows`,
			envRoot: `C:\Windows`,
			want:    `W:\Windows\System32\ntoskrnl.exe`,
		},
		{
			// Regression for #287: registry absent, non-default drive via env.
			name:    "SystemRoot env fallback",
			envRoot: `W:\Windows`,
			want:    `W:\Windows\System32\ntoskrnl.exe`,
		},
		{
			name:   "WINDIR fallback when SystemRoot empty",
			windir: `D:\WINNT`,
			want:   `D:\WINNT\System32\ntoskrnl.exe`,
		},
		{
			name: "hardcoded fallback when all sources absent",
			want: `C:\Windows\System32\ntoskrnl.exe`,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := resolveKernelExePath(tc.regRoot, tc.envRoot, tc.windir)
			if got != tc.want {
				t.Fatalf("resolveKernelExePath() = %q, want %q", got, tc.want)
			}
		})
	}
}

// TestKernelExePathLive checks that the live kernelExePath (registry + env)
// returns a non-empty path ending in ntoskrnl.exe.
func TestKernelExePathLive(t *testing.T) {
	p := kernelExePath()
	if p == "" {
		t.Fatal("kernelExePath() returned empty string")
	}
	if filepath.Base(p) != "ntoskrnl.exe" {
		t.Fatalf("kernelExePath() = %q, want path ending in ntoskrnl.exe", p)
	}
}
