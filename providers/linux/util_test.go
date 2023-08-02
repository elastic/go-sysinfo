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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseKeyValueNoEOL(t *testing.T) {
	vals := [][2]string{}
	err := parseKeyValue([]byte(
		"Name:	zsh\nUmask:	0022\nState:	S (sleeping)\nUid:	1000	1000	1000	1000",
	), ':', func(key, value []byte) error {
		vals = append(vals, [2]string{string(key), string(value)})
		return nil
	})
	assert.NoError(t, err)

	assert.Equal(t, [][2]string{
		{"Name", "zsh"},
		{"Umask", "0022"},
		{"State", "S (sleeping)"},
		{"Uid", "1000\t1000\t1000\t1000"},
	}, vals)
}

func TestParseKeyValueEmptyLine(t *testing.T) {
	vals := [][2]string{}
	err := parseKeyValue([]byte(
		"Name:	zsh\nUmask:	0022\nState:	S (sleeping)\n\nUid:	1000	1000	1000	1000",
	), ':', func(key, value []byte) error {
		vals = append(vals, [2]string{string(key), string(value)})
		return nil
	})
	assert.NoError(t, err)

	assert.Equal(t, [][2]string{
		{"Name", "zsh"},
		{"Umask", "0022"},
		{"State", "S (sleeping)"},
		{"Uid", "1000\t1000\t1000\t1000"},
	}, vals)
}

func TestParseKeyValueEOL(t *testing.T) {
	vals := [][2]string{}
	err := parseKeyValue([]byte(
		"Name:	zsh\nUmask:	0022\nState:	S (sleeping)\nUid:	1000	1000	1000	1000\n",
	), ':', func(key, value []byte) error {
		vals = append(vals, [2]string{string(key), string(value)})
		return nil
	})
	assert.NoError(t, err)

	assert.Equal(t, [][2]string{
		{"Name", "zsh"},
		{"Umask", "0022"},
		{"State", "S (sleeping)"},
		{"Uid", "1000\t1000\t1000\t1000"},
	}, vals)
}

// from cat /proc/$$/status
var testProcStatus = []byte(`Name:	zsh
Umask:	0022
State:	S (sleeping)
Tgid:	4023363
Ngid:	0
Pid:	4023363
PPid:	4023357
TracerPid:	0
Uid:	1000	1000	1000	1000
Gid:	1000	1000	1000	1000
FDSize:	64
Groups:	24 25 27 29 30 44 46 102 109 112 116 119 131 998 1000
NStgid:	4023363
NSpid:	4023363
NSpgid:	4023363
NSsid:	4023363
VmPeak:	   15596 kB
VmSize:	   15144 kB
VmLck:	       0 kB
VmPin:	       0 kB
VmHWM:	    9060 kB
VmRSS:	    8716 kB
RssAnon:	    3828 kB
RssFile:	    4888 kB
RssShmem:	       0 kB
VmData:	    3500 kB
VmStk:	     328 kB
VmExe:	     600 kB
VmLib:	    2676 kB
VmPTE:	      68 kB
VmSwap:	       0 kB
HugetlbPages:	       0 kB
CoreDumping:	0
THP_enabled:	1
Threads:	1
SigQ:	0/126683
SigPnd:	0000000000000000
ShdPnd:	0000000000000000
SigBlk:	0000000000000002
SigIgn:	0000000000384000
SigCgt:	0000000008013003
CapInh:	0000000000000000
CapPrm:	0000000000000000
CapEff:	0000000000000000
CapBnd:	000001ffffffffff
CapAmb:	0000000000000000
NoNewPrivs:	0
Seccomp:	0
Seccomp_filters:	0
Speculation_Store_Bypass:	thread vulnerable
Cpus_allowed:	fff
Cpus_allowed_list:	0-11
Mems_allowed:	00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000001
Mems_allowed_list:	0
voluntary_ctxt_switches:	223
nonvoluntary_ctxt_switches:	25
`)

func BenchmarkParseKeyValue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = parseKeyValue(testProcStatus, ':', func(key, value []byte) error {
			return nil
		})
	}
}

func FuzzParseKeyValue(f *testing.F) {
	testcases := []string{
		"no_separator",
		"no_value:",
		"empty_value: ",
		"normal:	223",
	}
	for _, tc := range testcases {
		f.Add(tc)
	}
	f.Fuzz(func(t *testing.T, orig string) {
		_ = parseKeyValue([]byte(orig), ':', func(key, value []byte) error {
			return nil
		})
	})
}
