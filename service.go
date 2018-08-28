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
	"net"
	"time"

	"github.com/google/uuid"
)

// Service ...
type Service struct {
	options  *Options
	listener *net.UDPConn
	ID       uuid.UUID
	ticker   *time.Ticker
}

// NewService ...
func NewService(options *Options) *Service {
	s := &Service{
		options: options,
		ID:      uuid.New(),
	}
	return s
}

// Start ...
func (s *Service) Start() {
	log.Printf("Service:\n")
	log.Printf("  - ID: %s - Address: ip:%d\n", s.ID, s.options.Port)

	// Ticker
	d := time.Duration(s.options.TimestampDelay) * time.Millisecond
	s.ticker = time.NewTicker(d)
	go func() {
		for t := range s.ticker.C {
			log.Printf("Tick at %s\n", t)
		}
	}()

	// Server listener ( UDP sockets transport )
	serverAddress, err := net.ResolveUDPAddr("udp", s.options.Address())
	if err != nil {
		log.Fatal(err)
	}
	s.listener, err = net.ListenUDP("udp", serverAddress)
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 1024)
	go func() {
		for {
			n, addr, err := s.listener.ReadFromUDP(buf)
			if n > 0 {
				message := string(buf[0 : n-1])
				log.Printf("Message received %q from %s", message, addr)
				if err != nil {
					fmt.Println("Error: ", err)
				}
			}
		}
	}()
}

// Stop ...
func (s *Service) Stop() {
	log.Printf("Stopping service: %s\n", s.ID)
	s.ticker.Stop()
	s.listener.Close()
}
