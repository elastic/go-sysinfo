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

//go:build freebsd

package freebsd

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elastic/go-sysinfo/types"
)

func TestOperatingSystem(t *testing.T) {
	t.Run("freebsd14", func(t *testing.T) {
		os, err := getOSInfo("")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, types.OSInfo{
			Type:     "",
			Family:   "freebsd",
			Platform: "freebsd",
			Name:     "FreeBSD",
			Version:  "14.0-RELEASE",
			Major:    14,
			Minor:    0,
			Patch:    0,
		}, *os)
		t.Logf("%#v", os)
	})
}
