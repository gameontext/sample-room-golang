// Copyright (c) 2016 IBM Corp. All rights reserved.
// Use of this source code is governed by the Apache License,
// Version 2.0, a copy of which can be found in the LICENSE file.

// Recognizes and routes room commands
package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"strings"
	log "github.com/sirupsen/logrus"
)

const (
	// GameOn Message Types
	MTPlayer         = "player"
	MTPlayerLocation = "playerLocation"
)

// Handles commands specific to, and implemented by, our room.
// On entry, req contains the unmarshalled JSON payload which
// must contain a non-empty Content field. If the Content field
// begins with a slash ("/") then the request is treated as a
// room command, otherwise it is treated as a chat request.
func handleRoom(conn *websocket.Conn, req *GameonRequest, room string) error {
	content := req.Content
	if len(content) < 1 {
		return JSPayloadError{"There is no content."}
	}
	if 0 == strings.Index(content, "/") {
		return handleSlashCommand(conn, req, room)
	}
	return handleChat(conn, req, room)
}

const (
	// Slash commands, without the actual '/', of course.
	slashExamine   = "EXAMINE"
	slashGo        = "GO"
	slashInventory = "INVENTORY"
	slashLook      = "LOOK"
	slashWink      = "WINK"
)

// The is the list of commands that we are willing to catch.
var commandsWeSupport = []string{slashExamine, slashGo, slashInventory, slashLook, slashWink}

type CommandDesc struct {
	cmd  string
	desc string
}

// This is the list of commands that we add, over and above the normal set.
// When a player enters our room, we will need to tell them game about these
// commands so that the game knows to add them to /help output.
var commandsWeAdd = []CommandDesc{
	{"/wink", "(You wonder what this would do.)"},
}

// Recognizes and dispatches a room slash command. Nil is returned
// if all goes well, otherwise an error is returned.
func handleSlashCommand(conn *websocket.Conn, req *GameonRequest, room string) error {
	locus := "HANDLE.SLASH"
	cmd, tail, err := parseCommandPrefix(req.Content)
	if err != nil {
		SendMessageToPlayer(conn, "What? I didn't understand that.", req.UserId)
		return err
	}
	checkpoint(locus, fmt.Sprintf("cmd=%s tail=%s", cmd, tail))
	switch cmd {
	case slashGo:
		return exitRoom(conn, req, tail, room)
	case slashLook:
		return lookAroundRoom(conn, req, tail, room)
	case slashInventory:
		return checkInventory(conn, req, tail, room)
	case slashExamine:
		return examineObject(conn, req, tail, room)
	case slashWink:
		return wink(conn, req, tail, room)
	default:
		SendMessageToPlayer(conn, "What? I didn't understand that.", req.UserId)
		return JSPayloadError{fmt.Sprintf("Unrecognized command: '%s'", cmd)}
	}
}

// Parse the Content field of a request and return the
// command string and the remainder if the command is
// recognized, otherwise return a non-nil error.
func parseCommandPrefix(s string) (cmd, tail string, err error) {
	const minCommandLen = 2
	if len(s) < minCommandLen {
		err = JSPayloadError{"Invalid command format: Less than minimum length."}
		return
	}
	if "/" != s[0:1] {
		err = JSPayloadError{"Invalid command format: Missing leading slash"}
		return
	}
	haystackAsis := s[1:]
	haystack := strings.ToUpper(s[1:])
	hlen := len(haystack)
	for _, key := range commandsWeSupport {
		if config.debug {
			log.Printf("parseCommandPrefix: '%s' vs '%s'\n", key, haystack)
		}
		n := len(key)
		if hlen < n {
			continue
		}
		if haystack[0:n] != key {
			continue
		}
		// Look for an exact match or else our
		// key followed by a space.
		if n == hlen {
			cmd = key
			return
		}
		if " " == haystack[n:n+1] {
			cmd = key
			tail = strings.Trim(haystackAsis[n:], " ")
			return
		}
	}
	err = JSPayloadError{fmt.Sprintf("Unrecognized command in '%s'", s)}
	return
}
