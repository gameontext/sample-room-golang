// Copyright (c) 2016 IBM Corp. All rights reserved.
// Use of this source code is governed by the Apache License,
// Version 2.0, a copy of which can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"strings"
)

type HelpResponse struct {
	Rtype    string            `json:"type,omitempty"`
	ExitId   string            `json:"exitId,omitempty"`
	Content  map[string]string `json:"content,omitempty"`
	Bookmark int               `json:"bookmark,omitempty"`
}

func helpCommand(conn *websocket.Conn, req *GameonRequest, tail, room string) error {
	var resp HelpResponse
	m := fmt.Sprintf("%s understands the following commands: %s.",
		MyRooms[room],
		strings.Join(slashCommands, ", "))
	resp.Rtype = "event"
	resp.Content = make(map[string]string)
	resp.Content[req.UserId] = m
	j, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		return err
	}
	return sendResp(conn, req.UserId, j, MTPlayer)
}
