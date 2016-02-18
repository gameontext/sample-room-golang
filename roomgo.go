// Copyright (c) 2016 IBM Corp. All rights reserved.
// Use of this source code is governed by the Apache License,
// Version 2.0, a copy of which can be found in the LICENSE file.

package main

// Room /go command
import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"strings"
)

type LocationResponse struct {
	Rtype    string `json:"type,omitempty"`
	ExitId   string `json:"exitId,omitempty"`
	Content  string `json:"content,omitempty"`
	Bookmark int    `json:"bookmark,omitempty"`
}

// Exits our room if the player requests a supported exit.
func exitRoom(conn *websocket.Conn, req *GameonRequest, tail, room string) (e error) {
	locus := "EXITROOM"
	// Content must be of the form "/go direction" or "/exit direction"
	// where direction is a valid exit.
	dir := strings.ToLower(tail)
	dir = strings.Trim(dir, " ")
	checkpoint(locus, dir)
	var lresp LocationResponse
	lresp.Rtype = "exit"
	lresp.ExitId = dir
	validExit := true
	banter := ""
	switch dir {
	case "n", "north":
		banter = "Going North!"
		lresp.ExitId = "n"
	case "s", "south":
		banter = "Going south! Later, Gator!!"
		lresp.ExitId = "s"
	case "e", "east":
		banter = "Going east!"
		lresp.ExitId = "e"
	case "w", "west":
		banter = "Going west, we think."
		lresp.ExitId = "w"
	case "home":
		banter = "You can't go home again."
		validExit = false
	case "away":
		banter = "Never!"
		validExit = false
	default:
		checkpoint(locus, "UNKNOWN.DIRECTION")
		validExit = false
		banter = fmt.Sprintf("'%s'?!? There is no exit with that name. Try again.", dir)
	}

	SendMessageToPlayer(conn, banter, req.UserId)

	if validExit {
		j, err := json.MarshalIndent(lresp, "", "    ")
		if err != nil {
			return err
		}
		e = SendMessage(conn, req.UserId, j, MTPlayerLocation)
	}
	return
}
