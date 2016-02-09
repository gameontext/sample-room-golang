// Copyright (c) 2016 IBM Corp. All rights reserved.
// Use of this source code is governed by the Apache License,
// Version 2.0, a copy of which can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

type ChatResponse struct {
	Rtype    string `json:"type,omitempty"`
	Username string `json:"username,omitempty"`
	Content  string `json:"content,omitempty"`
	Bookmark int    `json:"bookmark,omitempty"`
}

func handleChat(conn *websocket.Conn, req *GameonRequest) error {
	var resp ChatResponse
	resp.Rtype = "chat"
	resp.Username = req.Username
	resp.Content = req.Content
	resp.Bookmark = 0
	j, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		return err
	}
	return sendPlayerResp(conn, "*", j)
}
