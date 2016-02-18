// Copyright (c) 2016 IBM Corp. All rights reserved.
// Use of this source code is governed by the Apache License,
// Version 2.0, a copy of which can be found in the LICENSE file.

package main

// Room /chat command
import (
	"github.com/gorilla/websocket"
)

func handleChat(conn *websocket.Conn, req *GameonRequest, room string) error {
	BroadcastMessage(room, req.Content, req.Username, "*")
	return nil
}
