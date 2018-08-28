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

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/repejota/gossiper"
)

const (
	// DefaultPort ...
	DefaultPort = 30480

	// DefaultTimestampDelay ...
	DefaultTimestampDelay = 300

	// DefaultSimulationDelay ...
	DefaultSimulationDelay = 60

	// DefaultRetryDelaty ...
	DefaultRetryDelaty = 120

	// DefaultRetryAttempts ...
	DefaultRetryAttempts = 10

	// DefaultMaxDisplay ...
	DefaultMaxDisplay = 40
)

var (
	// Version is the current version number using the semver standard.
	Version string

	// Build is the current build id represented by the last commit id.
	Build string

	port            int
	retryDelay      int
	simulationDelay int
	timestampDelay  int
	helpFlag        bool
	versionFlag     bool
)

// Usage ...
func Usage() {
	fmt.Println("Usage: gossiper [flags] [args]")
}

func main() {
	timestampDelay = DefaultTimestampDelay
	simulationDelay = DefaultSimulationDelay
	retryDelay = DefaultRetryDelaty

	// Flags

	flag.IntVar(&port, "port", DefaultPort, "Port number.")
	flag.IntVar(&retryDelay, "retry-delay", DefaultRetryDelaty, "Retry delay in seconds.")
	flag.IntVar(&simulationDelay, "simulation-delay", DefaultSimulationDelay, "Simulation delay in seconds.")
	flag.IntVar(&timestampDelay, "timestamp-delay", DefaultTimestampDelay, "Timestamp delay in seconds.")
	flag.BoolVar(&helpFlag, "help", false, "Show usage informnation.")
	flag.BoolVar(&versionFlag, "version", false, "Show version informnation.")

	flag.Parse()

	// --help
	if helpFlag {
		Usage()
		os.Exit(0)
	}

	// --version
	if versionFlag {
		ShowVersionInfo(Version, Build)
		os.Exit(0)
	}

	// Args
	// - Build peers list

	g := gossiper.NewGossiper()

	g.Start()
}

// ShowVersionInfo prints version and build information
func ShowVersionInfo(version, build string) {
	tpl := "gossiper version %s build %s"
	output := fmt.Sprintf(tpl, version, build)
	fmt.Println(output)
}
