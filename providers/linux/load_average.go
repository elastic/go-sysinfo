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

package linux

import (
	"bufio"
	"bytes"
	"fmt"
	"reflect"

	"github.com/elastic/go-sysinfo/types"
)

var loadAverageInfoFieldIndexToTag = make(map[int]string)

func init() {
	var loadavg types.LoadAverageInfo
	val := reflect.ValueOf(loadavg)
	typ := reflect.TypeOf(loadavg)

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		if tag := field.Tag.Get("json"); tag != "" {
			loadAverageInfoFieldIndexToTag[i] = tag
		}
	}
}

// parseLoadAvg parses the content of /proc/loadavg
func parseLoadAvg(content []byte) (*types.LoadAverageInfo, error) {
	var loadAvg types.LoadAverageInfo
	refValues := reflect.ValueOf(&loadAvg).Elem()

	s := bufio.NewScanner(bytes.NewReader(content))
	s.Split(bufio.ScanWords)

	for index := 0; index < len(loadAverageInfoFieldIndexToTag); index++ {
		s.Scan()
		data := s.Bytes()
		val, err := parseBytesOrFloat(data)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %v: %w", data, err)
		}

		sval := refValues.Field(index)

		if sval.CanSet() {
			sval.SetFloat(val)
		}
	}

	return &loadAvg, s.Err()
}
