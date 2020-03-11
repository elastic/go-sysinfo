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

const systemdCgroup = `12:hugetlb:/service.slice/podc1281d63_01ab_11ea_ba0a_3cfdfe55a1c0.slice/e2b68f8a6e227921b236c686a243e8ff50f561f493d401da7ac3f8cae28f08b1
11:perf_event:/service.slice/podc1281d63_01ab_11ea_ba0a_3cfdfe55a1c0.slice/e2b68f8a6e227921b236c686a243e8ff50f561f493d401da7ac3f8cae28f08b1
10:pids:/service.slice/podc1281d63_01ab_11ea_ba0a_3cfdfe55a1c0.slice/e2b68f8a6e227921b236c686a243e8ff50f561f493d401da7ac3f8cae28f08b1
9:cpu,cpuacct:/service.slice/podc1281d63_01ab_11ea_ba0a_3cfdfe55a1c0.slice/e2b68f8a6e227921b236c686a243e8ff50f561f493d401da7ac3f8cae28f08b1
8:cpuset:/service.slice/podc1281d63_01ab_11ea_ba0a_3cfdfe55a1c0.slice/e2b68f8a6e227921b236c686a243e8ff50f561f493d401da7ac3f8cae28f08b1
7:memory:/service.slice/podc1281d63_01ab_11ea_ba0a_3cfdfe55a1c0.slice/e2b68f8a6e227921b236c686a243e8ff50f561f493d401da7ac3f8cae28f08b1
6:freezer:/service.slice/podc1281d63_01ab_11ea_ba0a_3cfdfe55a1c0.slice/e2b68f8a6e227921b236c686a243e8ff50f561f493d401da7ac3f8cae28f08b1
5:rdma:/
4:net_cls,net_prio:/service.slice/podc1281d63_01ab_11ea_ba0a_3cfdfe55a1c0.slice/e2b68f8a6e227921b236c686a243e8ff50f561f493d401da7ac3f8cae28f08b1
3:devices:/service.slice/podc1281d63_01ab_11ea_ba0a_3cfdfe55a1c0.slice/e2b68f8a6e227921b236c686a243e8ff50f561f493d401da7ac3f8cae28f08b1
2:blkio:/service.slice/podc1281d63_01ab_11ea_ba0a_3cfdfe55a1c0.slice/e2b68f8a6e227921b236c686a243e8ff50f561f493d401da7ac3f8cae28f08b1
1:name=systemd:/service.slice/podc1281d63_01ab_11ea_ba0a_3cfdfe55a1c0.slice/e2b68f8a6e227921b236c686a243e8ff50f561f493d401da7ac3f8cae28f08b1`

const emptyCgroup = ``

const kubernetesCgroup = `11:perf_event:/kubepods/burstable/podb83789a8-5f9d-11ea-bae1-0a0084deb344/9f99515d52142271cfeebef269bf4b7609b9b69b62008d6a5d316f561ccf061d
10:freezer:/kubepods/burstable/podb83789a8-5f9d-11ea-bae1-0a0084deb344/9f99515d52142271cfeebef269bf4b7609b9b69b62008d6a5d316f561ccf061d
9:hugetlb:/kubepods/burstable/podb83789a8-5f9d-11ea-bae1-0a0084deb344/9f99515d52142271cfeebef269bf4b7609b9b69b62008d6a5d316f561ccf061d
8:devices:/kubepods/burstable/podb83789a8-5f9d-11ea-bae1-0a0084deb344/9f99515d52142271cfeebef269bf4b7609b9b69b62008d6a5d316f561ccf061d
7:blkio:/kubepods/burstable/podb83789a8-5f9d-11ea-bae1-0a0084deb344/9f99515d52142271cfeebef269bf4b7609b9b69b62008d6a5d316f561ccf061d
6:cpuset:/kubepods/burstable/podb83789a8-5f9d-11ea-bae1-0a0084deb344/9f99515d52142271cfeebef269bf4b7609b9b69b62008d6a5d316f561ccf061d
5:cpu,cpuacct:/kubepods/burstable/podb83789a8-5f9d-11ea-bae1-0a0084deb344/9f99515d52142271cfeebef269bf4b7609b9b69b62008d6a5d316f561ccf061d
4:pids:/kubepods/burstable/podb83789a8-5f9d-11ea-bae1-0a0084deb344/9f99515d52142271cfeebef269bf4b7609b9b69b62008d6a5d316f561ccf061d
3:memory:/kubepods/burstable/podb83789a8-5f9d-11ea-bae1-0a0084deb344/9f99515d52142271cfeebef269bf4b7609b9b69b62008d6a5d316f561ccf061d
2:net_cls,net_prio:/kubepods/burstable/podb83789a8-5f9d-11ea-bae1-0a0084deb344/9f99515d52142271cfeebef269bf4b7609b9b69b62008d6a5d316f561ccf061d
1:name=systemd:/kubepods/burstable/podb83789a8-5f9d-11ea-bae1-0a0084deb344/9f99515d52142271cfeebef269bf4b7609b9b69b62008d6a5d316f561ccf061d
`

func TestIsContainerized(t *testing.T) {
	tests := []struct {
		cgroupStr     string
		containerized bool
	}{
		{
			cgroupStr:     nonContainerizedCgroup,
			containerized: false,
		},
		{
			cgroupStr:     containerCgroup,
			containerized: true,
		},
		{
			cgroupStr:     containerHostPIDNamespaceCgroup,
			containerized: false,
		},
		{
			cgroupStr:     lxcCgroup,
			containerized: true,
		},
		{
			cgroupStr:     systemdCgroup,
			containerized: true,
		},
		{
			cgroupStr:     emptyCgroup,
			containerized: false,
		},
		{
			cgroupStr:     kubernetesCgroup,
			containerized: true,
		},
	}

	for _, test := range tests {
		containerized, err := isContainerizedCgroup([]byte(test.cgroupStr))
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, test.containerized, containerized)
	}
}
