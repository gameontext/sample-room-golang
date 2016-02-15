// Copyright (c) 2016 IBM Corp. All rights reserved.
// Use of this source code is governed by the Apache License,
// Version 2.0, a copy of which can be found in the LICENSE file.

// Common response functions
package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
)

const NoMessage = ""

type MessageResponse struct {
	Rtype    string            `json:"type,omitempty"`
	Content  map[string]string `json:"content,omitempty"`
	Bookmark int               `json:"bookmark,omitempty"`
}

var bookmark = 1

// Sends an event message to a player or, possibly, to everyone.
// (Based on same from Game On! VerySimpleRoom).
func sendMessageToRoom(conn *websocket.Conn, mRoom, mUser, uid string) (e error) {
	locus := "SEND.MSG"
	var resp MessageResponse
	var j []byte
	var target string
	resp.Rtype = "event"
	resp.Bookmark = bookmark
	bookmark += 1
	resp.Content = make(map[string]string)

	// Default to sending to user.
	target = uid
	if mRoom != NoMessage {
		resp.Content["*"] = mRoom
		// Override default and send to everyone.
		target = "*"
	}
	if mUser != NoMessage {
		resp.Content[uid] = mUser
	}
	j, e = json.MarshalIndent(resp, "", "    ")
	if e != nil {
		return
	}
	m := fmt.Sprintf("%s,%s,%s", MTPlayer, target, string(j))
	e = conn.WriteMessage(ExpectedMessageType, []byte(m))
	if config.debug {
		checkpoint(locus, fmt.Sprintf("MSG=%s", m))
	}
	if e != nil {
		checkpoint(locus, fmt.Sprintf("FAILED err=%s", e.Error()))
	}
	return
}

// Sends a player response with a JSON payload.
func sendPlayerResp(conn *websocket.Conn, targetid string, j []byte) error {
	return sendResp(conn, targetid, j, MTPlayer)
}

// Sends a response with a JSON payload.
func sendResp(conn *websocket.Conn, targetid string, j []byte, rtype string) (e error) {
	locus := "SEND.RESP"
	var m = fmt.Sprintf("%s,%s,%s", rtype, targetid, string(j))
	e = conn.WriteMessage(ExpectedMessageType, []byte(m))
	if config.debug {
		checkpoint(locus, fmt.Sprintf("MSG=%s", m))
	}
	if e != nil {
		checkpoint(locus, fmt.Sprintf("FAILED err=%s", e.Error()))
	}
	return
}
