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
	"time"

	"github.com/google/uuid"
)

// Service ...
type Service struct {
	ID     uuid.UUID
	ticker *time.Ticker
}

// NewService ...
func NewService() *Service {
	s := &Service{
		ID: uuid.New(),
	}
	return s
}

// Start ...
func (s *Service) Start() {
	log.Printf("Starting service: %s\n", s.ID)

	// TODO:
	// - Tick duration should be setup by configuration/flags
	// - Always seconds or microseconds? ( it is a duration )
	s.ticker = time.NewTicker(time.Second * 30)
	go func() {
		for t := range s.ticker.C {
			log.Printf("Tick at %s\n", t)
		}
	}()
}

// Stop ...
func (s *Service) Stop() {
	log.Printf("Stopping service: %s\n", s.ID)
	s.ticker.Stop()
}
