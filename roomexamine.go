// Copyright (c) 2016 IBM Corp. All rights reserved.
// Use of this source code is governed by the Apache License,
// Version 2.0, a copy of which can be found in the LICENSE file.

package main

// Room /examine command
import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"strings"
)

type ExaminationResponse struct {
	Rtype    string            `json:"type,omitempty"`
	ExitId   string            `json:"exitId,omitempty"`
	Content  map[string]string `json:"content,omitempty"`
	Bookmark int               `json:"bookmark,omitempty"`
}

func examineObject(conn *websocket.Conn, req *GameonRequest, tail, room string) error {
	var resp ExaminationResponse
	resp.Rtype = "event"
	resp.Content = make(map[string]string)
	obj := strings.Trim(tail, " ")
	if len(obj) > 0 {
		var verb string
		if "s" == strings.ToLower(obj[len(obj)-1:]) {
			verb = "are"
		} else {
			verb = "is"
		}
		resp.Content[req.UserId] = fmt.Sprintf("There %s no %s here in %s. Keep moving.",
			verb, obj, MyRooms[room])
	} else {
		resp.Content[req.UserId] = fmt.Sprintf("There is nothing here in %s. Keep moving.",
			MyRooms[room])
	}

	j, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		return err
	}
	return SendMessage(conn, req.UserId, j, MTPlayer)
}
