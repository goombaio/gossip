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

// Package gossiper implements a gossip network protocol. It is
// intended only for fun and learning.
//
// Gossip protocols are network protocols where each machine, each
// *peer*, does not have a complete list of all peers. Instead, it knows
// only a subset of them. In order to spread a message to all the peers,
// every "gossiper" transmits the message to all the peers it knows, in
// turn, each transmits it to all the peers it knows and so on. As long
// as the set of peers is connected, the message will eventually reach
// everyone.
//
// The best-known examples of gossip protocols are the Network News (RFC
// 1036) and BGP (RFC 4271).
//
// An important part of a gossip protocol is the history: peers must
// remember which messages they sent, to avoid wasting time (or, worse,
// creating endless loops) with peers which already know the message.
//
// Note that a successful gossip protocol does not require every pair of
// peers to communicate by the same means (Network News is a good
// example: not everyone uses NNTP). But, in this simple example, the
// protocol between two peers is fixed. Every peer has an ID, set at
// startup. The "server" (the peer which replied to the connection) sends
// its ID followed by a comma. The "client" (the peer which initiated the
// connection) sends its ID followed by a comma and by the message (one
// line only). In this implementation, for each peer, only a tuple (IP
// address, port) is used to connect but it is not imposed by the
// protocol (machines are identified by the ID, not by the IP address).
//
// Peers remember the messages they have seen (in the global history) and
// the messages they sent to each peer (in a per-peer history).
package gossiper
