// Copyright (c) 2016 IBM Corp. All rights reserved.
// Use of this source code is governed by the Apache License,
// Version 2.0, a copy of which can be found in the LICENSE file.

// Room registration.

package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

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

type RoomRegistrationReq struct {
	Name              string      `json:"name,omitempty"`
	FullName          string      `json:"fullName,omitempty"`
	ConnectionDetails ConnDetails `json:"connectionDetails,omitempty"`
	Doors             DoorGroup   `json:"doors,omitempty"`
}

type RoomRegistrationResp struct {
	Id    string              `json:"_id,omitempty"`
	Rev   string              `json:"_rev,omitempty"`
	Owner string              `json:"owner,omitempty"`
	Info  RoomRegistrationReq `json:"info,omitempty"`
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
			rememberMyRooms()
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
		traceResponseBody(locus, resp, body)
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
		traceResponseBody(locus, resp, body)
		return
	default:
		checkpoint(locus, "UnsupportedStatusCode")
		printResponseBody(locus, resp, body)
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

	u := fmt.Sprintf("%s://%s/map/v1/sites", config.protocol, config.gameonAddr)

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

	addAuthenticationHeaders(req, registration)

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
		traceResponseBody(locus, resp, body)
		err = rememberRegistration(resp, body)
		return
	case http.StatusConflict:
		err = RegError{fmt.Sprintf("Bad status: %s", resp.Status)}
		checkpoint(locus, fmt.Sprintf("Internal Error. Attempt to reregister. Status=%s", resp.Status))
		printResponseBody(locus, resp, body)
		return
	default:
		err = RegError{fmt.Sprintf("Unhandled Status: %s", resp.Status)}
		checkpoint(locus, fmt.Sprintf("Unhandled Status=%s", resp.Status))
		printResponseBody(locus, resp, body)
		return
	}
}

func printRoomRegistrationResp(locus string, r *RoomRegistrationResp) {
	j, err := json.MarshalIndent(r, "", "    ")
	if err == nil {
		fmt.Printf("\n%s\n%s\n", locus, string(j))
	}
}

// Generate a JSON string containing our registration info.
func genRegistration() (rs string, err error) {
	var reg RoomRegistrationReq
	reg.Name = config.roomName
	reg.FullName = config.roomName
	// Door descriptions are collected from an inside-looking-out
	// perspective, but Game On! wants a description from the
	// connecting room's point of view. So, our commandline
	// North is what GameOn! wants for the South.
	reg.Doors.North = config.south
	reg.Doors.South = config.north
	reg.Doors.East = config.west
	reg.Doors.West = config.east
	reg.Doors.Up = config.down
	reg.Doors.Down = config.up
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
func traceResponseBody(locus string, r *http.Response, body string) {
	if config.debug {
		printResponseBody(locus, r, body)
	}
}

// Unconditionally print an http.Response body string.
func printResponseBody(locus string, r *http.Response, body string) {
	fmt.Printf("\n%s.RESP.StatusCode=%s\n", locus, r.Status)
	if r.StatusCode == http.StatusNotFound {
		// Avoid the noise.
		return
	}
	fmt.Printf("%s.RESP.Body='%s'\n", locus, body)
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

type RoomQueryResp struct {
	Id    string              `json:"_id,omitempty"`
	Owner string              `json:"owner,omitempty"`
	Info  RoomRegistrationReq `json:"info,omitempty"`
}

var MyRooms map[string]string

func rememberMyRooms() (err error) {
	locus := "REG.LISTMYROOMS"
	MyRooms = make(map[string]string)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	u := fmt.Sprintf("%s://%s/map/v1/sites?owner=%s",
		config.protocol, config.gameonAddr, config.id)
	checkpoint(locus, u)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		checkpoint(locus, fmt.Sprintf("NewRequest.Error err=%s", err.Error()))
		return
	}

	addAuthenticationHeaders(req, "")

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
	traceResponseBody(locus, resp, body)
	switch resp.StatusCode {
	case http.StatusOK:
		var qr [25]RoomQueryResp
		e := json.Unmarshal([]byte(body), &qr)
		if e != nil {
			fmt.Printf("%s: JSON unmarshalling error: %s\n",
				locus, e.Error())
			fmt.Printf("%s : Offending JSON: %s\n", locus, body)
			return
		}
		for _, r := range qr {
			if len(r.Id) > 0 {
				MyRooms[r.Id] = r.Info.FullName
			}
		}
		if config.debug {
			for k, v := range MyRooms {
				fmt.Printf("%s --> %s\n", k, v)
			}
		}
		return
	default:
		checkpoint(locus, "FAILED.")
	}
	return
}
