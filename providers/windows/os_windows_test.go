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

	"github.com/stretchr/testify/assert"

	"github.com/elastic/go-sysinfo/types"
)

func TestFixWindows11Naming(t *testing.T) {
	testCases := []struct {
		osInfo       types.OSInfo
		expectedName string
	}{
		{
			osInfo: types.OSInfo{
				Major: 10,
				Minor: 0,
				Build: "22000",
				Name:  "Windows 10 Pro",
			},
			expectedName: "Windows 11 Pro",
		},
		{
			osInfo: types.OSInfo{
				Major: 10,
				Minor: 0,
				Build: "22001",
				Name:  "Windows 10 Pro",
			},
			expectedName: "Windows 11 Pro",
		},
		{
			osInfo: types.OSInfo{
				Major: 10,
				Minor: 1,
				Build: "0",
				Name:  "Windows 10 Pro",
			},
			expectedName: "Windows 11 Pro",
		},
		{
			osInfo: types.OSInfo{
				Major: 11,
				Minor: 0,
				Build: "0",
				Name:  "Windows 10 Pro",
			},
			expectedName: "Windows 11 Pro",
		},
		{
			osInfo: types.OSInfo{
				Major: 11,
				Minor: 0,
				Build: "0",
				Name:  "Windows 12 Pro",
			},
			expectedName: "Windows 12 Pro",
		},
		{
			osInfo: types.OSInfo{
				Major: 9,
				Minor: 0,
				Build: "22000",
				Name:  "Windows 10 Pro",
			},
			expectedName: "Windows 10 Pro",
		},
	}

	for _, tc := range testCases {
		fixWindows11Naming(tc.osInfo.Build, &tc.osInfo)
		assert.Equal(t, tc.expectedName, tc.osInfo.Name)
	}
}
