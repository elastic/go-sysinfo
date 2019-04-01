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

const nonContainerizedCgroup = `11:freezer:/
10:pids:/init.scope
9:memory:/init.scope
8:cpuset:/
7:perf_event:/
6:hugetlb:/
5:blkio:/init.scope
4:net_cls,net_prio:/
3:devices:/init.scope
2:cpu,cpuacct:/init.scope
1:name=systemd:/init.scope
`

const containerCgroup = `14:name=systemd:/docker/81438f4655cd771c425607dcf7654f4dc03c073c0123edc45fcfad28132e8c60
13:pids:/docker/81438f4655cd771c425607dcf7654f4dc03c073c0123edc45fcfad28132e8c60
12:hugetlb:/docker/81438f4655cd771c425607dcf7654f4dc03c073c0123edc45fcfad28132e8c60
11:net_prio:/docker/81438f4655cd771c425607dcf7654f4dc03c073c0123edc45fcfad28132e8c60
10:perf_event:/docker/81438f4655cd771c425607dcf7654f4dc03c073c0123edc45fcfad28132e8c60
9:net_cls:/docker/81438f4655cd771c425607dcf7654f4dc03c073c0123edc45fcfad28132e8c60
8:freezer:/docker/81438f4655cd771c425607dcf7654f4dc03c073c0123edc45fcfad28132e8c60
7:devices:/docker/81438f4655cd771c425607dcf7654f4dc03c073c0123edc45fcfad28132e8c60
6:memory:/docker/81438f4655cd771c425607dcf7654f4dc03c073c0123edc45fcfad28132e8c60
5:blkio:/docker/81438f4655cd771c425607dcf7654f4dc03c073c0123edc45fcfad28132e8c60
4:cpuacct:/docker/81438f4655cd771c425607dcf7654f4dc03c073c0123edc45fcfad28132e8c60
3:cpu:/docker/81438f4655cd771c425607dcf7654f4dc03c073c0123edc45fcfad28132e8c60
2:cpuset:/docker/81438f4655cd771c425607dcf7654f4dc03c073c0123edc45fcfad28132e8c60
1:name=openrc:/docker
`

const containerHostPIDNamespaceCgroup = `14:name=systemd:/
13:pids:/
12:hugetlb:/
11:net_prio:/
10:perf_event:/
9:net_cls:/
8:freezer:/
7:devices:/
6:memory:/
5:blkio:/
4:cpuacct:/
3:cpu:/
2:cpuset:/
1:name=openrc:/
`

const lxcCgroup = `9:hugetlb:/lxc/81438f4655cd771c425607dcf7654f4dc03c073c0123edc45fcfad28132e8c60
8:perf_event:/lxc/81438f4655cd771c425607dcf7654f4dc03c073c0123edc45fcfad28132e8c60
7:blkio:/lxc/81438f4655cd771c425607dcf7654f4dc03c073c0123edc45fcfad28132e8c60
6:freezer:/lxc/81438f4655cd771c425607dcf7654f4dc03c073c0123edc45fcfad28132e8c60
5:devices:/lxc/81438f4655cd771c425607dcf7654f4dc03c073c0123edc45fcfad28132e8c60
4:memory:/lxc/81438f4655cd771c425607dcf7654f4dc03c073c0123edc45fcfad28132e8c60
3:cpuacct:/lxc/81438f4655cd771c425607dcf7654f4dc03c073c0123edc45fcfad28132e8c60
2:cpu:/lxc/81438f4655cd771c425607dcf7654f4dc03c073c0123edc45fcfad28132e8c60
1:cpuset:/lxc/81438f4655cd771c425607dcf7654f4dc03c073c0123edc45fcfad28132e8c60`

const emptyCgroup = ``

func TestIsContainerized(t *testing.T) {
	containerized, err := isContainerizedCgroup([]byte(nonContainerizedCgroup))
	if err != nil {
		t.Fatal(err)
	}
	assert.False(t, containerized)

	containerized, err = isContainerizedCgroup([]byte(containerCgroup))
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, containerized)

	containerized, err = isContainerizedCgroup([]byte(containerHostPIDNamespaceCgroup))
	if err != nil {
		t.Fatal(err)
	}
	assert.False(t, containerized)

	containerized, err = isContainerizedCgroup([]byte(lxcCgroup))
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, containerized)

	containerized, err = isContainerizedCgroup([]byte(emptyCgroup))
	if err != nil {
		t.Fatal(err)
	}
	assert.False(t, containerized)
}
