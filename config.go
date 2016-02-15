// Copyright (c) 2016 IBM Corp. All rights reserved.
// Use of this source code is governed by the Apache License,
// Version 2.0, a copy of which can be found in the LICENSE file.

// Configuration
package main

import (
	"flag"
	"fmt"
)

const localSecret = "MyRegistrationSecret"

// A RoomConfig struct contains important, frequently-needed, information
// about our current room registration as well as our interaction with
// our Game On! server. Typically these values are collected from commandline
// arguments.
type RoomConfig struct {
	// The basic address of the GameOn! server.
	gameonAddr   string
	callbackAddr string
	// The callback port is the externally-visible published port that
	// that will receive websocket traffic from the GameOn! server.
	callbackPort int
	// The listening port is the port that our websocket server listens
	// to internally.  For example, if this code is placed into a
	// container we could choose to always listen to one port, say 3000,
	// and map that to a different host port.
	listeningPort int
	// This is the name of our room; use this name when connecting
	// to our room from another using north, south, east and west.
	roomName string
	// Text descriptions of doors that connect us to other rooms.
	// Although up and down are provided, the GameOn! server typically
	// ignores any direction other than n,w,e or w.
	north, south, east, west, up, down string
	// If true, print a bunch of debugging information useful mostly
	// to programmers.
	debug bool
	// This is the authorization id that was obtained during your
	// GameOn! browser login. It might look something like this
	// if you logged in using your Google ID:
	// 'google:132043241444397884152'
	id string
	// This is the shared secret that was obtained during your GameOn!
	// browser login. If you logged in using your Google ID it might
	// look like this: 'LNIkaoiu62addlGp/rCZc7g,n3s9jUtOpXErr062kos='
	secret         string
	localServer    bool
	timeShift      int
	retries        int
	secondsBetween int
	// This is a room id and it is only used in the context of a
	// delete request.
	roomToDelete string
	// The protocol to be used when talking to the game server.
	protocol string
}

// config is our single, package-wide, source of configuration data.
// (cmdling processing) ==> config ==> (used by code in this package)
var config RoomConfig

// Processes our commandline and establishes the contents of
// the package-wide config struct.
//
// Returns nil if successful or an error otherwise
func processCommandline() (err error) {
	flag.StringVar(&config.gameonAddr, "g", "", "GameOn! server address")
	flag.StringVar(&config.callbackAddr, "c", "", "Our published callback address")
	flag.IntVar(&config.callbackPort, "cp", -1, "Our published callback port")
	flag.IntVar(&config.listeningPort, "lp", -1, "Our listening port")
	flag.StringVar(&config.roomName, "r", "", "Our room name.")
	flag.BoolVar(&config.debug, "d", false, "Enables debug mode")
	flag.StringVar(&config.north, "north", "A frost-covered door leads to the north.", "Describes our northern door")
	flag.StringVar(&config.south, "south", "A moss-covered door leads to the south", "Describes our southern door")
	flag.StringVar(&config.east, "east", "A badly-painted door opens to the east.", "Describes our eastern door")
	flag.StringVar(&config.west, "west", "An old swinging door leads west.", "Describes our western door")
	flag.StringVar(&config.up, "up", "There is a rickety set of steps leading up.", "GameOn! often ignores this door")
	flag.StringVar(&config.down, "down", "Heat eminates from an opening in the floor.", "GameOn! often ignores this door")
	flag.StringVar(&config.id, "id", "", "The id associated with our shared secret.")
	flag.StringVar(&config.secret, "secret", localSecret, "Our shared secret.")
	flag.BoolVar(&config.localServer, "local", false, "We are using a local server. Local servers expect http://; remote servers expect https://")
	flag.IntVar(&config.timeShift, "ts", 0, "The number of milleseconds to add or subtract from our timestamp so that we can better match the server clock")
	flag.IntVar(&config.retries, "retries", 5, "The number of initial registration attempts.")
	flag.IntVar(&config.secondsBetween, "between", 5, "The number of seconds between registration attempts.")
	flag.StringVar(&config.roomToDelete, "delete", "", "Delete the room with this id and exit.")

	flag.Parse()
	if config.gameonAddr == "" {
		err = ArgError{"Missing Game-on server address."}
		return
	}
	if len(config.roomToDelete) == 0 {
		// This is not a deletion request so make sure the information
		// we need to register a room and run the websocket server is valid.
		if config.callbackAddr == "" {
			err = ArgError{"Missing callback address."}
			return
		}
		if config.callbackPort < 0 {
			err = ArgError{"Missing or invalid callback port."}
			return
		}
		if config.listeningPort < 0 {
			// listening port defaults to callback port
			config.listeningPort = config.callbackPort
		}
		if config.roomName == "" {
			config.roomName = fmt.Sprintf("ROOM.%05d", config.callbackPort)
		}
	}
	if config.localServer {
		config.protocol = "http"
	} else {
		config.protocol = "https"
	}
	return
}

func printConfig(c *RoomConfig) {
	fmt.Printf("gameonAddr=%s\n", config.gameonAddr)
	// Many things are useless when we are just doing a delete.
	if config.roomToDelete == "" {
		fmt.Printf("callbackAddr=%s\n", config.callbackAddr)
		fmt.Printf("callbackPort=%d\n", config.callbackPort)
		fmt.Printf("listeningPort=%d\n", config.listeningPort)
		fmt.Printf("roomName=%s\n", config.roomName)
		fmt.Printf("north=%s\n", config.north)
		fmt.Printf("south=%s\n", config.south)
		fmt.Printf("east=%s\n", config.east)
		fmt.Printf("west=%s\n", config.west)
	}
	fmt.Printf("debug=%v\n", config.debug)
	fmt.Printf("roomToDelete=%v\n", config.roomToDelete)
	fmt.Printf("localServer=%v\n", config.localServer)
	fmt.Printf("timeShift=%d\n", config.timeShift)
	if config.debug {
		fmt.Printf("id=%s\n", config.id)
		fmt.Printf("secret=%s\n", config.secret)
	}
}
