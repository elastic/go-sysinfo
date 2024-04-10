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

//go:build freebsd

package freebsd

import (
	"strconv"
	"strings"

	"golang.org/x/sys/unix"

	"github.com/elastic/go-sysinfo/types"
)

const (
	ostypeMIB    = "kern.ostype"
	osreleaseMIB = "kern.osrelease"
)

func OperatingSystem() (*types.OSInfo, error) {
	return getOSInfo("")
}

func getOSInfo(baseDir string) (*types.OSInfo, error) {
	info := &types.OSInfo{
		Type:     "freebsd",
		Family:   "freebsd",
		Platform: "freebsd",
	}

	ostype, err := unix.Sysctl(ostypeMIB)
	if err != nil {
		return info, err
	}
	info.Name = ostype

	// Example: 13.0-RELEASE-p11
	osrelease, err := unix.Sysctl(osreleaseMIB)
	if err != nil {
		return info, err
	}
	info.Version = osrelease

	releaseParts := strings.Split(osrelease, "-")

	majorMinor := strings.Split(releaseParts[0], ".")
	if len(majorMinor) > 0 {
		info.Major, _ = strconv.Atoi(majorMinor[0])
	}
	if len(majorMinor) > 1 {
		info.Minor, _ = strconv.Atoi(majorMinor[1])
	}

	if len(releaseParts) > 1 {
		info.Build = releaseParts[1]
	}
	if len(releaseParts) > 2 {
		info.Patch, _ = strconv.Atoi(strings.TrimPrefix(releaseParts[2], "p"))
	}

	return info, nil
}
