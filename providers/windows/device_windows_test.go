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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeviceMapper(t *testing.T) {
	devMapper := deviceMapper{
		deviceProvider: testingDeviceProvider(map[byte]string{
			'A': `\Device\Floppy0`,
			'B': `\Device\Floppy1`,
			'C': `\Device\Harddisk0Volume2`,
			'D': `\Device\Harddisk1Volume1`,
			'E': `\Device\Cdrom0`,
			// Virtualbox-style share
			'W': `\Device\Share;w:\dataserver\programs`,
			// Network share
			'Z': `\Device\LANManRedirector;z:01234567812313123\officeserver\documents`,
		}),
	}
	for testIdx, testCase := range []struct {
		devicePath, expected string
	}{
		{`\DEVICE\FLOPPY0\README.TXT`, `A:\README.TXT`},
		{`\Device\cdrom0\autorun.INF`, `E:\autorun.INF`},
		{`\Device\Harddisk0Volume2\WINDOWS\System32\drivers\etc\hosts`, `C:\WINDOWS\System32\drivers\etc\hosts`},
		{`\Device\share\DATASERVER\PROGRAMS\elastic\packetbeat\PACKETBEAT.EXE`, `W:\elastic\packetbeat\PACKETBEAT.EXE`},
		{`\Device\MUP\OfficeServer\Documents\report.pdf`, `Z:\report.pdf`},
		{`\Device\share\othershare\files\run.EXE`, ``},
		{`\Device\MUP\networkserver\share\.git`, `\\networkserver\share\.git`},
		{`\Device\Harddisk1Volume1`, `D:\`},
		{`\Device\Harddisk1Volume1\`, `D:\`},
		{`\Device`, ``},
		{`C:\windows\calc.exe`, ``},
	} {
		msg := fmt.Sprintf("test case #%d: %v", testIdx, testCase)
		path, err := devMapper.DevicePathToDrivePath(testCase.devicePath)
		if err == nil {
			assert.Equal(t, testCase.expected, path, msg)
		} else {
			if len(testCase.expected) != 0 {
				t.Fatal(err, msg)
				continue
			}
			assert.Equal(t, testCase.expected, path, msg)
		}
	}
}
