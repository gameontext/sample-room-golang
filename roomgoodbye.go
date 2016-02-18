// Copyright (c) 2016 IBM Corp. All rights reserved.
// Use of this source code is governed by the Apache License,
// Version 2.0, a copy of which can be found in the LICENSE file.

// Room Good-Bye
package main

import (
	"fmt"
	"github.com/gorilla/websocket"
)

// Handles the "good-bye" request that is received each time
// a player leaves our room.
func handleGoodbye(conn *websocket.Conn, req *GoodbyeMessage, room string) error {
	locus := "GOODBYE"
	checkpoint(locus, fmt.Sprintf("room=%s userid=%s username=%s\n",
		MyRooms[room], req.UserId, req.Username))

	removePlayer(makePlayerKey(req.UserId, room))

	// Announce to the room that the player has left.
	m := fmt.Sprintf("%s has left %s.", req.Username, config.roomName)
	broadcastMessage(room, m, "tracker", "*")
	return nil
}
