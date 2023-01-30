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

//go:build cgo

package linux

// #include <unistd.h>
// #include <stdlib.h>
import "C"
import (
	"fmt"
	"os"
	"unsafe"
)

func fqdn() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}

	domain, err := domainname()
	if err != nil {
		return "", err
	}

	if domain == "" || domain == "(none)" { // mimicking 'hostname -f' behaviour
		domain = "lan"
	}

	return fmt.Sprintf("%s.%s", hostname, domain), nil
}

func hostname() (string, error) {
	const buffSize = 64
	buff := make([]byte, buffSize)
	size := C.size_t(buffSize)
	cString := C.CString(string(buff))
	defer C.free(unsafe.Pointer(cString))

	_, errno := C.gethostname(cString, size)
	if errno != nil {
		return "", fmt.Errorf("cgo call gethostname errored: %v", errno)
	}

	var name string = C.GoString(cString)

	if name == "(none)" {
		name = ""
	}

	return name, nil
}

func domainname() (string, error) {
	const buffSize = 64
	buff := make([]byte, buffSize)
	size := C.size_t(buffSize)
	cString := C.CString(string(buff))
	defer C.free(unsafe.Pointer(cString))

	_, errno := C.getdomainname(cString, size)
	if errno != nil {
		return "", fmt.Errorf("syscall getdomainname errored: %v", errno)
	}

	var domain string = C.GoString(cString)

	if domain == "(none)" { // mimicking 'hostname -f' behaviour
		domain = ""
	}

	return domain, nil
}
