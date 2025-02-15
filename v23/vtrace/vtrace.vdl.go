// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file was auto-generated by the vanadium vdl tool.
// Package: vtrace

//nolint:revive
package vtrace

import (
	"time"

	"v.io/v23/uniqueid"
	"v.io/v23/vdl"
	vdltime "v.io/v23/vdlroot/time"
)

var initializeVDLCalled = false
var _ = initializeVDL() // Must be first; see initializeVDL comments for details.

// Hold type definitions in package-level variables, for better performance.
// Declare and initialize with default values here so that the initializeVDL
// method will be considered ready to initialize before any of the type
// definitions that appear below.
//nolint:unused
var (
	vdlTypeStruct1  *vdl.Type = nil
	vdlTypeStruct2  *vdl.Type = nil
	vdlTypeStruct3  *vdl.Type = nil
	vdlTypeArray4   *vdl.Type = nil
	vdlTypeList5    *vdl.Type = nil
	vdlTypeList6    *vdl.Type = nil
	vdlTypeStruct7  *vdl.Type = nil
	vdlTypeList8    *vdl.Type = nil
	vdlTypeInt329   *vdl.Type = nil
	vdlTypeStruct10 *vdl.Type = nil
	vdlTypeStruct11 *vdl.Type = nil
)

// Type definitions
// ================
// An Annotation represents data that is relevant at a specific moment.
// They can be attached to spans to add useful debugging information.
type Annotation struct {
	// When the annotation was added.
	When time.Time
	// The annotation message.
	// TODO(mattr): Allow richer annotations.
	Message string
}

func (Annotation) VDLReflect(struct {
	Name string `vdl:"v.io/v23/vtrace.Annotation"`
}) {
}

func (x Annotation) VDLIsZero() bool { //nolint:gocyclo
	if !x.When.IsZero() {
		return false
	}
	if x.Message != "" {
		return false
	}
	return true
}

func (x Annotation) VDLWrite(enc vdl.Encoder) error { //nolint:gocyclo
	if err := enc.StartValue(vdlTypeStruct1); err != nil {
		return err
	}
	if !x.When.IsZero() {
		if err := enc.NextField(0); err != nil {
			return err
		}
		var wire vdltime.Time
		if err := vdltime.TimeFromNative(&wire, x.When); err != nil {
			return err
		}
		if err := wire.VDLWrite(enc); err != nil {
			return err
		}
	}
	if x.Message != "" {
		if err := enc.NextFieldValueString(1, vdl.StringType, x.Message); err != nil {
			return err
		}
	}
	if err := enc.NextField(-1); err != nil {
		return err
	}
	return enc.FinishValue()
}

func (x *Annotation) VDLRead(dec vdl.Decoder) error { //nolint:gocyclo
	*x = Annotation{}
	if err := dec.StartValue(vdlTypeStruct1); err != nil {
		return err
	}
	decType := dec.Type()
	for {
		index, err := dec.NextField()
		switch {
		case err != nil:
			return err
		case index == -1:
			return dec.FinishValue()
		}
		if decType != vdlTypeStruct1 {
			index = vdlTypeStruct1.FieldIndexByName(decType.Field(index).Name)
			if index == -1 {
				if err := dec.SkipValue(); err != nil {
					return err
				}
				continue
			}
		}
		switch index {
		case 0:
			var wire vdltime.Time
			if err := wire.VDLRead(dec); err != nil {
				return err
			}
			if err := vdltime.TimeToNative(wire, &x.When); err != nil {
				return err
			}
		case 1:
			switch value, err := dec.ReadValueString(); {
			case err != nil:
				return err
			default:
				x.Message = value
			}
		}
	}
}

// A SpanRecord is the wire format for a Span.
type SpanRecord struct {
	Id     uniqueid.Id // The Id of the Span.
	Parent uniqueid.Id // The Id of this Span's parent.
	Name   string      // The Name of this span.
	Start  time.Time   // The start time of this span.
	End    time.Time   // The end time of this span.
	// A series of annotations.
	Annotations []Annotation
	// RequestMetadata that will be sent along with the request.
	RequestMetadata []byte
}

