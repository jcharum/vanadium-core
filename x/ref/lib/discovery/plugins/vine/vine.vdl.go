// Copyright 2016 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file was auto-generated by the vanadium vdl tool.
// Package: vine

//nolint:revive
package vine

import (
	"time"

	v23 "v.io/v23"
	"v.io/v23/context"
	"v.io/v23/discovery"
	"v.io/v23/rpc"
	"v.io/v23/security/access"
	"v.io/v23/vdl"
	_ "v.io/v23/vdlroot/time"
	discovery_2 "v.io/x/ref/lib/discovery"
)

var initializeVDLCalled = false
var _ = initializeVDL() // Must be first; see initializeVDL comments for details.

// Interface definitions
// =====================

// StoreClientMethods is the client interface
// containing Store methods.
//
// Store is the interface for sharing advertisements between vine plugins.
type StoreClientMethods interface {
	// Add adds an advertisement with a given ttl.
	Add(_ *context.T, adinfo discovery_2.AdInfo, ttl time.Duration, _ ...rpc.CallOpt) error
	// Delete deletes the advertisement from the store.
	Delete(_ *context.T, id discovery.AdId, _ ...rpc.CallOpt) error
}

// StoreClientStub embeds StoreClientMethods and is a
// placeholder for additional management operations.
type StoreClientStub interface {
	StoreClientMethods
}

// StoreClient returns a client stub for Store.
func StoreClient(name string) StoreClientStub {
	return implStoreClientStub{name}
}

type implStoreClientStub struct {
	name string
}

func (c implStoreClientStub) Add(ctx *context.T, i0 discovery_2.AdInfo, i1 time.Duration, opts ...rpc.CallOpt) (err error) {
	err = v23.GetClient(ctx).Call(ctx, c.name, "Add", []interface{}{i0, i1}, nil, opts...)
	return
}

func (c implStoreClientStub) Delete(ctx *context.T, i0 discovery.AdId, opts ...rpc.CallOpt) (err error) {
	err = v23.GetClient(ctx).Call(ctx, c.name, "Delete", []interface{}{i0}, nil, opts...)
	return
}

// StoreServerMethods is the interface a server writer
// implements for Store.
//
// Store is the interface for sharing advertisements between vine plugins.
type StoreServerMethods interface {
	// Add adds an advertisement with a given ttl.
	Add(_ *context.T, _ rpc.ServerCall, adinfo discovery_2.AdInfo, ttl time.Duration) error
	// Delete deletes the advertisement from the store.
	Delete(_ *context.T, _ rpc.ServerCall, id discovery.AdId) error
}

// StoreServerStubMethods is the server interface containing
// Store methods, as expected by rpc.Server.
// There is no difference between this interface and StoreServerMethods
// since there are no streaming methods.
type StoreServerStubMethods StoreServerMethods

// StoreServerStub adds universal methods to StoreServerStubMethods.
type StoreServerStub interface {
	StoreServerStubMethods
	// DescribeInterfaces the Store interfaces.
	Describe__() []rpc.InterfaceDesc
}

// StoreServer returns a server stub for Store.
// It converts an implementation of StoreServerMethods into
// an object that may be used by rpc.Server.
func StoreServer(impl StoreServerMethods) StoreServerStub {
	stub := implStoreServerStub{
		impl: impl,
	}
	// Initialize GlobState; always check the stub itself first, to handle the
	// case where the user has the Glob method defined in their VDL source.
	if gs := rpc.NewGlobState(stub); gs != nil {
		stub.gs = gs
	} else if gs := rpc.NewGlobState(impl); gs != nil {
		stub.gs = gs
	}
	return stub
}

type implStoreServerStub struct {
	impl StoreServerMethods
	gs   *rpc.GlobState
}

func (s implStoreServerStub) Add(ctx *context.T, call rpc.ServerCall, i0 discovery_2.AdInfo, i1 time.Duration) error {
	return s.impl.Add(ctx, call, i0, i1)
}

func (s implStoreServerStub) Delete(ctx *context.T, call rpc.ServerCall, i0 discovery.AdId) error {
	return s.impl.Delete(ctx, call, i0)
}

func (s implStoreServerStub) Globber() *rpc.GlobState {
	return s.gs
}

func (s implStoreServerStub) Describe__() []rpc.InterfaceDesc {
	return []rpc.InterfaceDesc{StoreDesc}
}

// StoreDesc describes the Store interface.
var StoreDesc rpc.InterfaceDesc = descStore

// descStore hides the desc to keep godoc clean.
var descStore = rpc.InterfaceDesc{
	Name:    "Store",
	PkgPath: "v.io/x/ref/lib/discovery/plugins/vine",
	Doc:     "// Store is the interface for sharing advertisements between vine plugins.",
	Methods: []rpc.MethodDesc{
		{
			Name: "Add",
			Doc:  "// Add adds an advertisement with a given ttl.",
			InArgs: []rpc.ArgDesc{
				{Name: "adinfo", Doc: ``}, // discovery_2.AdInfo
				{Name: "ttl", Doc: ``},    // time.Duration
			},
			Tags: []*vdl.Value{vdl.ValueOf(access.Tag("Write"))},
		},
		{
			Name: "Delete",
			Doc:  "// Delete deletes the advertisement from the store.",
			InArgs: []rpc.ArgDesc{
				{Name: "id", Doc: ``}, // discovery.AdId
			},
			Tags: []*vdl.Value{vdl.ValueOf(access.Tag("Write"))},
		},
	},
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
