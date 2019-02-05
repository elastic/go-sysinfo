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

// +build freebsd,cgo

package freebsd

import "C"

import (
	"strconv"
	"strings"
	"syscall"

	"github.com/elastic/go-sysinfo/types"
)

const ostypeMIB = "kern.ostype"
const osreleaseMIB = "kern.osrelease"
const osrevisionMIB = "kern.osrevision"

func OperatingSystem() (*types.OSInfo, error) {
	info := &types.OSInfo{
		Family:   "freebsd",
		Platform: "freebsd",
	}

	ostype, err := syscall.Sysctl(ostypeMIB)
	if err != nil {
		return info, err
	}
	info.Name = ostype

	osrelease, err := syscall.Sysctl(osreleaseMIB)
	if err != nil {
		return info, err
	}
	info.Version = osrelease

	elems := strings.Split(osrelease, "-")
	majorminor := strings.Split(elems[0], ".")

	if len(majorminor) > 0 {
		info.Major, _ = strconv.Atoi(majorminor[0])
	}

	if len(majorminor) > 1 {
		info.Minor, _ = strconv.Atoi(majorminor[1])
	}

	info.Patch, _ = strconv.Atoi(strings.TrimPrefix(elems[2], "p"))
	return info, nil
}
