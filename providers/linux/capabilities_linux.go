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
	"strconv"

	"github.com/elastic/go-sysinfo/types"
)

// capabilityNames is mapping of capability constant values to names.
//
// Generated with:
//
//	curl -s https://raw.githubusercontent.com/torvalds/linux/master/include/uapi/linux/capability.h | \
//	grep -P '^#define CAP_\w+\s+\d+' | \
//	perl -pe 's/#define (\w+)\s+(\d+)/\2: "\1",/g'

var capabilityNames = map[int]string{
	0:  "CAP_CHOWN",
	1:  "CAP_DAC_OVERRIDE",
	2:  "CAP_DAC_READ_SEARCH",
	3:  "CAP_FOWNER",
	4:  "CAP_FSETID",
	5:  "CAP_KILL",
	6:  "CAP_SETGID",
	7:  "CAP_SETUID",
	8:  "CAP_SETPCAP",
	9:  "CAP_LINUX_IMMUTABLE",
	10: "CAP_NET_BIND_SERVICE",
	11: "CAP_NET_BROADCAST",
	12: "CAP_NET_ADMIN",
	13: "CAP_NET_RAW",
	14: "CAP_IPC_LOCK",
	15: "CAP_IPC_OWNER",
	16: "CAP_SYS_MODULE",
	17: "CAP_SYS_RAWIO",
	18: "CAP_SYS_CHROOT",
	19: "CAP_SYS_PTRACE",
	20: "CAP_SYS_PACCT",
	21: "CAP_SYS_ADMIN",
	22: "CAP_SYS_BOOT",
	23: "CAP_SYS_NICE",
	24: "CAP_SYS_RESOURCE",
	25: "CAP_SYS_TIME",
	26: "CAP_SYS_TTY_CONFIG",
	27: "CAP_MKNOD",
	28: "CAP_LEASE",
	29: "CAP_AUDIT_WRITE",
	30: "CAP_AUDIT_CONTROL",
	31: "CAP_SETFCAP",
	32: "CAP_MAC_OVERRIDE",
	33: "CAP_MAC_ADMIN",
	34: "CAP_SYSLOG",
	35: "CAP_WAKE_ALARM",
	36: "CAP_BLOCK_SUSPEND",
	37: "CAP_AUDIT_READ",
	38: "CAP_PERFMON",
	39: "CAP_BPF",
	40: "CAP_CHECKPOINT_RESTORE",
}

func capabilityName(num int) string {
	name, found := capabilityNames[num]
	if found {
		return name
	}

	return strconv.Itoa(num)
}

func readCapabilities(content []byte) (*types.CapabilityInfo, error) {
	var cap types.CapabilityInfo

	err := parseKeyValue(content, ':', func(key, value []byte) error {
		var err error
		switch string(key) {
		case "CapInh":
			cap.Inheritable, err = decodeBitMap(string(value), capabilityName)
			if err != nil {
				return err
			}
		case "CapPrm":
			cap.Permitted, err = decodeBitMap(string(value), capabilityName)
			if err != nil {
				return err
			}
		case "CapEff":
			cap.Effective, err = decodeBitMap(string(value), capabilityName)
			if err != nil {
				return err
			}
		case "CapBnd":
			cap.Bounding, err = decodeBitMap(string(value), capabilityName)
			if err != nil {
				return err
			}
		case "CapAmb":
			cap.Ambient, err = decodeBitMap(string(value), capabilityName)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return &cap, err
}
