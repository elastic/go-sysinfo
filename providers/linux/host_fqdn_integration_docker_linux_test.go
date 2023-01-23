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

func TestHost_FQDN_Domain_Cgo(t *testing.T) {
	host, err := newLinuxSystem("").Host()
	if err != nil {
		t.Fatal(fmt.Errorf("could not het host information: %w", err))
	}

	got := host.Info()
	if got.Hostname != wantHostname {
		t.Errorf("got wrong hostname want: %q, got %q", wantHostname, got.Hostname)
	}
	if got.Domain != wantDomainCgo {
		t.Errorf("got wrong domain want: %q, got %q", wantDomainCgo, got.Domain)
	}
	if got.FQDN != fmt.Sprintf("%s.%s", wantHostname, wantDomainCgo) {
		t.Errorf("FQDN shpould not be empty")
	}
}

func TestHost_FQDN_No_Domain_Cgo(t *testing.T) {
	host, err := newLinuxSystem("").Host()
	if err != nil {
		t.Fatal(fmt.Errorf("could not het host information: %w", err))
	}

	got := host.Info()
	if got.Hostname != wantHostname {
		t.Errorf("got wrong hostname want: %s, got %s", wantHostname, got.Hostname)
	}
	if got.Domain != "" {
		t.Errorf("got wrong domain should be empty but got %s", got.Domain)
	}
	wantFQDN := fmt.Sprintf("%s.%s", wantHostname, "lan")
	if got.FQDN != wantFQDN {
		t.Errorf("got wrong FQDN, want: %s, got %s", wantFQDN, got.FQDN)
	}
}

func TestHost_FQDN_Domain_NoCgo(t *testing.T) {
	host, err := newLinuxSystem("").Host()
	if err != nil {
		t.Fatal(fmt.Errorf("could not het host information: %w", err))
	}

	got := host.Info()
	if got.Hostname != wantHostname {
		t.Errorf("hostname want: %s, got %s", wantHostname, got.Hostname)
	}
	if got.Domain != "" {
		t.Errorf("domain should be empty but got %s", got.Domain)
	}
	if got.FQDN != "" {
		t.Errorf("FQDN should empty, got: %s", got.FQDN)
	}
}
