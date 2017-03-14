// Copyright (c) 2016 IBM Corp. All rights reserved.
// Use of this source code is governed by the Apache License,
// Version 2.0, a copy of which can be found in the LICENSE file.

// Documentation
package main

// This code takes a microservice approach to inserting a room
// into a running Game On! instance. We consume microservices
// provided by Game On! to register a room and we provide our own
// microservice that allows Game On! to bring our room to life.
// In the process of responding to requests from Game On!, we in
// turn will use additonal Game On! microservices to send notifications
// to players and rooms and also to manage certain player aspects,
// such as their next location, during their time in our room.

// The Dance
//
// In this discussion we assume that we are consistent in using
// the same websocket callback every time that our program is
// started.
//
// On starting up, we check to see if our room (by name) has
// already been registered. If not, then we register it using
// an authenticated registration POST.  Once we have determined
// that our room is registered, we use an unauthticated GET to
// gather the names of all rooms that we currently have registered.
// (This code is capable of handling multiple rooms as long as each
// room was registered using the same callback address.)
//
// At this point we start our websocket server to listen for
// service requests from Game On!; the websocket server runs
// forever until our program is terminated.
//
// Ack'ing new connections
//
// Game On! opens a new websocket connection each time a player
// enters our room. As soon as the connection is opened and ready
// to use the websocket protocol, we must send an ack message back
// to Game On!. At the time of this writing, the sole purpose of the
// ack message is to send the list of procotol versions that this
// room supports.
//
// Future work:
//
// Reconnects to previously-registered rooms are handled
// somewhat naively: if we find a room that we have previously
// registered, then we save info about that room so that
// we can customize later interactions.  However, if those rooms
// were registered to use a different callback address from the
// one that we are currently using, then we will never be asked
// to service those rooms.  A less-simplistic approach would
// save only those rooms that had a callback address matching
// the callback address that we are currently registering;
// a more ambitious approach would be to start another go routine
// to handle each additional callback address that we discover
// in our list of previously-registered rooms.

// Game On!
// Main site: https://gameontext.org

// JSON Marshalling Notes
//
// 1. All marshalling is performed using json.MarshalIndent so
//    that any JSON we log is formatted nicely.
//
// 2. The additional tagging that we use to annotate our Go
//    struct types is explained in the Go json documentation, which
//    is currently at https://golang.org/pkg/encoding/json/ . You
//    should ready the documentation for the Marshall function
//    carefully.
//
// 3. Some JSON responses have more information than we
//    need, so we only define marshalling for the fields we
//    care about. See RoomRegistrationResp for an example
//    of a struct designed to keep a subset of the information
//    returned in a response. (We do not need exits or coords
//    so we do not map those fields.)

// Room commands
//
// The handling for each room command (/go, /look, etc.) is kept in
// its own source file and it is typically named room<cmd>.go, as in
// roomchat.go, roomlook.go, etc.

// Chat (broadcast messages)
//
// It is up to the room to implement the notion of global chat to
// all parties in a given room. This requires us to track player
// movement to/from rooms and to remember which websocket connection
// is currently associated with a given player in a given room. Tracker
// functions (tracker.go) handle this for us and they use channels to
// isolate shared data structures. BroadcastMessage(r, m, sender,
// receiver string) should be used to send broadcast messages.
// Messages to a player who is still in the room should be sent
// more directly using SendMessageToPlayer(conn *websocket.Conn,
// mUser, uid string).

// Deleting rooms
//
// At the time that this code was being written, Game On! did not
// expose a console command to delete a room. Since that initial
// time a console command has been added to delete a room. This
// room's delete code will remain, however, as an example of
// deleting a room from outside of Game On!
