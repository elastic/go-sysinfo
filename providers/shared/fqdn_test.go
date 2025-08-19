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
	"context"
	"errors"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestFQDN(t *testing.T) {
	tests := map[string]struct {
		osHostname       string
		expectedFQDN     string
		expectedErrRegex string
		timeout          time.Duration
	}{
		// This test case depends on network, particularly DNS,
		// being available. If it starts to fail often enough
		// due to occasional network/DNS unavailability, we should
		// probably just delete this test case.
		"long_real_hostname": {
			osHostname:       "elastic.co",
			expectedFQDN:     "elastic.co",
			expectedErrRegex: "",
		},
		"long_nonexistent_hostname": {
			osHostname:       "foo.bar.elastic.co",
			expectedFQDN:     "",
			expectedErrRegex: makeErrorRegex("foo.bar.elastic.co", false),
		},
		"short_nonexistent_hostname": {
			osHostname:       "foobarbaz",
			expectedFQDN:     "",
			expectedErrRegex: makeErrorRegex("foobarbaz", false),
		},
		"long_mixed_case_hostname": {
			osHostname:       "eLaSTic.co",
			expectedFQDN:     "eLaSTic.co",
			expectedErrRegex: "",
		},
		"nonexistent_timeout": {
			osHostname:       "foobarbaz",
			expectedFQDN:     "",
			expectedErrRegex: makeErrorRegex("foobarbaz", true),
			timeout:          1 * time.Millisecond,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			timeout := test.timeout
			if timeout == 0 {
				timeout = 10 * time.Second
			}

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			actualFQDN, err := fqdn(ctx, test.osHostname)
			require.Equal(t, test.expectedFQDN, actualFQDN)

			if test.expectedErrRegex == "" {
				require.Nil(t, err)
			} else {
				require.Regexp(t, test.expectedErrRegex, err.Error())
			}
		})
	}
}

func TestMockFQDN_ValidCNAME(t *testing.T) {
	defer func() {
		defaultResolver = net.DefaultResolver
	}()

	tests := map[string]struct {
		osHostname   string
		cname        string
		expectedFQDN string
	}{
		"existing_cname": {
			osHostname:   "short_hostname",
			cname:        "short_hostname.elastic.co.",
			expectedFQDN: "short_hostname.elastic.co",
		},
		"existing_cname_upper_case": {
			osHostname:   "Short_Hostname",
			cname:        "short_hostname.",
			expectedFQDN: "Short_Hostname",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()

			defaultResolver = &mockResolver{}
			defaultResolver.(*mockResolver).On("LookupCNAME", ctx, test.osHostname).Once().Return(test.cname, nil)

			actualFQDN, err := fqdn(ctx, test.osHostname)
			require.NoError(t, err)
			assert.Equal(t, test.expectedFQDN, actualFQDN)

			mock.AssertExpectationsForObjects(t, defaultResolver)
		})
	}
}

