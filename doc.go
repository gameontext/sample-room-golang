// Copyright (c) 2016 IBM Corp. All rights reserved.
// Use of this source code is governed by the Apache License,
// Version 2.0, a copy of which can be found in the LICENSE file.

// Documentation
package main

// JSON Marshalling Notes
//
// All marshalling is performed using json.MarshalIndent so
// that any JSON we log is formatted nicely.
// The additional tagging that we use to annotate our Go
// struct types is explained in the Go json documentation, which
// is currently at https://golang.org/pkg/encoding/json/ . You
// should ready the documentation for the Marshall function
// carefully.
