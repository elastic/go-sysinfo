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
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestHost_hostname(t *testing.T) {
	testCases := []struct {
		name      string
		want      string
		lowercase bool
	}{
		{
			name:      "hostname",
			want:      "hostName",
			lowercase: false,
		},
		{
			name:      "lowercase hostname",
			want:      "hostname",
			lowercase: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			host, err := newLinuxSystem("", tc.lowercase).Host()
			if err != nil {
				t.Fatalf("could not get host info: %v\n", err)
			}

			got := host.Info()
			if got.Hostname != tc.want {
				t.Errorf("got hostname %q; want hostname %q",
					got.Hostname, tc.want)
			}
		})
	}
}
func TestHost_FQDN_set(t *testing.T) {
	host, err := newLinuxSystem("", false).Host()
	if err != nil {
		t.Fatal(fmt.Errorf("could not get host information: %w", err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	gotFQDN, err := host.FQDNWithContext(ctx)
	require.NoError(t, err)
	if gotFQDN != wantFQDN {
		t.Errorf("got FQDN %q, want: %q", gotFQDN, wantFQDN)
	}
}

func TestHost_FQDN_not_set(t *testing.T) {
	host, err := newLinuxSystem("", false).Host()
	if err != nil {
		t.Fatal(fmt.Errorf("could not get host information: %w", err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	gotFQDN, err := host.FQDNWithContext(ctx)
	require.NoError(t, err)
	hostname := host.Info().Hostname
	if gotFQDN != hostname {
		t.Errorf("name and FQDN should be the same but hostname: %s, FQDN %s", hostname, gotFQDN)
	}
}
