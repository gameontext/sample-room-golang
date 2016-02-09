// Copyright (c) 2016 IBM Corp. All rights reserved.
// Use of this source code is governed by the Apache License,
// Version 2.0, a copy of which can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
)

// Handles the "hello" request that is received each time
// a player enters our room. The incoming request is a string
// with the format <room>,<json>, where the JSON string contains
// the userid and username of the player entering our room. For
// example, "Room 3100,{\"username\": \"DevUser\",\"userId\": \"dummy.DevUser\"}"".
// Return true if successful.
func handleHello(conn *websocket.Conn, req *GameonRequest) error {
	if config.debug {
		fmt.Printf("HELLO room=%s userid=%s username=%s\n",
			config.roomName, req.UserId, req.Username)
	}
	return sayHello(conn, req)
}

type Exits struct {
	E string `json:"E,omitempty"`
	W string `json:"W,omitempty"`
	N string `json:"N,omitempty"`
	S string `json:"S,omitempty"`
}

type HelloResponse struct {
	Rtype       string `json:"type,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	OurExits    Exits  `json:exits,omitempty`
}

func sayHello(conn *websocket.Conn, req *GameonRequest) error {
	var resp HelloResponse
	resp.Rtype = "location"
	resp.Name = config.roomName
	resp.Description = "This is yet another crazy room."

	j, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		return err
	}
	return sendPlayerResp(conn, req.UserId, j)
}
