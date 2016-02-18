// Copyright (c) 2016 IBM Corp. All rights reserved.
// Use of this source code is governed by the Apache License,
// Version 2.0, a copy of which can be found in the LICENSE file.

// Room /look command
package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"time"
)

type LookResponse struct {
	Rtype    string            `json:"type,omitempty"`
	ExitId   string            `json:"exitId,omitempty"`
	Content  map[string]string `json:"content,omitempty"`
	Bookmark int               `json:"bookmark,omitempty"`
}

// Note that GameOn will strip white space from the ends,
// so adding spaces for indentation will not currently work.
var cheekyLookRemarks = []TimedText{
	{0, "*click*"},
	{750, "*POP*!"},
	{1500, "Hmmm. The light bulb has gone out."},
	{2000, "Looking around is useless in an unlighted room."},
}

func lookAroundRoom(conn *websocket.Conn, req *GameonRequest, tail, room string) error {
	locus := "LOOK"
	checkpoint(locus, "AROUND")
	for _, tt := range cheekyLookRemarks {
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
		err = sendMsg(conn, req.UserId, j, MTPlayer)
		if err != nil {
			return err
		}
	}
	return nil
}