func (SpanRecord) VDLReflect(struct {
	Name string `vdl:"v.io/v23/vtrace.SpanRecord"`
}) {
}

func (x SpanRecord) VDLIsZero() bool { //nolint:gocyclo
	if x.Id != (uniqueid.Id{}) {
		return false
	}
	if x.Parent != (uniqueid.Id{}) {
		return false
	}
	if x.Name != "" {
		return false
	}
	if !x.Start.IsZero() {
		return false
	}
	if !x.End.IsZero() {
		return false
	}
	if len(x.Annotations) != 0 {
		return false
	}
	if len(x.RequestMetadata) != 0 {
		return false
	}
	return true
}

func (x SpanRecord) VDLWrite(enc vdl.Encoder) error { //nolint:gocyclo
	if err := enc.StartValue(vdlTypeStruct3); err != nil {
		return err
	}
	if x.Id != (uniqueid.Id{}) {
		if err := enc.NextFieldValueBytes(0, vdlTypeArray4, x.Id[:]); err != nil {
			return err
		}
	}
	if x.Parent != (uniqueid.Id{}) {
		if err := enc.NextFieldValueBytes(1, vdlTypeArray4, x.Parent[:]); err != nil {
			return err
		}
	}
	if x.Name != "" {
		if err := enc.NextFieldValueString(2, vdl.StringType, x.Name); err != nil {
			return err
		}
	}
	if !x.Start.IsZero() {
		if err := enc.NextField(3); err != nil {
			return err
		}
		var wire vdltime.Time
		if err := vdltime.TimeFromNative(&wire, x.Start); err != nil {
			return err
		}
		if err := wire.VDLWrite(enc); err != nil {
			return err
		}
	}
	if !x.End.IsZero() {
		if err := enc.NextField(4); err != nil {
			return err
		}
		var wire vdltime.Time
		if err := vdltime.TimeFromNative(&wire, x.End); err != nil {
			return err
		}
		if err := wire.VDLWrite(enc); err != nil {
			return err
		}
	}
	if len(x.Annotations) != 0 {
		if err := enc.NextField(5); err != nil {
			return err
		}
		if err := vdlWriteAnonList1(enc, x.Annotations); err != nil {
			return err
		}
	}
	if len(x.RequestMetadata) != 0 {
		if err := enc.NextFieldValueBytes(6, vdlTypeList6, x.RequestMetadata); err != nil {
			return err
		}
	}
	if err := enc.NextField(-1); err != nil {
		return err
	}
	return enc.FinishValue()
}

func vdlWriteAnonList1(enc vdl.Encoder, x []Annotation) error {
	if err := enc.StartValue(vdlTypeList5); err != nil {
		return err
	}
	if err := enc.SetLenHint(len(x)); err != nil {
		return err
	}
	for _, elem := range x {
		if err := enc.NextEntry(false); err != nil {
			return err
		}
		if err := elem.VDLWrite(enc); err != nil {
			return err
		}
	}
	if err := enc.NextEntry(true); err != nil {
		return err
	}
	return enc.FinishValue()
}

func (x *SpanRecord) VDLRead(dec vdl.Decoder) error { //nolint:gocyclo
	*x = SpanRecord{}
	if err := dec.StartValue(vdlTypeStruct3); err != nil {
		return err
	}
	decType := dec.Type()
	for {
		index, err := dec.NextField()
		switch {
		case err != nil:
			return err
		case index == -1:
			return dec.FinishValue()
		}
		if decType != vdlTypeStruct3 {
			index = vdlTypeStruct3.FieldIndexByName(decType.Field(index).Name)
			if index == -1 {
				if err := dec.SkipValue(); err != nil {
					return err
				}
				continue
			}
		}
		switch index {
		case 0:
			bytes := x.Id[:]
			if err := dec.ReadValueBytes(16, &bytes); err != nil {
				return err
			}
		case 1:
			bytes := x.Parent[:]
			if err := dec.ReadValueBytes(16, &bytes); err != nil {
				return err
			}
		case 2:
			switch value, err := dec.ReadValueString(); {
			case err != nil:
				return err
			default:
				x.Name = value
			}
		case 3:
			var wire vdltime.Time
			if err := wire.VDLRead(dec); err != nil {
				return err
			}
			if err := vdltime.TimeToNative(wire, &x.Start); err != nil {
				return err
			}
		case 4:
			var wire vdltime.Time
			if err := wire.VDLRead(dec); err != nil {
				return err
			}
			if err := vdltime.TimeToNative(wire, &x.End); err != nil {
				return err
			}
		case 5:
			if err := vdlReadAnonList1(dec, &x.Annotations); err != nil {
				return err
			}
		case 6:
			if err := dec.ReadValueBytes(-1, &x.RequestMetadata); err != nil {
				return err
			}
		}
	}
}

