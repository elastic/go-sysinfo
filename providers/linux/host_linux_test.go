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
	s, err := host.VMStat()
	if err != nil {
		t.Fatal(err)
	}

	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
}
