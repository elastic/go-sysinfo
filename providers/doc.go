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

// Package providers
//
// # Hostname Behavior
//
// Starting from version v1.11.0, the host provider started automatically
// lowercasing hostnames. This behavior was reverted in v1.14.1.
//
// To provide flexibility and allow users to control this behavior, the
// `LowercaseHostname` and `SetLowerHostname` functions were added.
//
// By default, hostnames are not lowercased. If you require hostnames to be
// lowercased, explicitly set this using `SetLowerHostname(true)`.
package providers
