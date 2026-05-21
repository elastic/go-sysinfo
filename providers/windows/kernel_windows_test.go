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
	"testing"
)

func TestKernelExePath(t *testing.T) {
	tests := []struct {
		name       string
		systemRoot string
		windir     string
		want       string
	}{
		{
			// Regression for #287: when the system drive is not C:\ the
			// computed kernel path must follow %SystemRoot% instead of the
			// hardcoded C:\Windows.
			name:       "non-default SystemRoot drive",
			systemRoot: `W:\Windows`,
			want:       `W:\Windows\System32\ntoskrnl.exe`,
		},
		{
			name:       "default SystemRoot",
			systemRoot: `C:\Windows`,
			want:       `C:\Windows\System32\ntoskrnl.exe`,
		},
		{
			name:   "WINDIR fallback when SystemRoot empty",
			windir: `D:\WINNT`,
			want:   `D:\WINNT\System32\ntoskrnl.exe`,
		},
		{
			name: "fallback when neither env var is set",
			want: `C:\Windows\System32\ntoskrnl.exe`,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("SystemRoot", tc.systemRoot)
			t.Setenv("WINDIR", tc.windir)
			if got := kernelExePath(); got != tc.want {
				t.Fatalf("kernelExePath() = %q, want %q", got, tc.want)
			}
		})
	}
}
