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

package linux

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/elastic/go-sysinfo/types"
)

// Possible (current and historic) locations of the machine-id file.
// These will be searched in order.
var machineIDFiles = []string{"/etc/machine-id", "/var/lib/dbus/machine-id", "/var/db/dbus/machine-id"}

func machineID(hostfs string) (string, error) {
	var contents []byte
	var err error

	for _, file := range machineIDFiles {
		contents, err = os.ReadFile(filepath.Join(hostfs, file))
		if err != nil {
			if os.IsNotExist(err) {
				// Try next location
				continue
			}

			// Return with error on any other error
			return "", fmt.Errorf("failed to read %v: %w", file, err)
		}

		// Found it
		break
	}

	if os.IsNotExist(err) {
		// None of the locations existed
		return "", types.ErrNotImplemented
	}

	contents = bytes.TrimSpace(contents)
	return string(contents), nil
}

func MachineIDHostfs(hostfs string) (string, error) {
	return machineID(hostfs)
}

func MachineID() (string, error) {
	return machineID("")
}
