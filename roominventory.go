// Copyright (c) 2016 IBM Corp. All rights reserved.
// Use of this source code is governed by the Apache License,
// Version 2.0, a copy of which can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"time"
)

type InventoryResponse struct {
	Rtype    string            `json:"type,omitempty"`
	ExitId   string            `json:"exitId,omitempty"`
	Content  map[string]string `json:"content,omitempty"`
	Bookmark int               `json:"bookmark,omitempty"`
}

// TimedText allows us to specify a delay, in milliseconds,
// before the message containing the string is transmitted.
type TimedText struct {
	msPause int
	s       string
}

// Note that GameOn will strip white space from the ends,
// so adding spaces for indentation will not currently work.
var cheekyInventoryRemarks = []TimedText{
	{0, "Riddle me this."},
	{750, "\"How many pockets could a pickpocket pick"},
	{0, "if a pickpocket could pick pockets?\""},
	{2000, "(Enough, apparently. Your pockets are now empty.)"},
}

func checkInventory(conn *websocket.Conn, req *GameonRequest, tail, room string) error {
	for _, tt := range cheekyInventoryRemarks {
		var resp ExaminationResponse
		resp.Rtype = "event"
		resp.Content = make(map[string]string)
		resp.Content[req.UserId] = tt.s
		j, err := json.MarshalIndent(resp, "", "    ")
		if err != nil {
			return err
		}
		if tt.msPause > 0 {
			time.Sleep(time.Duration(tt.msPause) * time.Millisecond)
		}
		err = sendResp(conn, req.UserId, j, MTPlayer)
		if err != nil {
			return err
		}
	}
	return nil
}
