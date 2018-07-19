// Copyright (c) 2016 IBM Corp. All rights reserved.
// Use of this source code is governed by the Apache License,
// Version 2.0, a copy of which can be found in the LICENSE file.

// Keeps track of players so that we can interact with them.
package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

const (
	// Special broadcast sender id
	TrackerSender = "tracker.is.sender"
)

// Player connections are used to track a player's
// time in a room; they are mostly useful for implementing
// chat broadcasts.
type PlayerConnection struct {
	playerId string
	roomId   string
	conn     *websocket.Conn
}

type Broadcast struct {
	// Broadcasts are restricted to this room.
	roomId  string
	message string
	// Sender can be faked for effect (messages from the room)
	// or they can be a real user name for actual chats.
	sender string
	// receiver is either "*" for everyone in the room or else
	// it is one specific user id.
	receiver string
}

type Banter struct {
	message string
	sender  string
}

type Tracker struct {
	players   map[string]*PlayerConnection
	add       chan *PlayerConnection
	remove    chan string
	broadcast chan *Broadcast
	smalltalk chan *Banter
}

var tracker = Tracker{
	players:   make(map[string]*PlayerConnection),
	add:       make(chan *PlayerConnection),
	remove:    make(chan string),
	broadcast: make(chan *Broadcast),
	smalltalk: make(chan *Banter),
}

func BroadcastMessage(r, m, sender, receiver string) {
	var bc = Broadcast{
		roomId:   r,
		message:  m,
		sender:   normalizeBroadcastSender(sender, receiver),
		receiver: receiver}
	tracker.broadcast <- &bc
}

func TrackPlayer(pc *PlayerConnection) {
	tracker.add <- pc
}

func UntrackPlayer(roomId, playerId string) {
	tracker.remove <- makePlayerKey(playerId, roomId)
}

func MakeSmalltalk(m, sender string) {
	var banter = Banter{message: m, sender: sender}
	tracker.smalltalk <- &banter
}

// Runs the player tracker. This should be started as a new
// goroutine before any callbacks are enabled.
func TrackPlayers() {
	checkpoint("TRACKER", "STARTED")
	for {
		select {
		case pc := <-tracker.add:
			logPlayer(pc, "ADDING", config.debug)
			tracker.players[makePlayerKey(pc.playerId, pc.roomId)] = pc
		case k := <-tracker.remove:
			pc := tracker.players[k]
			if pc == nil {
				checkpoint("TRACKER", fmt.Sprintf("k=%s not found", k))
			} else {
				logPlayer(pc, "REMOVING", config.debug)
				delete(tracker.players, k)
			}
		case bc := <-tracker.broadcast:
			broadcast(bc)
		case banter := <-tracker.smalltalk:
			smalltalk(banter)
		}
	}
}

func makePlayerKey(playerId, roomId string) string {
	return fmt.Sprintf("%s-%s", playerId, roomId)
}

func logPlayer(pc *PlayerConnection, disposition string, logging bool) {
	if !logging {
		return
	}
	locus := "PLAYER"
	checkpoint(locus, fmt.Sprintf("%s playerId=%s", disposition, pc.playerId))
	checkpoint(locus, fmt.Sprintf("%s roomId=%s", disposition, pc.roomId))
}

func logBroadcast(bc *Broadcast, note string, logging bool) {
	if !logging {
		return
	}
	locus := "BROADCAST"
	if len(note) > 0 {
		checkpoint(locus, fmt.Sprintf("note=%s", note))
	}
	checkpoint(locus, fmt.Sprintf("roomId=%s", bc.roomId))
	checkpoint(locus, fmt.Sprintf("sender=%s", bc.sender))
	checkpoint(locus, fmt.Sprintf("receiver=%s", bc.receiver))
	checkpoint(locus, fmt.Sprintf("message=%s", bc.message))
}

type ChatMessage struct {
	Rtype    string `json:"type,omitempty"`
	Username string `json:"username,omitempty"`
	Content  string `json:"content,omitempty"`
	Bookmark int    `json:"bookmark,omitempty"`
}

func broadcast(bc *Broadcast) {
	logBroadcast(bc, "candidate", config.debug)
	for _, pc := range tracker.players {
		r := pc.roomId
		if len(r) == 0 || r == bc.roomId {
			logBroadcast(bc, "sending", config.debug)
			c := pc.conn
			var m ChatMessage
			m.Rtype = "chat"
			m.Username = bc.sender
			m.Content = bc.message
			m.Bookmark = 0
			j, err := json.MarshalIndent(m, "", "    ")
			if err != nil {
				log.Printf("BROADCAST JSON ERROR\n")
				return
			}
			SendMessage(c, bc.receiver, j, MTPlayer)
		} else {
			if config.debug {
				log.Printf("BROADCAST.%s REJECT\n", r)
			}
		}
	}
}

var bcCounter = 0

// Normalize morphs the special tracker id into a tracker
// id with a message number appended.
func normalizeBroadcastSender(s, r string) string {
	if s == TrackerSender {
		bcCounter += 1
		return fmt.Sprintf("tracker.%04d", bcCounter)
	}
	return s
}

func smalltalk(banter *Banter) {
	for _, pc := range tracker.players {
		var m ChatMessage
		m.Rtype = "chat"
		m.Username = banter.sender
		m.Content = banter.message
		m.Bookmark = 0
		j, err := json.MarshalIndent(m, "", "    ")
		if err != nil {
			log.Printf("smalltalk JSON ERROR\n")
			return
		}
		SendMessage(pc.conn, "*", j, MTPlayer)
	}
}
