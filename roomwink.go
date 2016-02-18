// Copyright (c) 2016 IBM Corp. All rights reserved.
// Use of this source code is governed by the Apache License,
// Version 2.0, a copy of which can be found in the LICENSE file.

// The /wink room command
package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
)

type WinkResponse struct {
	Rtype   string            `json:"type,omitempty"`
	Content map[string]string `json:"content,omitempty"`
}

func wink(conn *websocket.Conn, req *GameonRequest, tail, room string) error {
	var resp WinkResponse
	resp.Rtype = "event"
	resp.Content = make(map[string]string)
	resp.Content[req.UserId] = fmt.Sprintf("%s winks at you. Slyly.", MyRooms[room])
	j, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		return err
	}
	return SendMessage(conn, req.UserId, j, MTPlayer)
}
