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

//go:build !windows

package linux

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elastic/go-sysinfo/types"
)

func TestOperatingSystem(t *testing.T) {
	t.Run("almalinux9", func(t *testing.T) {
		// Data from 'docker pull almalinux:9'.
		os, err := getOSInfo("testdata/almalinux9")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, types.OSInfo{
			Type:     "linux",
			Family:   "redhat",
			Platform: "almalinux",
			Name:     "AlmaLinux",
			Version:  "9.1 (Lime Lynx)",
			Major:    9,
			Minor:    1,
			Codename: "Lime Lynx",
		}, *os)
		t.Logf("%#v", os)
	})
	t.Run("alpine3.17", func(t *testing.T) {
		// Data from 'docker pull alpine:3.17.3'.
		os, err := getOSInfo("testdata/alpine3.17")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, types.OSInfo{
			Type:     "linux",
			Platform: "alpine",
			Name:     "Alpine Linux",
			Version:  "3.17.3",
			Major:    3,
			Minor:    17,
			Patch:    3,
		}, *os)
		t.Logf("%#v", os)
	})
	t.Run("amazon2017.03", func(t *testing.T) {
		os, err := getOSInfo("testdata/amazon2017.03")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, types.OSInfo{
			Type:     "linux",
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
	t.Run("archlinux", func(t *testing.T) {
		os, err := getOSInfo("testdata/archlinux")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, types.OSInfo{
			Type:     "linux",
			Family:   "arch",
			Platform: "archarm",
			Name:     "Arch Linux ARM",
			Build:    "rolling",
		}, *os)
		t.Logf("%#v", os)
	})
	t.Run("centos6", func(t *testing.T) {
		os, err := getOSInfo("testdata/centos6")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, types.OSInfo{
			Type:     "linux",
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
			Type:     "linux",
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
	t.Run("centos7.8", func(t *testing.T) {
		os, err := getOSInfo("testdata/centos7.8")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, types.OSInfo{
			Type:     "linux",
			Family:   "redhat",
			Platform: "centos",
			Name:     "CentOS Linux",
			Version:  "7 (Core)",
			Major:    7,
			Minor:    8,
			Patch:    2003,
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
			Type:     "linux",
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
			Type:     "linux",
			Family:   "debian",
			Platform: "raspbian",
			Name:     "Raspbian GNU/Linux",
			Version:  "9 (stretch)",
			Major:    9,
			Codename: "stretch",
		}, *os)
		t.Logf("%#v", os)
	})
	t.Run("linuxmint20", func(t *testing.T) {
		os, err := getOSInfo("testdata/linuxmint20")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, types.OSInfo{
			Type:     "linux",
			Family:   "debian",
			Platform: "linuxmint",
			Name:     "Linux Mint",
			Version:  "20 (Ulyana)",
			Major:    20,
			Codename: "ulyana",
		}, *os)
		t.Logf("%#v", os)
	})
	t.Run("manjaro23", func(t *testing.T) {
		os, err := getOSInfo("testdata/manjaro23")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, types.OSInfo{
			Type:     "linux",
			Family:   "arch",
			Platform: "manjaro-arm",
			Name:     "Manjaro ARM",
			Version:  "23.02",
			Major:    23,
			Minor:    2,
		}, *os)
		t.Logf("%#v", os)
	})
	t.Run("redhat7", func(t *testing.T) {
		os, err := getOSInfo("testdata/redhat7")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, types.OSInfo{
			Type:     "linux",
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
	t.Run("redhat9", func(t *testing.T) {
		// Data from 'docker pull redhat/ubi9:9.0.0-1468'.
		os, err := getOSInfo("testdata/redhat9")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, types.OSInfo{
			Type:     "linux",
			Family:   "redhat",
			Platform: "rhel",
			Name:     "Red Hat Enterprise Linux",
			Version:  "9.0 (Plow)",
			Major:    9,
			Minor:    0,
			Codename: "Plow",
		}, *os)
		t.Logf("%#v", os)
	})
	t.Run("rockylinux9", func(t *testing.T) {
		// Data from 'docker pull rockylinux:9.0'.
		os, err := getOSInfo("testdata/rockylinux9")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, types.OSInfo{
			Type:     "linux",
			Family:   "redhat",
			Platform: "rocky",
			Name:     "Rocky Linux",
			Version:  "9.0 (Blue Onyx)",
			Major:    9,
			Minor:    0,
			Codename: "Blue Onyx",
		}, *os)
		t.Logf("%#v", os)
	})
	t.Run("openeuler20.03", func(t *testing.T) {
		os, err := getOSInfo("testdata/openeuler20.03")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, types.OSInfo{
			Type:     "linux",
			Family:   "redhat",
			Platform: "openEuler",
			Name:     "openEuler",
			Version:  "20.03 (LTS-SP3)",
			Major:    20,
			Minor:    3,
			Codename: "LTS-SP3",
		}, *os)
		t.Logf("%#v", os)
	})
	t.Run("opensuse-leap15.4", func(t *testing.T) {
		os, err := getOSInfo("testdata/opensuse-leap15.4")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, types.OSInfo{
			Type:     "linux",
			Family:   "suse",
			Platform: "opensuse-leap",
			Name:     "openSUSE Leap",
			Version:  "15.4",
			Major:    15,
			Minor:    4,
		}, *os)
		t.Logf("%#v", os)
	})
	t.Run("opensuse-tumbleweed", func(t *testing.T) {
		os, err := getOSInfo("testdata/opensuse-tumbleweed")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, types.OSInfo{
			Type:     "linux",
			Family:   "suse",
			Platform: "opensuse-tumbleweed",
			Name:     "openSUSE Tumbleweed",
			Version:  "20230108",
			Major:    20230108,
		}, *os)
		t.Logf("%#v", os)
	})
	t.Run("oraclelinux7", func(t *testing.T) {
		os, err := getOSInfo("testdata/oraclelinux7")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, types.OSInfo{
			Type:     "linux",
			Family:   "redhat",
			Platform: "ol",
			Name:     "Oracle Linux Server",
			Version:  "7.9",
			Major:    7,
			Minor:    9,
		}, *os)
		t.Logf("%#v", os)
	})
	t.Run("ubuntu1404", func(t *testing.T) {
		os, err := getOSInfo("testdata/ubuntu1404")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, types.OSInfo{
			Type:     "linux",
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
			Type:     "linux",
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
	t.Run("ubuntu2204", func(t *testing.T) {
		os, err := getOSInfo("testdata/ubuntu2204")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, types.OSInfo{
			Type:     "linux",
			Family:   "debian",
			Platform: "ubuntu",
			Name:     "Ubuntu",
			Version:  "22.04 LTS (Jammy Jellyfish)",
			Major:    22,
			Minor:    4,
			Codename: "jammy",
		}, *os)
		t.Logf("%#v", os)
	})
	t.Run("fedora30", func(t *testing.T) {
		os, err := getOSInfo("testdata/fedora30")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, types.OSInfo{
			Type:     "linux",
			Family:   "redhat",
			Platform: "fedora",
			Name:     "Fedora",
			Version:  "30 (Container Image)",
			Major:    30,
			Minor:    0,
			Patch:    0,
			Codename: "Thirty",
		}, *os)
		t.Logf("%#v", os)
	})
	t.Run("dir_release", func(t *testing.T) {
		os, err := getOSInfo("testdata/dir_release")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, types.OSInfo{
			Type:     "linux",
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
}
