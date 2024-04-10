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

package freebsd

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

func parseKeyValue(content []byte, separator string, callback func(key, value []byte) error) error {
	sc := bufio.NewScanner(bytes.NewReader(content))
	for sc.Scan() {
		parts := bytes.SplitN(sc.Bytes(), []byte(separator), 2)
		if len(parts) != 2 {
			continue
		}

		if err := callback(parts[0], bytes.TrimSpace(parts[1])); err != nil {
			return err
		}
	}

	return sc.Err()
}

func parseBytesOrNumber(data []byte) (uint64, error) {
	parts := bytes.Fields(data)

	if len(parts) == 0 {
		return 0, errors.New("empty value")
	}

	num, err := strconv.ParseUint(string(parts[0]), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse value: %w", err)
	}

	var multiplier uint64 = 1
	if len(parts) >= 2 {
		switch string(parts[1]) {
		case "kB":
			multiplier = 1024
		default:
			return 0, fmt.Errorf("unhandled unit %v", string(parts[1]))
		}
	}

	return num * multiplier, nil
}
