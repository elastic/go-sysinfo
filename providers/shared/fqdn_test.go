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
	"fmt"
	"testing"
	"time"

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
		// being available.  If it starts to fail often enough
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
			expectedFQDN:     "elastic.co",
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
