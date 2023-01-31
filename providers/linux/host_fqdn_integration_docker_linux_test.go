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

//go:build integration && docker

package linux

import (
	"fmt"
	"testing"
)

func TestHost_FQDN_set(t *testing.T) {
	host, err := newLinuxSystem("").Host()
	if err != nil {
		t.Fatal(fmt.Errorf("could not het host information: %w", err))
	}

	got := host.Info()
	if got.FQDN != wantFQDN {
		t.Errorf("got FQDN %q, want: %q", got.FQDN, wantFQDN)
	}
}

func TestHost_FQDN_not_set(t *testing.T) {
	host, err := newLinuxSystem("").Host()
	if err != nil {
		t.Fatal(fmt.Errorf("could not het host information: %w", err))
	}

	got := host.Info()
	if got.Hostname != got.FQDN {
		t.Errorf("name and FQDN should be the same but hostname: %s, FQDN %s", got.Hostname, got.FQDN)
	}
}

// ❯ docker run --hostname myhost.co --rm -v "$PWD":/usr/src/elastic/go-sysinfo -it -w /usr/src/elastic/go-sysinfo golang:1.19 /bin/bash
// root@myhost:/usr/src/elastic/go-sysinfo# hostname
// myhost.co
// root@myhost:/usr/src/elastic/go-sysinfo# hostname -s
// myhost
// root@myhost:/usr/src/elastic/go-sysinfo# hostname -f
// myhost.co
// root@myhost:/usr/src/elastic/go-sysinfo# hostname -d
// co
// root@myhost:/usr/src/elastic/go-sysinfo# cat /proc/sys/kernel/hostname
// myhost.co
// root@myhost:/usr/src/elastic/go-sysinfo# cat /proc/sys/kernel/domainname
// (none)
// root@myhost:/usr/src/elastic/go-sysinfo# cat /etc/hosts
// 127.0.0.1	localhost
// ::1	localhost ip6-localhost ip6-loopback
// fe00::0	ip6-localnet
// ff00::0	ip6-mcastprefix
// ff02::1	ip6-allnodes
// ff02::2	ip6-allrouters
// 172.17.0.2	myhost.co myhost
// root@myhost:/usr/src/elastic/go-sysinfo#

// ❯ docker run --hostname myhost --domainname co --rm -v "$PWD":/usr/src/elastic/go-sysinfo -it -w /usr/src/elastic/go-sysinfo golang:1.19 /bin/bash
// root@myhost:/usr/src/elastic/go-sysinfo# hostname
// myhost
// root@myhost:/usr/src/elastic/go-sysinfo# hostname -f
// myhost.co
// root@myhost:/usr/src/elastic/go-sysinfo# hostname -d
// co
// root@myhost:/usr/src/elastic/go-sysinfo# cat /proc/sys/kernel/hostname
// myhost
// root@myhost:/usr/src/elastic/go-sysinfo# cat /proc/sys/kernel/domainname
// co
// root@myhost:/usr/src/elastic/go-sysinfo# cat /etc/hosts
// 127.0.0.1	localhost
// ::1	localhost ip6-localhost ip6-loopback
// fe00::0	ip6-localnet
// ff00::0	ip6-mcastprefix
// ff02::1	ip6-allnodes
// ff02::2	ip6-allrouters
// 172.17.0.2	myhost.co myhost