func vdlReadAnonList1(dec vdl.Decoder, x *[]Annotation) error {
	if err := dec.StartValue(vdlTypeList5); err != nil {
		return err
	}
	if len := dec.LenHint(); len > 0 {
		*x = make([]Annotation, 0, len)
	} else {
		*x = nil
	}
	for {
		switch done, err := dec.NextEntry(); {
		case err != nil:
			return err
		case done:
			return dec.FinishValue()
		default:
			var elem Annotation
			if err := elem.VDLRead(dec); err != nil {
				return err
			}
			*x = append(*x, elem)
		}
	}
}

type TraceRecord struct {
	Id    uniqueid.Id
	Spans []SpanRecord
}

func (TraceRecord) VDLReflect(struct {
	Name string `vdl:"v.io/v23/vtrace.TraceRecord"`
}) {
}

func (x TraceRecord) VDLIsZero() bool { //nolint:gocyclo
	if x.Id != (uniqueid.Id{}) {
		return false
	}
	if len(x.Spans) != 0 {
		return false
	}
	return true
}

func (x TraceRecord) VDLWrite(enc vdl.Encoder) error { //nolint:gocyclo
	if err := enc.StartValue(vdlTypeStruct7); err != nil {
		return err
	}
	if x.Id != (uniqueid.Id{}) {
		if err := enc.NextFieldValueBytes(0, vdlTypeArray4, x.Id[:]); err != nil {
			return err
		}
	}
	if len(x.Spans) != 0 {
		if err := enc.NextField(1); err != nil {
			return err
		}
		if err := vdlWriteAnonList2(enc, x.Spans); err != nil {
			return err
		}
	}
	if err := enc.NextField(-1); err != nil {
		return err
	}
	return enc.FinishValue()
}

func vdlWriteAnonList2(enc vdl.Encoder, x []SpanRecord) error {
	if err := enc.StartValue(vdlTypeList8); err != nil {
		return err
	}
	if err := enc.SetLenHint(len(x)); err != nil {
		return err
	}
	for _, elem := range x {
		if err := enc.NextEntry(false); err != nil {
			return err
		}
		if err := elem.VDLWrite(enc); err != nil {
			return err
		}
	}
	if err := enc.NextEntry(true); err != nil {
		return err
	}
	return enc.FinishValue()
}

func (x *TraceRecord) VDLRead(dec vdl.Decoder) error { //nolint:gocyclo
	*x = TraceRecord{}
	if err := dec.StartValue(vdlTypeStruct7); err != nil {
		return err
	}
	decType := dec.Type()
	for {
		index, err := dec.NextField()
		switch {
		case err != nil:
			return err
		case index == -1:
			return dec.FinishValue()
		}
		if decType != vdlTypeStruct7 {
			index = vdlTypeStruct7.FieldIndexByName(decType.Field(index).Name)
			if index == -1 {
				if err := dec.SkipValue(); err != nil {
					return err
				}
				continue
			}
		}
		switch index {
		case 0:
			bytes := x.Id[:]
			if err := dec.ReadValueBytes(16, &bytes); err != nil {
				return err
			}
		case 1:
			if err := vdlReadAnonList2(dec, &x.Spans); err != nil {
				return err
			}
		}
	}
}

