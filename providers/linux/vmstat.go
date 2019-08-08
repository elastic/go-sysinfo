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
	"reflect"

	"github.com/elastic/go-sysinfo/types"
	"github.com/pkg/errors"
)

// parseVMStat parses the contents of /proc/vmstat
func parseVMStat(content []byte) (*types.VMStatInfo, error) {

	vmStat := &types.VMStatInfo{}
	refVal := reflect.ValueOf(vmStat).Elem()

	err := parseKeyValue(content, " ", func(key, value []byte) error {
		// turn our []byte value into an int
		val, err := parseBytesOrNumber(value)
		if err != nil {
			return errors.Wrapf(err, "failed to parse %v value of %v", string(key), string(value))
		}

		// Search The struct object to see if we have a field with a tag that matches the raw key coming off the file input
		// This is the best way I've found to "search" for a a struct field based on a struct tag value.
		// In this case, the /proc/vmstat keys are struct tags.
		fieldToSet := refVal.FieldByNameFunc(func(name string) bool {
			testField, exists := reflect.TypeOf(vmStat).Elem().FieldByName(name)
			if !exists {
				return false
			}
			if testField.Tag.Get("vmstat") == string(key) {
				return true
			}
			return false
		})

		// This protects us from fields in /proc/vmstat that we don't have added in our struct
		//This is just a way to make sure we actually found a field in the above `FieldByNameFunc`
		if fieldToSet.CanSet() {
			fieldToSet.SetUint(val)
		}
		return nil
	})

	return vmStat, err
}
