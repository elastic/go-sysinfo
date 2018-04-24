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

// +build darwin,amd64,cgo

package darwin

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elastic/go-sysinfo/internal/registry"
)

var _ registry.HostProvider = darwinSystem{}
var _ registry.ProcessProvider = darwinSystem{}

func TestKernProcInfo(t *testing.T) {
	var p process
	if err := kern_procargs(os.Getpid(), &p); err != nil {
		t.Fatal(err)
	}

	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, exe, p.exe)
	assert.Equal(t, os.Args, p.args)
}
