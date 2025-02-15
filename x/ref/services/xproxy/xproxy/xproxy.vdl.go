// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file was auto-generated by the vanadium vdl tool.
// Package: xproxy

//nolint:revive
package xproxy

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
	ErrNotListening              = verror.NewIDAction("v.io/x/ref/services/xproxy/xproxy.NotListening", verror.NoRetry)
	ErrUnexpectedMessage         = verror.NewIDAction("v.io/x/ref/services/xproxy/xproxy.UnexpectedMessage", verror.NoRetry)
	ErrFailedToResolveToEndpoint = verror.NewIDAction("v.io/x/ref/services/xproxy/xproxy.FailedToResolveToEndpoint", verror.NoRetry)
	ErrProxyAlreadyClosed        = verror.NewIDAction("v.io/x/ref/services/xproxy/xproxy.ProxyAlreadyClosed", verror.NoRetry)
	ErrProxyResponse             = verror.NewIDAction("v.io/x/ref/services/xproxy/xproxy.ProxyResponse", verror.NoRetry)
)

// ErrorfNotListening calls ErrNotListening.Errorf with the supplied arguments.
func ErrorfNotListening(ctx *context.T, format string) error {
	return ErrNotListening.Errorf(ctx, format)
}

// MessageNotListening calls ErrNotListening.Message with the supplied arguments.
func MessageNotListening(ctx *context.T, message string) error {
	return ErrNotListening.Message(ctx, message)
}

// ParamsErrNotListening extracts the expected parameters from the error's ParameterList.
func ParamsErrNotListening(argumentError error) (verrorComponent string, verrorOperation string, returnErr error) {
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

// ErrorfUnexpectedMessage calls ErrUnexpectedMessage.Errorf with the supplied arguments.
func ErrorfUnexpectedMessage(ctx *context.T, format string, msgType string) error {
	return ErrUnexpectedMessage.Errorf(ctx, format, msgType)
}

// MessageUnexpectedMessage calls ErrUnexpectedMessage.Message with the supplied arguments.
func MessageUnexpectedMessage(ctx *context.T, message string, msgType string) error {
	return ErrUnexpectedMessage.Message(ctx, message, msgType)
}

// ParamsErrUnexpectedMessage extracts the expected parameters from the error's ParameterList.
func ParamsErrUnexpectedMessage(argumentError error) (verrorComponent string, verrorOperation string, msgType string, returnErr error) {
	params := verror.Params(argumentError)
	if params == nil {
		returnErr = fmt.Errorf("no parameters found in: %T: %v", argumentError, argumentError)
		return
	}
	iter := &paramListIterator{params: params, max: len(params)}

	if verrorComponent, verrorOperation, returnErr = iter.preamble(); returnErr != nil {
		return
	}

	var (
		tmp interface{}
		ok  bool
	)
	tmp, returnErr = iter.next()
	if msgType, ok = tmp.(string); !ok {
		if returnErr != nil {
			return
		}
		returnErr = fmt.Errorf("parameter list contains the wrong type for return value msgType, has %T and not string", tmp)
		return
	}

	return
}

// ErrorfFailedToResolveToEndpoint calls ErrFailedToResolveToEndpoint.Errorf with the supplied arguments.
func ErrorfFailedToResolveToEndpoint(ctx *context.T, format string, name string) error {
	return ErrFailedToResolveToEndpoint.Errorf(ctx, format, name)
}

// MessageFailedToResolveToEndpoint calls ErrFailedToResolveToEndpoint.Message with the supplied arguments.
func MessageFailedToResolveToEndpoint(ctx *context.T, message string, name string) error {
	return ErrFailedToResolveToEndpoint.Message(ctx, message, name)
}

// ParamsErrFailedToResolveToEndpoint extracts the expected parameters from the error's ParameterList.
func ParamsErrFailedToResolveToEndpoint(argumentError error) (verrorComponent string, verrorOperation string, name string, returnErr error) {
	params := verror.Params(argumentError)
	if params == nil {
		returnErr = fmt.Errorf("no parameters found in: %T: %v", argumentError, argumentError)
		return
	}
	iter := &paramListIterator{params: params, max: len(params)}

	if verrorComponent, verrorOperation, returnErr = iter.preamble(); returnErr != nil {
		return
	}

	var (
		tmp interface{}
		ok  bool
	)
	tmp, returnErr = iter.next()
	if name, ok = tmp.(string); !ok {
		if returnErr != nil {
			return
		}
		returnErr = fmt.Errorf("parameter list contains the wrong type for return value name, has %T and not string", tmp)
		return
	}

	return
}

// ErrorfProxyAlreadyClosed calls ErrProxyAlreadyClosed.Errorf with the supplied arguments.
func ErrorfProxyAlreadyClosed(ctx *context.T, format string) error {
	return ErrProxyAlreadyClosed.Errorf(ctx, format)
}

// MessageProxyAlreadyClosed calls ErrProxyAlreadyClosed.Message with the supplied arguments.
func MessageProxyAlreadyClosed(ctx *context.T, message string) error {
	return ErrProxyAlreadyClosed.Message(ctx, message)
}

// ParamsErrProxyAlreadyClosed extracts the expected parameters from the error's ParameterList.
func ParamsErrProxyAlreadyClosed(argumentError error) (verrorComponent string, verrorOperation string, returnErr error) {
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

// ErrorfProxyResponse calls ErrProxyResponse.Errorf with the supplied arguments.
func ErrorfProxyResponse(ctx *context.T, format string, msg string) error {
	return ErrProxyResponse.Errorf(ctx, format, msg)
}

// MessageProxyResponse calls ErrProxyResponse.Message with the supplied arguments.
func MessageProxyResponse(ctx *context.T, message string, msg string) error {
	return ErrProxyResponse.Message(ctx, message, msg)
}

// ParamsErrProxyResponse extracts the expected parameters from the error's ParameterList.
func ParamsErrProxyResponse(argumentError error) (verrorComponent string, verrorOperation string, msg string, returnErr error) {
	params := verror.Params(argumentError)
	if params == nil {
		returnErr = fmt.Errorf("no parameters found in: %T: %v", argumentError, argumentError)
		return
	}
	iter := &paramListIterator{params: params, max: len(params)}

	if verrorComponent, verrorOperation, returnErr = iter.preamble(); returnErr != nil {
		return
	}

	var (
		tmp interface{}
		ok  bool
	)
	tmp, returnErr = iter.next()
	if msg, ok = tmp.(string); !ok {
		if returnErr != nil {
			return
		}
		returnErr = fmt.Errorf("parameter list contains the wrong type for return value msg, has %T and not string", tmp)
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
