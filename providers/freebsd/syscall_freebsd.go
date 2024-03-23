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

//go:build freebsd && cgo
// +build freebsd,cgo

package freebsd

//#include <sys/sysctl.h>
//#include <stdlib.h>
import "C"

import (
	"bytes"
	"encoding/binary"
	"sync"
	"unsafe"
)

// Buffer Pool

var bufferPool = sync.Pool{
	New: func() interface{} {
		return &poolMem{
			buf: make([]byte, 512),
		}
	},
}

type poolMem struct {
	buf  []byte
	pool *sync.Pool
}

func getPoolMem() *poolMem {
	pm := bufferPool.Get().(*poolMem)
	pm.buf = pm.buf[0:cap(pm.buf)]
	pm.pool = &bufferPool
	return pm
}

func (m *poolMem) Release() { m.pool.Put(m) }

func sysctlbyname(name string, value interface{}) error {
	mem := getPoolMem()
	defer mem.Release()

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	size := C.ulong(len(mem.buf))
	if n, err := C.sysctlbyname(cname, unsafe.Pointer(&mem.buf[0]), &size, nil, C.ulong(0)); n != 0 {
		return err
	}

	data := mem.buf[0:size]

	switch v := value.(type) {
	case *[]byte:
		out := make([]byte, len(data))
		copy(out, data)
		*v = out
		return nil
	default:
		return binary.Read(bytes.NewReader(data), binary.LittleEndian, v)
	}
}

func sysctlByName(name string, out interface{}) error {
	return sysctlbyname(name, out)
}
