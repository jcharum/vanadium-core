// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file was auto-generated by the vanadium vdl tool.
// Package: uniqueid

//nolint:revive
package uniqueid

import (
	"v.io/v23/vdl"
)

var initializeVDLCalled = false
var _ = initializeVDL() // Must be first; see initializeVDL comments for details.

// Hold type definitions in package-level variables, for better performance.
// Declare and initialize with default values here so that the initializeVDL
// method will be considered ready to initialize before any of the type
// definitions that appear below.
//nolint:unused
var (
	vdlTypeArray1 *vdl.Type = nil
)

// Type definitions
// ================
// An Id is a likely globally unique identifier.
type Id [16]byte

func (Id) VDLReflect(struct {
	Name string `vdl:"v.io/v23/uniqueid.Id"`
}) {
}

func (x Id) VDLIsZero() bool { //nolint:gocyclo
	return x == Id{}
}

func (x Id) VDLWrite(enc vdl.Encoder) error { //nolint:gocyclo
	if err := enc.WriteValueBytes(vdlTypeArray1, x[:]); err != nil {
		return err
	}
	return nil
}

func (x *Id) VDLRead(dec vdl.Decoder) error { //nolint:gocyclo
	bytes := x[:]
	if err := dec.ReadValueBytes(16, &bytes); err != nil {
		return err
	}
	return nil
}

// initializeVDL performs vdl initialization.  It is safe to call multiple times.
// If you have an init ordering issue, just insert the following line verbatim
// into your source files in this package, right after the "package foo" clause:
//
//    var _ = initializeVDL()
//
// The purpose of this function is to ensure that vdl initialization occurs in
// the right order, and very early in the init sequence.  In particular, vdl
// registration and package variable initialization needs to occur before
// functions like vdl.TypeOf will work properly.
//
// This function returns a dummy value, so that it can be used to initialize the
// first var in the file, to take advantage of Go's defined init order.
func initializeVDL() struct{} {
	if initializeVDLCalled {
		return struct{}{}
	}
	initializeVDLCalled = true

	// Register types.
	vdl.Register((*Id)(nil))

	// Initialize type definitions.
	vdlTypeArray1 = vdl.TypeOf((*Id)(nil))

	return struct{}{}
}
