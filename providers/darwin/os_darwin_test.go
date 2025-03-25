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

package darwin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const SystemVersionPlist = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
        <key>ProductBuildVersion</key>
        <string>16G1114</string>
        <key>ProductCopyright</key>
        <string>1983-2017 Apple Inc.</string>
        <key>ProductName</key>
        <string>Mac OS X</string>
        <key>ProductUserVisibleVersion</key>
        <string>10.12.6</string>
        <key>ProductVersion</key>
        <string>10.12.6</string>
</dict>
</plist>
`

func TestOperatingSystem(t *testing.T) {
	osInfo, err := getOSInfo([]byte(SystemVersionPlist))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "macos", osInfo.Type)
	assert.Equal(t, "darwin", osInfo.Family)
	assert.Equal(t, "darwin", osInfo.Platform)
	assert.Equal(t, "Mac OS X", osInfo.Name)
	assert.Equal(t, "10.12.6", osInfo.Version)
	assert.Equal(t, 10, osInfo.Major)
	assert.Equal(t, 12, osInfo.Minor)
	assert.Equal(t, 6, osInfo.Patch)
	assert.Equal(t, "16G1114", osInfo.Build)
}
