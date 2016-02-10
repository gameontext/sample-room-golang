// Copyright (c) 2016 IBM Corp. All rights reserved.
// Use of this source code is governed by the Apache License,
// Version 2.0, a copy of which can be found in the LICENSE file.

// Room registration.

package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// A RoomExit describes an exit out of a GameOn! room.
// Here is an example of a corresponding JSON fragment:
//  {
//    "id":"firstroom",
//    "name":"First Room",
//    "fullName":"The First Room",
//    "door":"A warped wooden door with a friendly face branded on the corner"
//  }
type RoomExit struct {
	Id       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	FullName string `json:"fullName,omitempty"`
	Door     string `json:"door,omitempty"`
}

type ConnDetails struct {
	Type   string `json:"type,omitempty"`
	Target string `json:"target,omitempty"`
}

type DoorGroup struct {
	North string `json:"n,omitempty"`
	South string `json:"s,omitempty"`
	East  string `json:"e,omitempty"`
	West  string `json:"w,omitempty"`
	// Up is allowed, but is currently ignored.
	Up string `json:"u,omitempty"`
	// Down is allowed, but is currently ignored
	Down string `json:"d,omitempty"`
}

// TODO remove this
// Exit information, which we receive as part of a room regstration response,
// is currently kept in a struct. Consider using a map.
type ExitGroup struct {
	North RoomExit `json:"n,omitempty"`
	South RoomExit `json:"s,omitempty"`
	East  RoomExit `json:"e,omitempty"`
	West  RoomExit `json:"w,omitempty"`
	Up    RoomExit `json:"u,omitempty"`
	Down  RoomExit `json:"d,omitempty"`
}

type RoomRegistrationReq struct {
	Name              string      `json:"name,omitempty"`
	FullName          string      `json:"fullName,omitempty"`
	ConnectionDetails ConnDetails `json:"connectionDetails,omitempty"`
	Doors             DoorGroup   `json:"doors,omitempty"`
}

type RoomCoord struct {
	X int `json:x,omitempty`
	Y int `json:y,omitempty`
}

type RoomRegistrationResp struct {
	Id    string              `json:"_id,omitempty"`
	Rev   string              `json:"_rev,omitempty"`
	Owner string              `json:"owner,omitempty"`
	Info  RoomRegistrationReq `json:"info,omitempty"`
	Exits ExitGroup           `json:"exits,omitempty"`
	Coord RoomCoord           `json:"coord,omitempty`
	Type  string              `json:"type,omitempty"`
}

// Registers our room with the GameOn! server, with retries on failure.
func registerWithRetries() (e error) {
	locus := "REG_W_RETRIES"
	checkpoint(locus, fmt.Sprintf("retries=%d secondsBetween=%d", config.retries, config.secondsBetween))
	for i := 0; i < config.retries; i++ {
		checkpoint(locus, fmt.Sprintf("Begin attempt %d of %d",
			i+1, config.retries))
		e = register()
		if e == nil {
			checkpoint(locus, "Registration was successful.")
			return
		}
		checkpoint(locus, fmt.Sprintf("sleeping %d seconds.", config.secondsBetween))
		if i+1 < config.retries {
			time.Sleep(time.Duration(config.secondsBetween) * time.Second)
		}
	}
	checkpoint(locus, "Registration failed.")
	e = RegError{fmt.Sprintf("Timed out. Last error: %s", e.Error())}
	return
}

// Registers our room with the game-on server if the room is not
// already registered.
func register() (err error) {
	locus := "REG"
	checkpoint(locus, "Begin")
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	var registered bool
	registered, err = checkForPriorRegistration(client)
	if err != nil {
		checkpoint(locus, fmt.Sprintf("err=%s", err.Error()))
		return
	}
	if registered {
		checkpoint(locus, "WasAlreadyRegistered")
		return
	}
	checkpoint(locus, "WeNeedToRegister")
	err = registerOurRoom(client)
	if err == nil {
		checkpoint(locus, "Registered")
	} else {
		checkpoint(locus, fmt.Sprintf("RegistrationFailed err=%s", err.Error()))
	}
	return
}

func checkForPriorRegistration(client *http.Client) (registered bool, err error) {
	locus := "REG.CHECKPRIOR"
	checkpoint(locus, "Begin")
	queryParams := fmt.Sprintf("name=%s&owner=%s", config.roomName, config.id)
	u := fmt.Sprintf("http://%s/map/v1/sites?%s",
		config.gameonAddr,
		queryParams)
	if config.debug {
		checkpoint(locus, fmt.Sprintf(".GET.URL %s", u))
	}

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		checkpoint(locus, fmt.Sprintf("NewRequest.Error err=%s", err.Error()))
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		checkpoint(locus, fmt.Sprintf("Get.Error err=%s", err.Error()))
		return
	}
	body, err := extractBody(resp)
	if err != nil {
		checkpoint(locus, fmt.Sprintf("Body.Error err=%s", err.Error()))
		return
	}
	switch resp.StatusCode {
	case http.StatusOK:
		checkpoint(locus, "AlreadyRegistered")
		registered = true
		traceResponseBody(resp, body)
		if config.debug {
			// We expect a one element array from our query.
			// Do not propagate any errors from this debug logging.
			var regResp [1]RoomRegistrationResp
			e := json.Unmarshal([]byte(body), &regResp)
			if e != nil {
				fmt.Printf("%s: JSON unmarshalling error: %s\n",
					locus, e.Error())
				fmt.Printf("%s : Offending JSON: %s\n", locus, body)
				return
			}
			rememberedRegistration.Id = regResp[0].Id
			printRoomRegistrationResp(locus, &regResp[0])
		}
		return
	case http.StatusNoContent:
		checkpoint(locus, "NotCurrentlyRegistered")
		traceResponseBody(resp, body)
		return
	default:
		checkpoint(locus, "UnsupportedStatusCode")
		printResponseBody(resp, body)
		return
	}
}

