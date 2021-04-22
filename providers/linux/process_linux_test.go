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

	"github.com/elastic/go-sysinfo/internal/registry"
	"github.com/elastic/go-sysinfo/types"
)

var (
	_ registry.HostProvider    = linuxSystem{}
	_ registry.ProcessProvider = linuxSystem{}
)

func TestProcessNetstat(t *testing.T) {
	proc, err := newLinuxSystem("").Self()
	if err != nil {
		t.Fatal(err)
	}
	procNetwork, ok := proc.(types.NetworkCounters)
	if !ok {
		t.Fatalf("error, cannot cast to types.NetworkCounters")
	}
	stats, err := procNetwork.NetworkCounters()
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, stats.SNMP.ICMP, "ICMP")
	assert.NotEmpty(t, stats.SNMP.IP, "IP")
	assert.NotEmpty(t, stats.SNMP.TCP, "TCP")
	assert.NotEmpty(t, stats.SNMP.UDP, "UDP")
}
