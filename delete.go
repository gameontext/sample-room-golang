// Copyright (c) 2016 IBM Corp. All rights reserved.
// Use of this source code is governed by the Apache License,
// Version 2.0, a copy of which can be found in the LICENSE file.

// Delete our room

package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

// Deletes our room with the GameOn! server, with retries on connection failure.
func deleteWithRetries() (e error) {
	locus := "DELETE_W_RETRIES"
	checkpoint(locus, fmt.Sprintf("retries=%d secondsBetween=%d", config.retries, config.secondsBetween))
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	for i := 0; i < config.retries; i++ {
		checkpoint(locus, fmt.Sprintf("Begin attempt %d of %d",
			i+1, config.retries))
		var stopTrying bool
		stopTrying, e = deleteRoom(client)
		if e == nil {
			if stopTrying {
				checkpoint(locus, "Room deletion failed")
			} else {
				checkpoint(locus, "Room deletion was successful.")
			}
			return
		}
		checkpoint(locus, fmt.Sprintf("sleeping %d seconds.", config.secondsBetween))
		if i+1 < config.retries {
			time.Sleep(time.Duration(config.secondsBetween) * time.Second)
		}
	}
	checkpoint(locus, "Room deletion failed.")
	e = RegError{fmt.Sprintf("Timed out. Last error: %s", e.Error())}
	return
}

func deleteRoom(client *http.Client) (stopTrying bool, err error) {
	locus := "DELETE.ROOM"
	checkpoint(locus, "Begin")
	ts := makeTimestamp()
	bodyHash := hash("")
	tokens := []string{config.id, ts}
	sig := buildHmac(tokens, config.secret)
	var u string
	if config.localServer {
		u = fmt.Sprintf("http://%s/map/v1/sites/%s", config.gameonAddr, config.roomToDelete)
	} else {
		u = fmt.Sprintf("https://%s/map/v1/sites/%s", config.gameonAddr, config.roomToDelete)
	}

	checkpoint(locus, fmt.Sprintf("URL %s", u))

	req, err := http.NewRequest("DELETE", u, nil)
	if err != nil {
		checkpoint(locus, fmt.Sprintf("NewRequest.Error err=%s", err.Error()))
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json,text/plain")
	req.Header.Set("gameon-id", config.id)
	req.Header.Set("gameon-date", ts)
	req.Header.Set("gameon-sig-body", bodyHash)
	req.Header.Set("gameon-signature", sig)
	if config.debug {
		for _, k := range []string{"gameon-id", "gameon-date", "gameon-sig-body", "gameon-signature"} {
			fmt.Printf("%s=%s\n", k, req.Header.Get(k))
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		checkpoint(locus, fmt.Sprintf("DELETE.Error err=%s", err.Error()))
		return
	}
	body, err := extractBody(resp)
	if err != nil {
		checkpoint(locus, fmt.Sprintf("Body.Error err=%s", err.Error()))
		return
	}
	checkpoint(locus, fmt.Sprintf("Status=%s", resp.Status))
	switch resp.StatusCode {
	case http.StatusNoContent:
		checkpoint(locus, "Deleted")
		return
	case http.StatusForbidden:
	case http.StatusNotFound:
		checkpoint(locus, "Sigh. There is no use trying any more.")
		printResponseBody(locus, resp, body)
		stopTrying = true
		return
	default:
		err = RegError{fmt.Sprintf("Unhandled Status: %s", resp.Status)}
		checkpoint(locus, fmt.Sprintf("Unhandled Status=%s", resp.Status))
		printResponseBody(locus, resp, body)
		return
	}
	return
}
