// Copyright (c) 2016 IBM Corp. All rights reserved.
// Use of this source code is governed by the Apache License,
// Version 2.0, a copy of which can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"strings"
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
func handleRoom(conn *websocket.Conn, req *GameonRequest) error {
	content := req.Content
	if len(content) < 1 {
		return JSPayloadError{"There is no content."}
	}
	if 0 == strings.Index(content, "/") {
		return handleSlashCommand(conn, req)
	}
	return handleChat(conn, req)
}

const (
	// Slash commands, without the actual '/', of course.
	slashExamine   = "EXAMINE"
	slashExit      = "EXIT"
	slashGo        = "GO"
	slashHelp      = "HELP"
	slashInventory = "INVENTORY"
	slashLook      = "LOOK"
)

var roomCommands = []string{slashExamine, slashExit, slashGo, slashHelp,
	slashInventory, slashLook}
var slashCommands = []string{"/" + slashExamine, "/" + slashExit, "/" + slashGo,
	"/" + slashHelp, "/" + slashInventory, "/" + slashLook}

// Recognizes and dispatches a room slash command. Nil is returned
// if all goes well, otherwise an error is returned.
func handleSlashCommand(conn *websocket.Conn, req *GameonRequest) error {
	locus := "HANDLE.SLASH"
	cmd, tail, err := parseCommandPrefix(req.Content)
	if err != nil {
		return err
	}
	checkpoint(locus, fmt.Sprintf("cmd=%s tail=%s", cmd, tail))
	switch cmd {
	case slashGo, slashExit:
		return exitRoom(conn, req, tail)
	case slashLook:
		return lookAroundRoom(conn, req, tail)
	case slashHelp:
		return helpCommand(conn, req, tail)
	case slashInventory:
		return checkInventory(conn, req, tail)
	case slashExamine:
		return examineObject(conn, req, tail)
	default:
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
	haystack := strings.ToUpper(s[1:])
	hlen := len(haystack)
	for _, key := range roomCommands {
		if config.debug {
			fmt.Printf("parseCommandPrefix: '%s' vs '%s'\n", key, haystack)
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
			tail = strings.Trim(haystack[n:], " ")
			return
		}
	}
	err = JSPayloadError{fmt.Sprintf("Unrecognized command in '%s'", s)}
	return
}

func sendPlayerResp(conn *websocket.Conn, targetid string, j []byte) error {
	return sendResp(conn, targetid, j, MTPlayer)
}

func sendResp(conn *websocket.Conn, targetid string, j []byte, rtype string) error {
	var m = fmt.Sprintf("%s,%s,%s", rtype, targetid, string(j))
	err := conn.WriteMessage(expectedMessageType, []byte(m))
	if config.debug {
		fmt.Printf("sendResp(%s)\n", m)
		if err == nil {
			fmt.Printf("SENT OKAY\n")
		} else {
			fmt.Printf("SEND FAILED: %s\n", err.Error())
		}
	}
	return err
}
