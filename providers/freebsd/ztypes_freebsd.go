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

// Code generated by cmd/cgo -godefs; and then patched up to fix
// an alignment issue
// cgo -godefs defs_freebsd.go

package freebsd

type vmTotal struct {
	Rq     int16
	Dw     int16
	Pw     int16
	Sl     int16
	_      int16 // cgo doesn't generate the same alignment as C does
	Sw     int16
	Vm     int32
	Avm    int32
	Rm     int32
	Arm    int32
	Vmshr  int32
	Avmshr int32
	Rmshr  int32
	Armshr int32
	Free   int32
}

type kvmSwap struct {
	Devname   [32]int8
	Used      uint32
	Total     uint32
	Flags     int32
	Reserved1 uint32
	Reserved2 uint32
}

type clockInfo struct {
	Hz     int32
	Tick   int32
	Spare  int32
	Stathz int32
	Profhz int32
}
