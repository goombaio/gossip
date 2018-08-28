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
	"net"

	"github.com/google/uuid"
)

// Peer ...
type Peer struct {
	ID      uuid.UUID
	Address net.Addr
}

// NewPeer ...
func NewPeer() *Peer {
	p := &Peer{
		ID: uuid.New(),
	}
	return p
}

// String implements stringer interface
func (p *Peer) String() string {
	str := fmt.Sprintf("%s,%s://%s", p.ID, p.Address.Network(), p.Address.String())
	return str
}
