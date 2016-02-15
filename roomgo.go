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
func exitRoom(conn *websocket.Conn, req *GameonRequest, tail, room string) error {
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
	case "n":
		banter = "Going North!"
	case "s":
		banter = "Going south! Later, Gator!!"
	case "e":
		banter = "Going east!"
	case "w":
		banter = "Going west, we think."
	default:
		checkpoint(locus, "UNKNOWN.DIRECTION")
		validExit = false
		banter = fmt.Sprintf("'%s'?!? There is no exit with that name. Try again.", dir)
	}

	var cresp ChatResponse
	cresp.Rtype = "chat"
	cresp.Username = req.Username
	cresp.Content = banter
	j, err := json.MarshalIndent(cresp, "", "    ")
	if err != nil {
		return err
	}
	err = sendPlayerResp(conn, req.UserId, j)
	if err != nil {
		return err
	}
	if validExit {
		j, err := json.MarshalIndent(lresp, "", "    ")
		if err != nil {
			return err
		}
		err = sendResp(conn, req.UserId, j, MTPlayerLocation)
	}
	return err
}
