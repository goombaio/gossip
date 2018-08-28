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

var (
	// Version is the current version number using the semver standard.
	Version string

	// Build is the current build id represented by the last commit id.
	Build string

	helpFlag    bool
	versionFlag bool
)

// Usage ...
func Usage() {
	fmt.Println("Usage: gossiper [flags] <peers>")
	fmt.Println()
	fmt.Println("The peers argument list follows the format address:port, each item separated by a space.")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Printf("  -host\tHostname or IP, defaults to %s.\n", gossiper.DefaultHost)
	fmt.Printf("  -port\tPort number, defaults to %d.\n", gossiper.DefaultPort)
	fmt.Printf("  -help\tShow this help message.\n")
	fmt.Println()
	fmt.Println("Example:")
	fmt.Println("  gossiper -port 60900 127.0.0.1:60901")
	fmt.Println("  gossiper -port 60901 127.0.0.1:60900 127.0.0.1:60902")
	fmt.Println("  gossiper -port 60902 127.0.0.1:60903")
	fmt.Println("  gossiper -port 60903 127.0.0.1:60902")
	fmt.Println()
}

func main() {
	// Options / Configuration
	options := &gossiper.Options{
		TimestampDelay:  gossiper.DefaultTimestampDelay,
		SimulationDelay: gossiper.DefaultSimulationDelay,
		RetryDelay:      gossiper.DefaultRetryDelay,
		RetryAttempts:   gossiper.DefaultRetryAttempts,
		MaxDisplay:      gossiper.DefaultMaxDisplay,
	}

	// Flags
	flag.IntVar(&options.Port, "port", gossiper.DefaultPort, "Port number.")
	flag.StringVar(&options.Host, "host", gossiper.DefaultHost, "Hostname or IP.")
	flag.IntVar(&options.RetryDelay, "retry-delay", gossiper.DefaultRetryDelay, "Retry delay in seconds.")
	flag.IntVar(&options.SimulationDelay, "simulation-delay", gossiper.DefaultSimulationDelay, "Simulation delay in seconds.")
	flag.IntVar(&options.TimestampDelay, "timestamp-delay", gossiper.DefaultTimestampDelay, "Timestamp delay in seconds.")

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

	if len(flag.Args()) == 0 {
		fmt.Println("ERROR: Not enough peers indicated")
		Usage()
		os.Exit(1)
	}

	g := gossiper.NewGossiper(options)

	for _, peerAddress := range flag.Args() {
		g.AddPeer(peerAddress)
	}

	g.Start()
}

// ShowVersionInfo prints version and build information
func ShowVersionInfo(version, build string) {
	tpl := "gossiper version %s build %s"
	output := fmt.Sprintf(tpl, version, build)
	fmt.Println(output)
}
