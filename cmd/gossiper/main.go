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
)

const (
	defaultPort = 30480
)

func init() {
	var (
		helpFlag bool
	)
	flag.BoolVar(&helpFlag, "help", false, "Show usage informnation.")

	if helpFlag {
		Usage()
		os.Exit(0)
	}
}

// Usage ...
func Usage() {
	fmt.Println("Usage: gossiper [flags]")
}

func main() {
	fmt.Println("Usage: gossiper")

	port := defaultPort
	flag.IntVar(&port, "port", defaultPort, "port number.")

	/*
		myid := nil

		timestampDelay := 300
		simulationDelay := 60
		retryDelay := 120
	*/

}
