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

// +build !windows

package linux

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elastic/go-sysinfo/types"
)

func TestOperatingSystem(t *testing.T) {
	t.Run("amazon2017.03", func(t *testing.T) {
		os, err := getOSInfo("testdata/amazon2017.03")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, types.OSInfo{
			Family:   "redhat",
			Platform: "amzn",
			Name:     "Amazon Linux AMI",
			Version:  "2017.03",
			Major:    2017,
			Minor:    3,
			Patch:    0,
		}, *os)
		t.Logf("%#v", os)
	})
	t.Run("centos6", func(t *testing.T) {
		os, err := getOSInfo("testdata/centos6")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, types.OSInfo{
			Family:   "redhat",
			Platform: "centos",
			Name:     "CentOS",
			Version:  "6.9 (Final)",
			Major:    6,
			Minor:    9,
			Codename: "Final",
		}, *os)
		t.Logf("%#v", os)
	})
	t.Run("centos7", func(t *testing.T) {
		os, err := getOSInfo("testdata/centos7")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, types.OSInfo{
			Family:   "redhat",
			Platform: "centos",
			Name:     "CentOS Linux",
			Version:  "7 (Core)",
			Major:    7,
			Minor:    4,
			Patch:    1708,
			Codename: "Core",
		}, *os)
		t.Logf("%#v", os)
	})
	t.Run("debian9", func(t *testing.T) {
		os, err := getOSInfo("testdata/debian9")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, types.OSInfo{
			Family:   "debian",
			Platform: "debian",
			Name:     "Debian GNU/Linux",
			Version:  "9 (stretch)",
			Major:    9,
			Codename: "stretch",
		}, *os)
		t.Logf("%#v", os)
	})
	t.Run("raspbian9", func(t *testing.T) {
		os, err := getOSInfo("testdata/raspbian9")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, types.OSInfo{
			Family:   "debian",
			Platform: "raspbian",
			Name:     "Raspbian GNU/Linux",
			Version:  "9 (stretch)",
			Major:    9,
			Codename: "stretch",
		}, *os)
		t.Logf("%#v", os)
	})
	t.Run("redhat7", func(t *testing.T) {
		os, err := getOSInfo("testdata/redhat7")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, types.OSInfo{
			Family:   "redhat",
			Platform: "rhel",
			Name:     "Red Hat Enterprise Linux Server",
			Version:  "7.6 (Maipo)",
			Major:    7,
			Minor:    6,
			Codename: "Maipo",
		}, *os)
		t.Logf("%#v", os)
	})
	t.Run("ubuntu1404", func(t *testing.T) {
		os, err := getOSInfo("testdata/ubuntu1404")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, types.OSInfo{
			Family:   "debian",
			Platform: "ubuntu",
			Name:     "Ubuntu",
			Version:  "14.04.5 LTS, Trusty Tahr",
			Major:    14,
			Minor:    4,
			Patch:    5,
			Codename: "trusty",
		}, *os)
		t.Logf("%#v", os)
	})
	t.Run("ubuntu1710", func(t *testing.T) {
		os, err := getOSInfo("testdata/ubuntu1710")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, types.OSInfo{
			Family:   "debian",
			Platform: "ubuntu",
			Name:     "Ubuntu",
			Version:  "17.10 (Artful Aardvark)",
			Major:    17,
			Minor:    10,
			Patch:    0,
			Codename: "artful",
		}, *os)
		t.Logf("%#v", os)
	})
}