func vdlReadAnonList2(dec vdl.Decoder, x *[]SpanRecord) error {
	if err := dec.StartValue(vdlTypeList8); err != nil {
		return err
	}
	if len := dec.LenHint(); len > 0 {
		*x = make([]SpanRecord, 0, len)
	} else {
		*x = nil
	}
	for {
		switch done, err := dec.NextEntry(); {
		case err != nil:
			return err
		case done:
			return dec.FinishValue()
		default:
			var elem SpanRecord
			if err := elem.VDLRead(dec); err != nil {
				return err
			}
			*x = append(*x, elem)
		}
	}
}

// TraceFlags represents a bit mask that determines
type TraceFlags int32

func (TraceFlags) VDLReflect(struct {
	Name string `vdl:"v.io/v23/vtrace.TraceFlags"`
}) {
}

func (x TraceFlags) VDLIsZero() bool { //nolint:gocyclo
	return x == 0
}

func (x TraceFlags) VDLWrite(enc vdl.Encoder) error { //nolint:gocyclo
	if err := enc.WriteValueInt(vdlTypeInt329, int64(x)); err != nil {
		return err
	}
	return nil
}

func (x *TraceFlags) VDLRead(dec vdl.Decoder) error { //nolint:gocyclo
	switch value, err := dec.ReadValueInt(32); {
	case err != nil:
		return err
	default:
		*x = TraceFlags(value)
	}
	return nil
}

// Request is the object that carries trace information between processes.
type Request struct {
	SpanId          uniqueid.Id // The Id of the span that originated the RPC call.
	TraceId         uniqueid.Id // The Id of the trace this call is a part of.
	RequestMetadata []byte      // Any metadata to be sent with the request.
	Flags           TraceFlags
	LogLevel        int32
}

func (Request) VDLReflect(struct {
	Name string `vdl:"v.io/v23/vtrace.Request"`
}) {
}

func (x Request) VDLIsZero() bool { //nolint:gocyclo
	if x.SpanId != (uniqueid.Id{}) {
		return false
	}
	if x.TraceId != (uniqueid.Id{}) {
		return false
	}
	if len(x.RequestMetadata) != 0 {
		return false
	}
	if x.Flags != 0 {
		return false
	}
	if x.LogLevel != 0 {
		return false
	}
	return true
}

func (x Request) VDLWrite(enc vdl.Encoder) error { //nolint:gocyclo
	if err := enc.StartValue(vdlTypeStruct10); err != nil {
		return err
	}
	if x.SpanId != (uniqueid.Id{}) {
		if err := enc.NextFieldValueBytes(0, vdlTypeArray4, x.SpanId[:]); err != nil {
			return err
		}
	}
	if x.TraceId != (uniqueid.Id{}) {
		if err := enc.NextFieldValueBytes(1, vdlTypeArray4, x.TraceId[:]); err != nil {
			return err
		}
	}
	if len(x.RequestMetadata) != 0 {
		if err := enc.NextFieldValueBytes(2, vdlTypeList6, x.RequestMetadata); err != nil {
			return err
		}
	}
	if x.Flags != 0 {
		if err := enc.NextFieldValueInt(3, vdlTypeInt329, int64(x.Flags)); err != nil {
			return err
		}
	}
	if x.LogLevel != 0 {
		if err := enc.NextFieldValueInt(4, vdl.Int32Type, int64(x.LogLevel)); err != nil {
			return err
		}
	}
	if err := enc.NextField(-1); err != nil {
		return err
	}
	return enc.FinishValue()
}

func (x *Request) VDLRead(dec vdl.Decoder) error { //nolint:gocyclo
	*x = Request{}
	if err := dec.StartValue(vdlTypeStruct10); err != nil {
		return err
	}
	decType := dec.Type()
	for {
		index, err := dec.NextField()
		switch {
		case err != nil:
			return err
		case index == -1:
			return dec.FinishValue()
		}
		if decType != vdlTypeStruct10 {
			index = vdlTypeStruct10.FieldIndexByName(decType.Field(index).Name)
			if index == -1 {
				if err := dec.SkipValue(); err != nil {
					return err
				}
				continue
			}
		}
		switch index {
		case 0:
			bytes := x.SpanId[:]
			if err := dec.ReadValueBytes(16, &bytes); err != nil {
				return err
			}
		case 1:
			bytes := x.TraceId[:]
			if err := dec.ReadValueBytes(16, &bytes); err != nil {
				return err
			}
		case 2:
			if err := dec.ReadValueBytes(-1, &x.RequestMetadata); err != nil {
				return err
			}
		case 3:
			switch value, err := dec.ReadValueInt(32); {
			case err != nil:
				return err
			default:
				x.Flags = TraceFlags(value)
			}
		case 4:
			switch value, err := dec.ReadValueInt(32); {
			case err != nil:
				return err
			default:
				x.LogLevel = int32(value)
			}
		}
	}
}

