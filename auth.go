// Copyright (c) 2016 IBM Corp. All rights reserved.
// Use of this source code is governed by the Apache License,
// Version 2.0, a copy of which can be found in the LICENSE file.

// Request authentication

package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
)

func addAuthenticationHeaders(req *http.Request, body string) {
	// Build the authentication values that Game On! requires. If
	// the body is empty, do not include it in the calculations.
	var bodyHash, sig string
	ts := makeTimestamp()
	if len(body) > 0 {
		bodyHash = hash(body)
		tokens := []string{config.id, ts, bodyHash}
		sig = buildHmac(tokens, config.secret)
	} else {
		bodyHash = hash("")
		tokens := []string{config.id, ts}
		sig = buildHmac(tokens, config.secret)
	}
	// Set the required headers.
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
}

func hash(message string) string {
	h := sha256.New()
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func buildHmac(tokens []string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	s := ""
	for _, t := range tokens {
		s += t
	}
	h.Write([]byte(s))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// Returns the current time as a UTC-formatted string.
// If config.timeShift is non-zero, then the timestamp will
// be shifted by config.timeShift milliseconds. This can be
// used to slide our registration timestamp closer to the
// clock on a remote GameOn! server.
func makeTimestamp() string {
	if config.timeShift == 0 {
		return time.Now().UTC().Format(time.RFC3339Nano)
	}
	locus := "MAKE.TIMESTAMP"
	t1 := time.Now()
	t2 := t1.Add(time.Duration(config.timeShift) * time.Millisecond)
	ourTime := t1.UTC().Format(time.RFC3339Nano)
	serverTime := t2.UTC().Format(time.RFC3339Nano)
	checkpoint(locus, fmt.Sprintf("ourTime    %s", ourTime))
	checkpoint(locus, fmt.Sprintf("serverTime %s", serverTime))
	return serverTime
}