func TestMockFQDN_EmptyCNAME(t *testing.T) {
	defer func() {
		defaultResolver = net.DefaultResolver
	}()

	type ipNames struct {
		names []string
		ip    string
		err   error
	}

	tests := map[string]struct {
		osHostname   string
		ips          []net.IP
		ipsNames     []ipNames
		expectedFQDN string
	}{
		"single_ip": {
			osHostname:   "short_hostname",
			ips:          []net.IP{net.ParseIP("192.168.1.29")},
			ipsNames:     []ipNames{{ip: "192.168.1.29", names: []string{"short_hostname.elastic.co."}}},
			expectedFQDN: "short_hostname.elastic.co",
		},
		"localhost_skipped": {
			osHostname: "short_hostname",
			ips: []net.IP{
				net.ParseIP("127.0.0.1"),
				net.ParseIP("::1"),
				net.ParseIP("192.168.1.29"),
				net.ParseIP("172.1.1.2"),
			},
			ipsNames:     []ipNames{{ip: "192.168.1.29", names: []string{"short_hostname.elastic.co."}}},
			expectedFQDN: "short_hostname.elastic.co",
		},
		"skip_ips_w/o_names": {
			osHostname:   "short_hostname",
			ips:          []net.IP{net.ParseIP("192.168.1.30"), net.ParseIP("192.168.1.29"), net.ParseIP("172.1.1.2")},
			ipsNames:     []ipNames{{ip: "192.168.1.30"}, {ip: "192.168.1.29", names: []string{"short_hostname.elastic.co."}}},
			expectedFQDN: "short_hostname.elastic.co",
		},
		"lookup_errors_are_skipped": {
			osHostname: "short_hostname",
			ips:        []net.IP{net.ParseIP("192.168.1.30"), net.ParseIP("192.168.1.29"), net.ParseIP("172.1.1.2")},
			ipsNames: []ipNames{
				{ip: "192.168.1.30", names: []string{"short_hostname.elastic.co."}, err: errors.New("skipped error")},
				{ip: "192.168.1.29", names: []string{"short_hostname.elastic.co."}},
			},
			expectedFQDN: "short_hostname.elastic.co",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()

			defaultResolver = &mockResolver{}
			defaultResolver.(*mockResolver).On("LookupCNAME", ctx, test.osHostname).Once().Return("", nil)
			defaultResolver.(*mockResolver).On("LookupIP", ctx, "ip", test.osHostname).Once().Return(test.ips, nil)
			for _, ipNames := range test.ipsNames {
				defaultResolver.(*mockResolver).On("LookupAddr", ctx, ipNames.ip).Once().Return(ipNames.names, ipNames.err)
			}

			actualFQDN, err := fqdn(ctx, test.osHostname)
			require.NoError(t, err)
			assert.Equal(t, test.expectedFQDN, actualFQDN)

			mock.AssertExpectationsForObjects(t, defaultResolver)
		})
	}
}

func Test_CNAMELookupError(t *testing.T) {
	defer func() {
		defaultResolver = net.DefaultResolver
	}()

	ctx := context.Background()
	hostname := "short_hostname"
	cnameErr := errors.New("cname error")

	defaultResolver = &mockResolver{}
	// When CNAME lookup fails and LookupIP does not return any IPs, we should return an error
	defaultResolver.(*mockResolver).On("LookupCNAME", ctx, hostname).Once().Return("", cnameErr)
	defaultResolver.(*mockResolver).On("LookupIP", ctx, "ip", hostname).Once().Return([]net.IP{}, nil)

	_, err := fqdn(ctx, hostname)
	assert.ErrorIs(t, err, cnameErr)

	mock.AssertExpectationsForObjects(t, defaultResolver)
}

func Test_LookupIPError(t *testing.T) {
	defer func() {
		defaultResolver = net.DefaultResolver
	}()

	ctx := context.Background()
	hostname := "short_hostname"
	lookupIPErr := errors.New("lookup ip error")

	defaultResolver = &mockResolver{}
	// When CNAME lookup fails and LookupIP does not return any IPs, we should return an error
	defaultResolver.(*mockResolver).On("LookupCNAME", ctx, hostname).Once().Return("", nil)
	defaultResolver.(*mockResolver).On("LookupIP", ctx, "ip", hostname).Once().Return([]net.IP{}, lookupIPErr)

	_, err := fqdn(ctx, hostname)
	assert.ErrorIs(t, err, lookupIPErr)

	mock.AssertExpectationsForObjects(t, defaultResolver)
}

func makeErrorRegex(osHostname string, withTimeout bool) string {
	timeoutStr := ""
	if withTimeout {
		timeoutStr = ": i/o timeout"
	}

	return fmt.Sprintf(
		"could not get FQDN, all methods failed: "+
			"failed looking up CNAME: lookup %s.*: "+
			"failed looking up IP: lookup %s"+timeoutStr,
		osHostname,
		osHostname,
	)
}

type mockResolver struct {
	mock.Mock
}

func (m *mockResolver) LookupCNAME(ctx context.Context, host string) (string, error) {
	args := m.Called(ctx, host)
	return args.String(0), args.Error(1)
}

func (m *mockResolver) LookupIP(ctx context.Context, network, host string) ([]net.IP, error) {
	args := m.Called(ctx, network, host)
	return args.Get(0).([]net.IP), args.Error(1)
}

func (m *mockResolver) LookupAddr(ctx context.Context, addr string) ([]string, error) {
	args := m.Called(ctx, addr)
	return args.Get(0).([]string), args.Error(1)
}
