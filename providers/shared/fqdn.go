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
	"fmt"
	"net"
	"os"
	"strings"
)

func FQDN() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("could not get hostname to look for FQDN: %w", err)
	}

	var errs error
	cname, err := net.LookupCNAME(hostname)
	if err != nil {
		errs = fmt.Errorf("could not get FQDN, all methods failed: failed looking up CNAME: %w",
			err)
	}
	if cname != "" {
		return strings.TrimSuffix(cname, "."), nil
	}

	ips, err := net.LookupIP(hostname)
	if err != nil {
		errs = fmt.Errorf("%s: failed looking up IP: %w", errs, err)
	}

	for _, ip := range ips {
		names, err := net.LookupAddr(ip.String())
		if err != nil || len(names) == 0 {
			continue
		}
		return strings.TrimSuffix(names[0], "."), nil
	}

	return "", errs
}
