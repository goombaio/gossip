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

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Gossiper ...
type Gossiper struct {
	signal chan os.Signal
	quit   chan bool
}

// New ...
func New() *Gossiper {
	g := &Gossiper{
		signal: make(chan os.Signal, 1),
		quit:   make(chan bool, 1),
	}
	return g
}

// Start ...
func (g *Gossiper) Start() {
	log.Printf("Starting gossiper..\n")

	// Peers
	// - Setup queue on each peer
	// - Setup sender on each peer
	// - Start queues an senders on each peer

	// Stamper

	// Server
	g.handleGracefulShutdown()
}

// Stop ...
func (g *Gossiper) Stop() {
	log.Printf("Stopping gossiper..\n")
}

// handleGracefulShutdown enables graceful shutdown
//
// listens for two signals
// SIGTERM: generic signal used to cause program termination.
// SIGINT: signal used when the user types C-c
func (g *Gossiper) handleGracefulShutdown() {
	signal.Notify(g.signal, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		_ = <-g.signal
		g.Stop()
		g.quit <- true
	}()
	<-g.quit
}
