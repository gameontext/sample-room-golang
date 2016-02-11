// Copyright (c) 2016 IBM Corp. All rights reserved.
// Use of this source code is governed by the Apache License,
// Version 2.0, a copy of which can be found in the LICENSE file.

// Documentation
package main

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
