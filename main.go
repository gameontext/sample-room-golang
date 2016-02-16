// Copyright (c) 2016 IBM Corp. All rights reserved.
// Use of this source code is governed by the Apache License,
// Version 2.0, a copy of which can be found in the LICENSE file.

// Main - it all starts here.
package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
)

// We operate in one of two modes:
//   - for a delete request, we get the work done and exit the program
//   - for a registration request, we register the room and then run
//     a server to handle websocket callbacks. Currently this server
//     runs pretty much forever until the program is killed.
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

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	if len(config.roomToDelete) > 0 {
		checkpoint(locus, fmt.Sprintf("deleteWithRetries %s", config.roomToDelete))
		err = deleteWithRetries(client, config.roomToDelete)
		if err != nil {
			checkpoint(locus, fmt.Sprintf("DELETE.FAILED err=%s", err.Error()))
		}
		return
	}

	checkpoint(locus, "registerWithRetries")

	err = registerWithRetries(client)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	checkpoint(locus, "startServer")
	startServer()
}

// Prints a simple checkpoint message.
func checkpoint(locus, s string) {
	fmt.Printf("CHECKPOINT: %s.%s\n", locus, s)
}
