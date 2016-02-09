// Copyright (c) 2016 IBM Corp. All rights reserved.
// Use of this source code is governed by the Apache License,
// Version 2.0, a copy of which can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
)

// Handles the "good-bye" request that is received each time
// a player leaves our room.
func handleGoodbye(conn *websocket.Conn, req *GameonRequest) error {
	locus := "GOODBYE"
	checkpoint(locus, fmt.Sprintf("room=%s userid=%s username=%s\n",
		config.roomName, req.UserId, req.Username))
	return sayGoodbye(conn, req)
}

type Goodbye struct {
	Rtype   string            `json:"type,omitempty"`
	Content map[string]string `json:"content,omitempty"`
}

func sayGoodbye(conn *websocket.Conn, req *GameonRequest) error {
	var g Goodbye
	g.Content = make(map[string]string)
	g.Rtype = "event"
	g.Content["*"] = fmt.Sprintf("%s has left %s.", req.Username, config.roomName)
	j, err := json.MarshalIndent(g, "", "    ")
	if err != nil {
		return err
	}
	return sendPlayerResp(conn, req.UserId, j)
}
