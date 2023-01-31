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
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elastic/go-sysinfo/internal/registry"
	"github.com/elastic/go-sysinfo/types"
)

var _ registry.HostProvider = linuxSystem{}

func TestHost(t *testing.T) {
	host, err := newLinuxSystem("").Host()
	if err != nil {
		t.Fatal(err)
	}
	info := host.Info()
	data, _ := json.MarshalIndent(info, "", "  ")
	t.Log(string(data))
}

func TestHostMemoryInfo(t *testing.T) {
	host, err := newLinuxSystem("testdata/ubuntu1710").Host()
	if err != nil {
		t.Fatal(err)
	}
	m, err := host.Memory()
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, 4139057152, m.Total)
	assert.NotContains(t, m.Metrics, "MemTotal")
	assert.Contains(t, m.Metrics, "Slab")
}

func TestHostVMStat(t *testing.T) {
	host, err := newLinuxSystem("testdata/ubuntu1710").Host()
	if err != nil {
		t.Fatal(err)
	}
	s, err := host.(types.VMStat).VMStat()
	if err != nil {
		t.Fatal(err)
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
}

func TestHostLoadAverage(t *testing.T) {
	host, err := newLinuxSystem("testdata/ubuntu1710").Host()
	if err != nil {
		t.Fatal(err)
	}
	s, err := host.(types.LoadAverage).LoadAverage()
	if err != nil {
		t.Fatal(err)
	}

	data, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
}

func TestHostNetworkCounters(t *testing.T) {
	host, err := newLinuxSystem("testdata/fedora30").Host()
	if err != nil {
		t.Fatal(err)
	}

	s, err := host.(types.NetworkCounters).NetworkCounters()
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, s.Netstat.IPExt)
	assert.NotEmpty(t, s.Netstat.TCPExt)
	assert.NotEmpty(t, s.SNMP.IP)

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
}

func TestParseLine(t *testing.T) {
	tcs := []struct {
		name     string
		hostname string
		line     string
		want     string
	}{
		{
			name:     "find fqdn - spaces",
			hostname: "thishost",
			want:     "thishost.mydomain.org",
			line:     "127.0.1.1       thishost.mydomain.org  thishost",
		},
		{
			name:     "find fqdn - tabs",
			hostname: "thishost",
			want:     "thishost.mydomain.org",
			line:     "127.0.1.1	thishost.mydomain.org	thishost",
		},
		{
			name:     "find fqdn - tabs and spaces",
			hostname: "thishost",
			want:     "thishost.mydomain.org",
			line:     "127.0.1.1	 thishost.mydomain.org	  thishost",
		},
		{
			name:     "find fqdn - line with comment",
			hostname: "bar",
			want:     "bar.mydomain.org",
			line:     "192.168.1.13    bar.mydomain.org       bar # comment 1",
		},
		{
			name:     "find fqdn - fqdn with hostname but no alias",
			hostname: "ahostWith",
			want:     "ahostWith.no.alias",
			line:     "213.456.178.9 ahostWith.no.alias",
		},
		{
			name:     "ignore invalid line",
			hostname: "ignore",
			want:     "",
			line:     "209.237.226.91INVALIDwww.opensource.org",
		},
		{
			name:     "comment line",
			hostname: "ignore",
			want:     "",
			line:     "# The following lines are desirable for IPv4 capable hosts",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			got := findInHostsLine(tc.hostname, tc.line)
			if got != tc.want {
				t.Errorf("got %s, want %s", got, tc.want)
			}
		})
	}
}
