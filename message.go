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

type PlayerMessage struct {
	Rtype    string            `json:"type,omitempty"`
	Content  map[string]string `json:"content,omitempty"`
	Bookmark int               `json:"bookmark,omitempty"`
}

var bookmark = 1

// Sends an event message to a player using the current websocket.
func sendMessageToPlayer(conn *websocket.Conn, mUser, uid string) (e error) {
	//TODO//cleanup//locus := "SEND.MSG-TO-PLAYER"
	var msg PlayerMessage
	var j []byte
	msg.Rtype = "event"
	msg.Bookmark = bookmark
	bookmark += 1
	msg.Content = make(map[string]string)

	msg.Content[uid] = mUser

	j, e = json.MarshalIndent(msg, "", "    ")
	if e != nil {
		return
	}
	e = sendMsg(conn, uid, j, MTPlayer)
	//TODO cleanup
	//m := fmt.Sprintf("%s,%s,%s", MTPlayer, target, string(j))
	//e = conn.WriteMessage(ExpectedMessageType, []byte(m))
	//if config.debug {
	//	checkpoint(locus, fmt.Sprintf("MSG=%s", m))
	//}
	//if e != nil {
	//	checkpoint(locus, fmt.Sprintf("FAILED err=%s", e.Error()))
	//}
	return
}

// Sends a message with a JSON payload.
func sendMsg(conn *websocket.Conn, targetid string, j []byte, messageType string) (e error) {
	locus := "SEND.MSG"
	var m = fmt.Sprintf("%s,%s,%s", messageType, targetid, string(j))
	e = conn.WriteMessage(ExpectedMessageType, []byte(m))
	if config.debug {
		checkpoint(locus, fmt.Sprintf("m=%s", m))
	}
	if e != nil {
		checkpoint(locus, fmt.Sprintf("FAILED err=%s", e.Error()))
	}
	return
}
