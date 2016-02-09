// Copyright (c) 2016 IBM Corp. All rights reserved.
// Use of this source code is governed by the Apache License,
// Version 2.0, a copy of which can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
)

const localSecret = "MyRegistrationSecret"

type RoomConfig struct {
	gameonAddr   string
	callbackAddr string
	// The callback port is our externally-visible published port.
	// We send this to the gameOn server so that it can call back.
	callbackPort int
	// The listening port is the port that our server listens to
	// internally.  For example, if this code is placed into a
	// container we would listen internally to one port, say 3000,
	// and map that to a diffent host port. In this case that different
	// host port must be our callbackPort.
	listeningPort int
	// This is the name of our room; use this name when connecting
	// to our room from another using north, south, east and west.
	roomName string
	// Rooms we connect to from any given direction.
	north, south, east, west, up, down string
	debug                              bool
	id                                 string
	key                                string
	owner                              string
	localServer                        bool
}

var config RoomConfig

func main() {
	locus := "MAIN"
	checkpoint(locus, "processCommandLine")
	err := processCommandline()
	if err != nil {
		fmt.Println(err.Error())
		flag.Usage()
		return
	}
	printConfig(&config)

	checkpoint(locus, "registerWithRetries")
	err = registerWithRetries()
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	checkpoint(locus, "startServer")
	startServer()
}

// Processes the commandline.
// Sets values in the config struct and returns nil
// if successful or an error otherwise
func processCommandline() (err error) {
	err = nil
	flag.StringVar(&config.gameonAddr, "g", "", "Game-on server address")
	flag.StringVar(&config.callbackAddr, "c", "", "Our published callback address")
	flag.IntVar(&config.callbackPort, "cp", -1, "Our published callback port")
	flag.IntVar(&config.listeningPort, "lp", -1, "Our listening port")
	flag.StringVar(&config.roomName, "r", "", "Our room name.")
	flag.BoolVar(&config.debug, "d", false, "Enables debug mode")
	flag.StringVar(&config.north, "north", "A frost-covered door leads to the north.", "/exit N takes us here")
	flag.StringVar(&config.south, "south", "A moss-covered door leads to the south", "/exit S takes us here")
	flag.StringVar(&config.east, "east", "A badly-painted door opens to the east.", "/exit E takes us here")
	flag.StringVar(&config.west, "west", "An old swinging door leads west.", "/exit W takes us here")
	flag.StringVar(&config.up, "up", "There is a rickety set of steps leading up.", "/exit W takes us here")
	flag.StringVar(&config.down, "down", "Heat eminates from an opening in the floor.", "/exit W takes us here")
	flag.StringVar(&config.id, "id", "", "The id associated with our key.")
	flag.StringVar(&config.key, "key", localSecret, "Our secret key.")
	flag.StringVar(&config.owner, "owner", "", "Our user name")
	flag.BoolVar(&config.localServer, "local", false, "We are using a local server")

	flag.Parse()
	if config.gameonAddr == "" {
		err = ArgError{"Missing Game-on server address."}
		return
	}
	if config.callbackAddr == "" {
		err = ArgError{"Missing callback address."}
		return
	}
	if config.callbackPort < 0 {
		err = ArgError{"Missing or invalid callback port."}
		return
	}
	if config.listeningPort < 0 {
		config.listeningPort = config.callbackPort
	}
	if config.roomName == "" {
		config.roomName = fmt.Sprintf("ROOM.%05d", config.callbackPort)
	}
	if config.owner == "" {
		config.owner = config.id
	}
	return
}

func printConfig(c *RoomConfig) {
	fmt.Printf("gameonAddr=%s\n", config.gameonAddr)
	fmt.Printf("callbackAddr=%s\n", config.callbackAddr)
	fmt.Printf("callbackPort=%d\n", config.callbackPort)
	fmt.Printf("roomName=%s\n", config.roomName)
	fmt.Printf("north=%s\n", config.north)
	fmt.Printf("south=%s\n", config.south)
	fmt.Printf("east=%s\n", config.east)
	fmt.Printf("west=%s\n", config.west)
	fmt.Printf("debug=%v\n", config.debug)
}

// Print a simple checkpoint message.
func checkpoint(locus, s string) {
	fmt.Printf("\nCHECKPOINT: %s.%s\n", locus, s)
}
