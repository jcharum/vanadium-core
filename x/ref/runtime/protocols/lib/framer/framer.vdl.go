// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file was auto-generated by the vanadium vdl tool.
// Package: framer

//nolint:revive
package framer

import (
	"fmt"

	"v.io/v23/context"
	"v.io/v23/verror"
)

var initializeVDLCalled = false
var _ = initializeVDL() // Must be first; see initializeVDL comments for details.

// Error definitions
// =================

var (
	ErrLargerThan3ByteUInt = verror.NewIDAction("v.io/x/ref/runtime/protocols/lib/framer.LargerThan3ByteUInt", verror.NoRetry)
)

// ErrorfLargerThan3ByteUInt calls ErrLargerThan3ByteUInt.Errorf with the supplied arguments.
func ErrorfLargerThan3ByteUInt(ctx *context.T, format string) error {
	return ErrLargerThan3ByteUInt.Errorf(ctx, format)
}

// MessageLargerThan3ByteUInt calls ErrLargerThan3ByteUInt.Message with the supplied arguments.
func MessageLargerThan3ByteUInt(ctx *context.T, message string) error {
	return ErrLargerThan3ByteUInt.Message(ctx, message)
}

// ParamsErrLargerThan3ByteUInt extracts the expected parameters from the error's ParameterList.
func ParamsErrLargerThan3ByteUInt(argumentError error) (verrorComponent string, verrorOperation string, returnErr error) {
	params := verror.Params(argumentError)
	if params == nil {
		returnErr = fmt.Errorf("no parameters found in: %T: %v", argumentError, argumentError)
		return
	}
	iter := &paramListIterator{params: params, max: len(params)}

	if verrorComponent, verrorOperation, returnErr = iter.preamble(); returnErr != nil {
		return
	}

	return
}

type paramListIterator struct {
	err      error
	idx, max int
	params   []interface{}
}

func (pl *paramListIterator) next() (interface{}, error) {
	if pl.err != nil {
		return nil, pl.err
	}
	if pl.idx+1 > pl.max {
		pl.err = fmt.Errorf("too few parameters: have %v", pl.max)
		return nil, pl.err
	}
	pl.idx++
	return pl.params[pl.idx-1], nil
}

func (pl *paramListIterator) preamble() (component, operation string, err error) {
	var tmp interface{}
	if tmp, err = pl.next(); err != nil {
		return
	}
	var ok bool
	if component, ok = tmp.(string); !ok {
		return "", "", fmt.Errorf("ParamList[0]: component name is not a string: %T", tmp)
	}
	if tmp, err = pl.next(); err != nil {
		return
	}
	if operation, ok = tmp.(string); !ok {
		return "", "", fmt.Errorf("ParamList[1]: operation name is not a string: %T", tmp)
	}
	return
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

	return struct{}{}
}
