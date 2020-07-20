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

// +build ignore

package aix

/*
#include <sys/types.h>
#include <utmp.h>
#include <sys/procfs.h>
*/
import "C"

type prcred C.prcred_t

type pstatus C.pstatus_t
type prTimestruc64 C.pr_timestruc64_t
type prSigset C.pr_sigset_t
type fltset C.fltset_t
type lwpstatus C.lwpstatus_t
type prSiginfo64 C.pr_siginfo64_t
type prStack64 C.pr_stack64_t
type prSigaction64 C.struct_pr_sigaction64
type prgregset C.prgregset_t
type prfpregset C.prfpregset_t
type pfamily C.pfamily_t
