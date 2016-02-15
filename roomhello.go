// Copyright (c) 2016 IBM Corp. All rights reserved.
// Use of this source code is governed by the Apache License,
// Version 2.0, a copy of which can be found in the LICENSE file.

// Room Hello
package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
)

type HelloResponse struct {
	Rtype       string `json:"type,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`

	// We have intentially omitted the exits response field
	// because we do not wish to override our initial exit setup.

	Commands map[string]string `json:"commands,omitempty"`
}

// Handles the "hello" request that is received each time
// a player enters our room. The incoming request is a string
// with the format <room>,<json>, where the JSON string contains
// the userid and username of the player entering our room. For
// example: "Room 3100",{"username": "ebullient","userId": "github:808713"}
//
// Return an error if a problem occurs, otherwise return nil.
func handleHello(conn *websocket.Conn, req *GameonRequest, room string) (e error) {
	checkpoint("HELLO", fmt.Sprintf("room=%s userid=%s username=%s\n",
		config.roomName, req.UserId, req.Username))

	// Announce to the room and the player that the player has entered
	// the room. Ignore errors for this.
	mUser := fmt.Sprintf("Welcome to %s, %s. Take your time. Look around.",
		MyRooms[room], req.Username)
	mRoom := fmt.Sprintf("%s has entered %s.", req.Username, MyRooms[room])
	sendMessageToRoom(conn, mRoom, mUser, req.UserId)

	// Send back the required response. Do not ignore these errors.
	var resp HelloResponse
	var j []byte
	resp.Rtype = "location"
	resp.Name = config.roomName
	resp.Description = fmt.Sprintf("This is %s", MyRooms[room])

	// The /help command's output is somewhat canned.  That is, it will
	// always list a minimal set of commands that the room should respond
	// to:
	//   - /help, /exits and /sos are implemented by the Game On! server and
	//     so they will always exist and function for free.  /
	//   - /go, /examine and /look are included in the minimal command set,
	//     but the room must catch and respond to these commands in order to
	//     do anything useful.
	//
	// If your room supports addtional commands then their descriptions should
	// be added to the location response so that Game On! knows to include them
	// in the output from a /help command in your room.
	//
	// You can change the description of one of the minimal commands, by including
	// it in the response and by giving it your own descriptive text.  There is no
	// way, currently, to remove a command from the list that you are choosing not
	// to support.
	resp.Commands = make(map[string]string)
	for _, c := range commandsWeAdd {
		resp.Commands[c.cmd] = c.desc
	}
	j, e = json.MarshalIndent(resp, "", "    ")
	if e != nil {
		return
	}
	e = sendPlayerResp(conn, req.UserId, j)
	return
}
