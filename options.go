// Copyright 2018, gossiper project Authors. All rights reserved.
//
// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with this
// work for additional information regarding copyright ownership.  The ASF
// licenses this file to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.  See the
// License for the specific language governing permissions and limitations
// under the License.

package gossiper

import "fmt"

const (
	// DefaultPort ...
	DefaultPort = 90900

	// DefaultHost ...
	DefaultHost = "127.0.0.1"

	// DefaultTimestampDelay (miliseconds) ...
	DefaultTimestampDelay = 5000

	// DefaultSimulationDelay (miliseconds) ...
	DefaultSimulationDelay = 1000

	// DefaultRetryDelay (miliseconds) ...
	DefaultRetryDelay = 2000

	// DefaultRetryAttempts ...
	DefaultRetryAttempts = 10

	// DefaultMaxDisplay ...
	DefaultMaxDisplay = 40
)

// Options ...
type Options struct {
	Port            int
	Host            string
	TimestampDelay  int
	SimulationDelay int
	RetryDelay      int
	RetryAttempts   int
	MaxDisplay      int
}

// Address ...
func (o *Options) Address() string {
	a := fmt.Sprintf("%s:%d", o.Host, o.Port)
	return a
}
