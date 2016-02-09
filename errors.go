// Copyright (c) 2016 IBM Corp. All rights reserved.
// Use of this source code is governed by the Apache License,
// Version 2.0, a copy of which can be found in the LICENSE file.

// All of our errors are defined in this file
package main

import (
	"fmt"
)

// RegError describes a registration error.
type RegError struct {
	message string
}

func (e RegError) Error() string { return fmt.Sprintf("REG.ERROR: %s", e.message) }

// Payload error describes an error concerning the top-level structure of an
// incoming webservice payload. Payload errors do not deal with the contents
// of the javascript payload component
type PayloadError struct {
	message string
}

func (e PayloadError) Error() string { return fmt.Sprintf("PAYLOAD.ERROR: %s", e.message) }

// JSPayloadError describes a problem with the structure of the
// javascript payload component.
type JSPayloadError struct {
	message string
}

func (e JSPayloadError) Error() string { return fmt.Sprintf("JS.PAYLOAD.ERROR: %s", e.message) }

// ArgError describes an error with the commandline arguments.
type ArgError struct {
	message string
}

func (e ArgError) Error() string { return fmt.Sprintf("ARG.ERROR: %s", e.message) }