type Response struct {
	// Flags give options for trace collection, the client should alter its
	// collection for this trace according to the flags sent back from the
	// server.
	Flags TraceFlags
	// Trace is collected trace data.  This may be empty.
	Trace TraceRecord
}

func (Response) VDLReflect(struct {
	Name string `vdl:"v.io/v23/vtrace.Response"`
}) {
}

func (x Response) VDLIsZero() bool { //nolint:gocyclo
	if x.Flags != 0 {
		return false
	}
	if !x.Trace.VDLIsZero() {
		return false
	}
	return true
}

func (x Response) VDLWrite(enc vdl.Encoder) error { //nolint:gocyclo
	if err := enc.StartValue(vdlTypeStruct11); err != nil {
		return err
	}
	if x.Flags != 0 {
		if err := enc.NextFieldValueInt(0, vdlTypeInt329, int64(x.Flags)); err != nil {
			return err
		}
	}
	if !x.Trace.VDLIsZero() {
		if err := enc.NextField(1); err != nil {
			return err
		}
		if err := x.Trace.VDLWrite(enc); err != nil {
			return err
		}
	}
	if err := enc.NextField(-1); err != nil {
		return err
	}
	return enc.FinishValue()
}

func (x *Response) VDLRead(dec vdl.Decoder) error { //nolint:gocyclo
	*x = Response{}
	if err := dec.StartValue(vdlTypeStruct11); err != nil {
		return err
	}
	decType := dec.Type()
	for {
		index, err := dec.NextField()
		switch {
		case err != nil:
			return err
		case index == -1:
			return dec.FinishValue()
		}
		if decType != vdlTypeStruct11 {
			index = vdlTypeStruct11.FieldIndexByName(decType.Field(index).Name)
			if index == -1 {
				if err := dec.SkipValue(); err != nil {
					return err
				}
				continue
			}
		}
		switch index {
		case 0:
			switch value, err := dec.ReadValueInt(32); {
			case err != nil:
				return err
			default:
				x.Flags = TraceFlags(value)
			}
		case 1:
			if err := x.Trace.VDLRead(dec); err != nil {
				return err
			}
		}
	}
}

// Const definitions
// =================

const Empty = TraceFlags(0)
const CollectInMemory = TraceFlags(1)
const AWSXRay = TraceFlags(2)

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
	vdl.Register((*Annotation)(nil))
	vdl.Register((*SpanRecord)(nil))
	vdl.Register((*TraceRecord)(nil))
	vdl.Register((*TraceFlags)(nil))
	vdl.Register((*Request)(nil))
	vdl.Register((*Response)(nil))

	// Initialize type definitions.
	vdlTypeStruct1 = vdl.TypeOf((*Annotation)(nil)).Elem()
	vdlTypeStruct2 = vdl.TypeOf((*vdltime.Time)(nil)).Elem()
	vdlTypeStruct3 = vdl.TypeOf((*SpanRecord)(nil)).Elem()
	vdlTypeArray4 = vdl.TypeOf((*uniqueid.Id)(nil))
	vdlTypeList5 = vdl.TypeOf((*[]Annotation)(nil))
	vdlTypeList6 = vdl.TypeOf((*[]byte)(nil))
	vdlTypeStruct7 = vdl.TypeOf((*TraceRecord)(nil)).Elem()
	vdlTypeList8 = vdl.TypeOf((*[]SpanRecord)(nil))
	vdlTypeInt329 = vdl.TypeOf((*TraceFlags)(nil))
	vdlTypeStruct10 = vdl.TypeOf((*Request)(nil)).Elem()
	vdlTypeStruct11 = vdl.TypeOf((*Response)(nil)).Elem()

	return struct{}{}
}
