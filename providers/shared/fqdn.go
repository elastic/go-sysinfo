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

//go:build linux || darwin

package shared

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"strings"
)

const etcHosts = "/etc/hosts"

func FQDN() (string, error) {
	f, err := os.Open(etcHosts)
	if err != nil {
		return "", fmt.Errorf("could open %q to get FQDN: %w", etcHosts, err)
	}

	hname, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("could get hostname to look for FQDN: %w", err)
	}

	fqdn, err := fqdnFromHosts(hname, f)
	if err != nil {
		return "", fmt.Errorf("error when looking for FQDN on %s: %w", etcHosts, err)
	}

	if fqdn == "" {
		// FQDN not found on hosts file, fall back to net.Lookup?
		// add an error?
	}

	return fqdn, nil
}

// fqdnFromHosts looks for the FQDN for hostname on hostFile.
// If successfully it returns FQDN, nil. If no FQDN for hostname is found
// it returns "", nil. It returns "", err if any error happens.
func fqdnFromHosts(hostname string, hostsFile fs.File) (string, error) {
	s := bufio.NewScanner(hostsFile)

	for s.Scan() {
		fqdn := findInHostsLine(hostname, s.Text())
		if fqdn != "" {
			return fqdn, nil
		}
	}
	if err := s.Err(); err != nil {
		return "", fmt.Errorf("error reading hosts file lines: %w", err)
	}

	return "", nil
}

// findInHostsLine takes a HOSTS(5) line and searches for an alias matching
// hostname, if found it returns the canonical_hostname. The canonical_hostname
// should be the FQDN, see HOSTNAME(1).
// TODO: check k8s: https://kubernetes.io/docs/tasks/network/customize-hosts-file-for-pods/
func findInHostsLine(hostname, hostsEntry string) string {
	line, _, _ := strings.Cut(hostsEntry, "#")
	if len(line) < 1 {
		fmt.Printf("skip comment or empty: %q\n", hostsEntry)
		return ""
	}

	fileds := strings.FieldsFunc(line, func(r rune) bool {
		return r == ' ' || r == '\t'
	})

	if len(fileds) < 2 {
		// invalid hostsEntry
		return ""
	}

	// fields[0] is the ip address
	cannonical, aliases := fileds[1], fileds[1:]

	// TODO: confirm: a name should not repeat on different addresses.
	if len(fileds) == 2 {
		if fileds[1] == hostname {
			return cannonical
		}

		// If hostname was not set as an alias for FQDN, but the fist name
		// before the dot is the hostname:
		//   192.168.1.10    foo.mydomain.org	#  foo
		if hname, _, _ := strings.Cut(cannonical, "."); hname == hostname {
			return cannonical
		}
	}

	for _, h := range aliases {
		if h == hostname {
			return cannonical
		}
	}

	return ""
}
