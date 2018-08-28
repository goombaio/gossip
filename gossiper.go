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
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/google/uuid"
)

// Gossiper ...
type Gossiper struct {
	options *Options
	service *Service

	signal chan os.Signal
	quit   chan bool

	mu    sync.Mutex
	peers map[uuid.UUID]*Peer
}

// NewGossiper ...
func NewGossiper(options *Options) *Gossiper {
	g := &Gossiper{
		options: options,
		signal:  make(chan os.Signal, 1),
		quit:    make(chan bool, 1),
		peers:   make(map[uuid.UUID]*Peer),
	}

	g.service = NewService(options)

	return g
}

// Start ...
func (g *Gossiper) Start() {
	log.Printf("Starting gossiper..\n")

	g.service.Start()

	// Peers
	log.Printf("%d peers:\n", len(g.peers))
	for _, peer := range g.peers {
		log.Printf("  - ID: %s - Address: %s", peer.ID, peer.Address)
	}

	// - Setup queue on each peer
	// - Setup sender on each peer
	// - Start queues an senders on each peer

	g.hanleSignals()
}

// Stop ...
func (g *Gossiper) Stop() {
	log.Printf("Stopping gossiper..\n")

	g.service.Stop()
}

// AddPeer ...
func (g *Gossiper) AddPeer(peerAddress string) {
	p := NewPeer(peerAddress)
	g.mu.Lock()
	defer g.mu.Unlock()
	g.peers[p.ID] = p
}

// hanleSignalss enables graceful shutdown listening for OS signals.
//
// listens for two signals
// SIGTERM: generic signal used to cause program termination.
// SIGINT: signal used when the user types C-c
func (g *Gossiper) hanleSignals() {
	signal.Notify(g.signal, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-g.signal
		fmt.Printf(" ==> Trap signal: %s\n", s)
		g.Stop()
		g.quit <- true
	}()
	<-g.quit
}