func registerOurRoom(client *http.Client) (err error) {
	locus := "REG.REGROOM"
	checkpoint(locus, "Begin")
	var registration string
	registration, err = genRegistration()
	if err != nil {
		fmt.Printf("\nREG.REGROOM.InternalError\n")
		fmt.Printf("Internal registration error: %s\n", err.Error())
		return
	}
	ts := makeTimestamp()
	bodyHash := hash(registration)
	tokens := []string{config.id, ts, bodyHash}
	sig := buildHmac(tokens, config.secret)
	var u string
	if config.localServer {
		u = fmt.Sprintf("http://%s/map/v1/sites", config.gameonAddr)
	} else {
		u = fmt.Sprintf("https://%s/map/v1/sites", config.gameonAddr)
	}

	if config.debug {
		fmt.Printf("\nREG.POST URL: %s\n", u)
		fmt.Println("----- registration json begin -----")
		fmt.Println(registration)
		fmt.Println("----- registration json end -----")
	}

	req, err := http.NewRequest("POST", u, strings.NewReader(registration))
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
		checkpoint(locus, fmt.Sprintf("Post.Error err=%s", err.Error()))
		return
	}
	body, err := extractBody(resp)
	if err != nil {
		checkpoint(locus, fmt.Sprintf("Body.Error err=%s", err.Error()))
		return
	}
	checkpoint(locus, fmt.Sprintf("Status=%s", resp.Status))
	switch resp.StatusCode {
	case http.StatusCreated:
		checkpoint(locus, "Registered")
		err = rememberRegistration(resp, body)
		return
	case http.StatusConflict:
		err = RegError{fmt.Sprintf("Bad status: %s", resp.Status)}
		checkpoint(locus, fmt.Sprintf("Internal Error. Attempt to reregister. Status=%s", resp.Status))
		printResponseBody(resp, body)
		return
	default:
		err = RegError{fmt.Sprintf("Unhandled Status: %s", resp.Status)}
		checkpoint(locus, fmt.Sprintf("Unhandled Status=%s", resp.Status))
		printResponseBody(resp, body)
		return
	}
}

func printRoomRegistrationResp(locus string, r *RoomRegistrationResp) {
	j, err := json.MarshalIndent(r, "", "    ")
	if err == nil {
		fmt.Printf("\n%s\n%s", locus, string(j))
	}
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

// Generate a JSON string containing our registration info.
func genRegistration() (rs string, err error) {
	var reg RoomRegistrationReq
	reg.Name = config.roomName
	reg.FullName = config.roomName
	reg.Doors.North = config.north
	reg.Doors.South = config.south
	reg.Doors.East = config.east
	reg.Doors.West = config.west
	reg.Doors.Up = config.up
	reg.Doors.Down = config.down
	reg.ConnectionDetails.Type = "websocket"
	reg.ConnectionDetails.Target = fmt.Sprintf("ws://%s:%d",
		config.callbackAddr, config.callbackPort)
	j, err := json.MarshalIndent(reg, "", "    ")
	if err != nil {
		return
	}
	rs = string(j)
	return
}

// Conditionally print an http.Response body string
// if config.debug is true.
func traceResponseBody(r *http.Response, body string) {
	if config.debug {
		printResponseBody(r, body)
	}
}

// Unconditionally print an http.Response body string.
func printResponseBody(r *http.Response, body string) {
	fmt.Printf("RESP.StatusCode=%s\n", r.Status)
	if r.StatusCode == http.StatusNotFound {
		// Avoid the noise.
		return
	}
	fmt.Printf("RESP.Body='%s'\n", body)
}

// We stash our room registration response in rememberedRegistration
// so that we can use it later if we choose to do so.
var rememberedRegistration RoomRegistrationResp

// Remember our registration data in case we chose to use it later.
func rememberRegistration(r *http.Response, body string) (err error) {
	locus := "REG.REMEMBER"
	var rememberedRegistration RoomRegistrationResp
	err = json.Unmarshal([]byte(body), &rememberedRegistration)
	if err != nil {
		fmt.Printf("rememberRegistration : A JSON unmarshalling error occured: %s\n",
			err.Error())
		fmt.Printf("rememberRegistration :Offending JSON: %s\n", body)
		return
	}
	if config.debug {
		printRoomRegistrationResp(locus, &rememberedRegistration)
	}
	return
}

// Returns the body string from an http.Response. This can
// only be called once as the body is closed after reading
// its contents.
func extractBody(r *http.Response) (body string, e error) {
	defer r.Body.Close()
	b, e := ioutil.ReadAll(r.Body)
	if e == nil {
		body = string(b)
	}
	return
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
