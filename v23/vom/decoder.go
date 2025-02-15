// Copyright 2016 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vom

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"

	"v.io/v23/vdl"
)

const (
	// IEEE 754 represents float64 using 52 bits to represent the mantissa, with
	// an extra implied leading bit.  That gives us 53 bits to store integers
	// without overflow - i.e. [0, (2^53)-1].  And since 2^53 is a small power of
	// two, it can also be stored without loss via mantissa=1 exponent=53.  Thus
	// we have our max and min values.  Ditto for float32, which uses 23 bits with
	// an extra implied leading bit.
	float64MaxInt = (1 << 53)
	float64MinInt = -(1 << 53)
	float32MaxInt = (1 << 24)
	float32MinInt = -(1 << 24)
)

var (
	errEmptyDecoderStack          = errors.New("vom: empty decoder stack")
	errReadRawBytesAlreadyStarted = errors.New("vom: read into vom.RawBytes after StartValue called")
	errReadRawBytesFromNonAny     = errors.New("vom: read into vom.RawBytes only supported on any values")
)

// This is only used for debugging; add this as the first line of NewDecoder to
// dump formatted vom bytes to stdout:
//   r = teeDump(r)
//nolint:unused
func teeDump(r io.Reader) io.Reader {
	return io.TeeReader(r, NewDumper(NewDumpWriter(os.Stdout)))
}

// Decoder manages the receipt and unmarshalling of typed values from the other
// side of a connection.
type Decoder struct {
	dec decoder81
}

const reservedCodec81StackSize = 2

type decoder81 struct {
	buf          *decbuf
	flag         decFlag
	stackStorage [reservedCodec81StackSize]decStackEntry
	stack        []decStackEntry
	refTypes     referencedTypes
	refAnyLens   referencedAnyLens
	typeDec      *TypeDecoder
}

type decStackEntry struct {
	Type          *vdl.Type     // Type of the value that we're decoding.
	Index         int           // Index of the key, elem or field.
	LenHint       int           // Length of the value, or -1 for unknown.
	NextEntryType *vdl.Type     // type for NextEntryValue* methods
	Flag          decStackFlag  // properties of this stack entry
	NextEntryData nextEntryData // properties of the next entry
}

// decFlag holds properties of the decoder.
type decFlag uint

const (
	decFlagIgnoreNextStartValue decFlag = 0x1 // ignore the next call to StartValue
	decFlagIsParentBytes        decFlag = 0x2 // parent type is []byte or [N]byte
	decFlagTypeIncomplete       decFlag = 0x4 // type has dependencies on unsent types
	decFlagSeparateTypeDec      decFlag = 0x8 // type decoder is separate

	// In FinishValue we need to clear both of these bits.
	decFlagFinishValue decFlag = decFlagIgnoreNextStartValue | decFlagIsParentBytes
)

func (f decFlag) Set(bits decFlag) decFlag   { return f | bits }
func (f decFlag) Clear(bits decFlag) decFlag { return f &^ bits }

func (f decFlag) IgnoreNextStartValue() bool { return f&decFlagIgnoreNextStartValue != 0 }
func (f decFlag) IsParentBytes() bool        { return f&decFlagIsParentBytes != 0 }
func (f decFlag) TypeIncomplete() bool       { return f&decFlagTypeIncomplete != 0 }
func (f decFlag) SeparateTypeDec() bool      { return f&decFlagSeparateTypeDec != 0 }

// decStackFlag holds type or value properties of the stack entry.
type decStackFlag uint

const (
	decStackFlagIsMapKey   decStackFlag = 0x1 // key or elem for dfsNextType
	decStackFlagIsAny      decStackFlag = 0x2 // the static type is Any
	decStackFlagIsOptional decStackFlag = 0x4 // the static type is Optional
	decStackFlagFastRead   decStackFlag = 0x8 // subtypes use ReadValue fastpath
)

func (f decStackFlag) FlipIsMapKey() decStackFlag { return f ^ decStackFlagIsMapKey }
func (f decStackFlag) IsMapKey() bool             { return f&decStackFlagIsMapKey != 0 }
func (f decStackFlag) IsAny() bool                { return f&decStackFlagIsAny != 0 }
func (f decStackFlag) IsOptional() bool           { return f&decStackFlagIsOptional != 0 }
func (f decStackFlag) FastRead() bool             { return f&decStackFlagFastRead != 0 }

// NewDecoder returns a new Decoder that reads from the given reader.  The
// Decoder understands all formats generated by the Encoder.
func NewDecoder(r io.Reader) *Decoder {
	buf := newDecbuf(r)
	dec := &decoder81{buf: buf}
	typeDec := newTypeDecoderInternal(dec)
	return &Decoder{decoder81{
		buf:     buf,
		typeDec: typeDec,
	}}
}

// NewDecoderWithTypeDecoder returns a new Decoder that reads from the given
// reader.  Types are decoded separately through the typeDec.
func NewDecoderWithTypeDecoder(r io.Reader, typeDec *TypeDecoder) *Decoder {
	return &Decoder{decoder81{
		buf:     newDecbuf(r),
		typeDec: typeDec,
		flag:    decFlagSeparateTypeDec,
	}}
}

// Decoder returns d as a vdl.Decoder.
func (d *Decoder) Decoder() vdl.Decoder {
	return &d.dec
}

// Decode reads the next value and stores it in value v.  The type of v need not
// exactly match the type of the originally encoded value; decoding succeeds as
// long as the values are convertible.
func (d *Decoder) Decode(v interface{}) error {
	return vdl.Read(&d.dec, v)
}

func (d *decoder81) reset(buf *decbuf) {
	d.buf = buf
	d.flag = 0
	d.stack = d.stackStorage[:0]
	d.refTypes = referencedTypes{}
	d.refAnyLens = referencedAnyLens{}
	d.typeDec = nil
}

func (d *decoder81) appendStack(entry decStackEntry) {
	// TODO(cnicolaou): get rid of this test by ensuring that the
	//      stack is properly initialized when a decoder81 is
	//      created.
	if d.stack == nil {
		d.stack = d.stackStorage[:0]
	}
	d.stack = append(d.stack, entry)
}

func (d *decoder81) IgnoreNextStartValue() {
	d.flag = d.flag.Set(decFlagIgnoreNextStartValue)
}

func (d *decoder81) decodeWireType(wt *wireType) (TypeId, error) {
	// Type messages are just a regularly encoded wireType, which is a union.  To
	// decode we pre-populate the stack with an entry for the wire type, and run
	// the code-generated __VDLRead_wireType method.
	tid, err := d.nextMessage()
	if err != nil {
		return 0, err
	}
	d.appendStack(decStackEntry{
		Type:    wireTypeType,
		Index:   -1,
		LenHint: 1, // wireType is a union
	})
	d.flag = d.flag.Set(decFlagIgnoreNextStartValue)
	if err := vdlReadwireType(d, wt); err != nil {
		return 0, err
	}
	return tid, nil
}

// readRawBytes fills in raw with the next value.  It can be called for both
// top-level and internal values.
func (d *decoder81) readRawBytes(raw *RawBytes) error {
	if d.flag.IgnoreNextStartValue() {
		// If the user has already called StartValue on the decoder, it's harder to
		// capture all the raw bytes, since the optional flag and length hints have
		// already been decoded.  So we simply disallow this from happening.
		return errReadRawBytesAlreadyStarted
	}
	tt, err := d.dfsNextType()
	if err != nil {
		return err
	}
	// Handle top-level values.  All types of values are supported, since we can
	// simply copy the message bytes.
	if len(d.stack) == 0 {
		anyLen, err := d.peekValueByteLen(tt)
		if err != nil {
			return err
		}
		if err := d.decodeRaw(tt, anyLen, raw); err != nil {
			return err
		}
		return d.endMessage()
	}
	// Handle internal values.  Only any values are supported at the moment, since
	// they come with a header that tells us the exact length to read.
	//
	// TODO(toddw): Handle other types, either by reading and skipping bytes based
	// on the type, or by falling back to a decode / re-encode slowpath.
	if tt.Kind() != vdl.Any {
		return errReadRawBytesFromNonAny
	}
	ttElem, anyLen, err := d.readAnyHeader()
	if err != nil {
		return err
	}
	if ttElem == nil {
		// This is a nil any value, which has already been read by readAnyHeader.
		// We simply fill in RawBytes with the single WireCtrlNil byte.
		raw.Version = d.buf.version
		raw.Type = vdl.AnyType
		raw.RefTypes = nil
		raw.AnyLengths = nil
		raw.Data = []byte{WireCtrlNil}
		return nil
	}
	return d.decodeRaw(ttElem, anyLen, raw)
}

func (d *decoder81) StartValue(want *vdl.Type) error {
	if d.flag.IgnoreNextStartValue() {
		d.flag = d.flag.Clear(decFlagIgnoreNextStartValue)
		return nil
	}
	tt, err := d.dfsNextType()
	if err != nil {
		return err
	}
	tt, lenHint, flag, err := d.setupType(tt, want)
	if err != nil {
		return err
	}
	d.appendStack(decStackEntry{
		Type:    tt,
		Index:   -1,
		LenHint: lenHint,
		Flag:    flag,
	})
	return nil
}

func (d *decoder81) setupType(tt, want *vdl.Type) (_ *vdl.Type, lenHint int, flag decStackFlag, _ error) { //nolint:gocyclo
	// Handle any, which may be nil.  We "dereference" non-nil any to the inner
	// type.  If that happens to be an optional, it's handled below.
	if tt.Kind() == vdl.Any {
		flag |= decStackFlagIsAny
		var err error
		switch tt, _, err = d.readAnyHeader(); {
		case err != nil:
			return nil, 0, 0, err
		case tt == nil:
			tt = vdl.AnyType // nil any
		}
	}
	// Handle optional, which may be nil.  Similar to any, we "dereference"
	// non-nil optional to the inner type, which is never allowed to be another
	// optional or any type.
	if tt.Kind() == vdl.Optional {
		flag |= decStackFlagIsOptional
		// Read the WireCtrlNil code, but if it's not WireCtrlNil we need to keep
		// the buffer as-is, since it's the first byte of the value, which may
		// itself be another control code.
		switch ctrl, err := binaryPeekControl(d.buf); {
		case err != nil:
			return nil, 0, 0, err
		case ctrl == WireCtrlNil:
			d.buf.SkipAvailable(1) // nil optional
		default:
			tt = tt.Elem() // non-nil optional
		}
	}
	// Check compatibility between the actual type and the want type.  Since
	// compatibility applies to the entire static type, we only need to perform
	// this check for top-level decoded values, and subsequently for decoded any
	// values.  We skip checking non-composite want types, since those will be
	// naturally caught by the Decode* calls anyways.
	if want != nil && (len(d.stack) == 0 || flag.IsAny()) {
		switch want.Kind() {
		case vdl.Optional, vdl.Array, vdl.List, vdl.Set, vdl.Map, vdl.Struct, vdl.Union:
			if tt == want {
				// Set FastRead flag, which will let us use the fastpath for ReadValue*
				// in common cases.  We can only use this fastpath if tt and want are
				// identical, which ensures we don't need to perform any conversions.
				if !flag.IsAny() && isFastReadParent(tt) {
					flag |= decStackFlagFastRead
				}
				// Regardless of whether we can use the fastpath, there's no need to
				// check compatibility if tt and want are identical.
			} else {
				if !vdl.Compatible(tt, want) {
					return nil, 0, 0, errIncompatibleDecode(tt, want)
				}
			}
		}
	}
	// Initialize LenHint for composite types.
	switch tt.Kind() {
	case vdl.Array, vdl.List, vdl.Set, vdl.Map:
		// TODO(toddw): Handle sentry-terminated collections without a length hint.
		len, err := binaryDecodeLenOrArrayLen(d.buf, tt)
		if err != nil {
			return nil, 0, 0, err
		}
		lenHint = len
	case vdl.Union:
		// Union shouldn't have a LenHint, but we abuse it in NextField as a
		// convenience for detecting when fields are done, so we initialize it here.
		// It has to be at least 1, since 0 will cause NextField to think that the
		// union field has already been decoded.
		lenHint = 1
	case vdl.Struct:
		// Struct shouldn't have a LenHint, but we abuse it in NextField as a
		// convenience for detecting when fields are done, so we initialize it here.
		lenHint = tt.NumField()
	default:
		lenHint = -1
	}
	if top := d.top(); top != nil && top.Type.IsBytes() {
		d.flag = d.flag.Set(decFlagIsParentBytes)
	} else {
		d.flag = d.flag.Clear(decFlagIsParentBytes)
	}
	return tt, lenHint, flag, nil
}

func errIncompatibleDecode(tt *vdl.Type, want interface{}) error {
	return fmt.Errorf("vom: incompatible decode from %v into %v", tt, want)
}

func (d *decoder81) FinishValue() error {
	d.flag = d.flag.Clear(decFlagFinishValue)
	stackTop := len(d.stack) - 1
	if stackTop == -1 {
		return errEmptyDecoderStack
	}
	d.stack = d.stack[:stackTop]
	if stackTop == 0 {
		return d.endMessage()
	}
	return nil
}

func (d *decoder81) top() *decStackEntry {
	if stackTop := len(d.stack) - 1; stackTop >= 0 {
		return &d.stack[stackTop]
	}
	return nil
}

// dfsNextType determines the type of the next value that we will decode, by
// walking the static type in DFS order.  To bootstrap we retrieve the top-level
// type from the VOM value message.
func (d *decoder81) dfsNextType() (*vdl.Type, error) {
	top := d.top()
	if top == nil {
		// Bootstrap: start decoding a new top-level value.
		if !d.flag.SeparateTypeDec() {
			if err := d.decodeTypeDefs(); err != nil {
				return nil, err
			}
		}
		tid, err := d.nextMessage()
		if err != nil {
			return nil, err
		}
		return d.typeDec.lookupType(tid)
	}
	// Return the next type from our composite types.
	tt := top.Type
	switch tt.Kind() {
	case vdl.Array, vdl.List:
		return tt.Elem(), nil
	case vdl.Set:
		return tt.Key(), nil
	case vdl.Map:
		top.Flag = top.Flag.FlipIsMapKey()
		if top.Flag.IsMapKey() {
			return tt.Key(), nil
		}
		return tt.Elem(), nil
	case vdl.Union, vdl.Struct:
		return tt.Field(top.Index).Type, nil
	}
	return nil, fmt.Errorf("vom: can't StartValue on %v", tt)
}

func (d *decoder81) NextEntry() (bool, error) {
	// Our strategy is to increment top.Index until it reaches top.LenHint.
	// Currently the LenHint is always set, so it's stronger than a hint.
	//
	// TODO(toddw): Handle sentry-terminated collections without a LenHint.
	top := d.top()
	if top == nil {
		return false, errEmptyDecoderStack
	}
	// Increment index and check errors.
	top.Index++
	switch top.Type.Kind() {
	case vdl.Array, vdl.List, vdl.Set, vdl.Map:
		if top.Index > top.LenHint && top.LenHint >= 0 {
			return false, fmt.Errorf("vom: NextEntry called after done, stack: %+v", d.stack)
		}
	default:
		return false, fmt.Errorf("vom: NextEntry called on invalid type, stack: %+v", d.stack)
	}
	return top.Index == top.LenHint, nil
}

func (d *decoder81) NextField() (int, error) { //nolint:gocyclo
	top := d.top()
	if top == nil {
		return -1, errEmptyDecoderStack
	}
	// Increment index and check errors.  Note that the actual top.Index is
	// decoded from the buf data stream; we use top.LenHint to help detect when
	// the fields are done, and to detect invalid calls after we're done.
	top.Index++
	switch top.Type.Kind() {
	case vdl.Union, vdl.Struct:
		if top.Index > top.LenHint {
			return -1, fmt.Errorf("vom: NextField called after done, stack: %+v", d.stack)
		}
	default:
		return -1, fmt.Errorf("vom: NextField called on invalid type, stack: %+v", d.stack)
	}
	var field int
	switch top.Type.Kind() {
	case vdl.Union:
		if top.Index == top.LenHint {
			// We know we're done since we set LenHint=Index+1 the first time around,
			// and we incremented the index above.
			return -1, nil
		}
		// Decode the union field index.
		switch index, err := binaryDecodeUint(d.buf); {
		case err != nil:
			return -1, err
		case index >= uint64(top.Type.NumField()):
			return -1, errIndexOutOfRange
		default:
			// Set LenHint=Index+1 so that we'll know we're done next time around.
			field = int(index)
			top.Index = field
			top.LenHint = field + 1
		}
	case vdl.Struct:
		// Handle the end-of-struct sentry.
		switch ok, err := binaryDecodeControlOnly(d.buf, WireCtrlEnd); {
		case err != nil:
			return -1, err
		case ok:
			// Set Index=LenHint to ensure repeated calls will fail.
			top.Index = top.LenHint
			return -1, nil
		}
		// Decode the struct field index.
		switch index, err := binaryDecodeUint(d.buf); {
		case err != nil:
			return -1, err
		case index >= uint64(top.Type.NumField()):
			return -1, errIndexOutOfRange
		default:
			field = int(index)
			top.Index = field
		}
	}
	return field, nil
}

func (d *decoder81) Type() *vdl.Type {
	if top := d.top(); top != nil {
		return top.Type
	}
	return nil
}

func (d *decoder81) IsAny() bool {
	if top := d.top(); top != nil {
		return top.Flag.IsAny()
	}
	return false
}

func (d *decoder81) IsOptional() bool {
	if top := d.top(); top != nil {
		return top.Flag.IsOptional()
	}
	return false
}

func (d *decoder81) IsNil() bool {
	if top := d.top(); top != nil {
		// Because of the "dereferencing" we do, the only time the type is any or
		// optional is when it's nil.
		return top.Type == vdl.AnyType || top.Type.Kind() == vdl.Optional
	}
	return false
}

func (d *decoder81) Index() int {
	if top := d.top(); top != nil {
		return top.Index
	}
	return -1
}

func (d *decoder81) LenHint() int {
	if top := d.top(); top != nil {
		// Note that union and struct shouldn't have a LenHint, but we abuse it in
		// NextField as a convenience for detecting when fields are done, so an
		// "arbitrary" value is returned here.  Users shouldn't be looking at it for
		// union and struct anyways.
		return top.LenHint
	}
	return -1
}

func (d *decoder81) DecodeBool() (bool, error) {
	tt := d.Type()
	if tt == nil {
		return false, errEmptyDecoderStack
	}
	if tt.Kind() == vdl.Bool {
		return binaryDecodeBool(d.buf)
	}
	return false, errIncompatibleDecode(tt, "bool")
}

func (d *decoder81) DecodeString() (string, error) {
	tt := d.Type()
	if tt == nil {
		return "", errEmptyDecoderStack
	}
	switch tt.Kind() {
	case vdl.String:
		return binaryDecodeString(d.buf)
	case vdl.Enum:
		return d.binaryDecodeEnum(tt)
	}
	return "", errIncompatibleDecode(tt, "string")
}

func (d *decoder81) binaryDecodeEnum(tt *vdl.Type) (string, error) {
	index, err := binaryDecodeUint(d.buf)
	switch {
	case err != nil:
		return "", err
	case index >= uint64(tt.NumEnumLabel()):
		return "", fmt.Errorf("vom: enum index %d out of range, %v", index, tt)
	}
	return tt.EnumLabel(int(index)), nil
}

func (d *decoder81) binaryDecodeByte() (byte, error) {
	// Handle a special-case where normally single bytes are written out as
	// variable sized numbers, which use 2 bytes to encode bytes > 127.  But each
	// byte contained in a list or array is written out as one byte.  E.g.
	//   byte(0x81)         -> 0xFF81   : single byte with variable-size
	//   []byte("\x81\x82") -> 0x028182 : each elem byte encoded as one byte
	if d.flag.IsParentBytes() {
		return d.buf.ReadByte()
	}
	x, err := binaryDecodeUint(d.buf)
	return byte(x), err
}

func (d *decoder81) DecodeUint(bitlen int) (uint64, error) {
	tt := d.Type()
	if tt == nil {
		return 0, errEmptyDecoderStack
	}
	return d.decodeUint(tt, uint(bitlen))
}

func (d *decoder81) decodeUint(tt *vdl.Type, ubitlen uint) (uint64, error) {
	const errFmt = "vom: conversion from %v into uint%d loses precision: %v"
	switch tt.Kind() {
	case vdl.Byte:
		x, err := d.binaryDecodeByte()
		if err != nil {
			return 0, err
		}
		return uint64(x), err
	case vdl.Uint16, vdl.Uint32, vdl.Uint64:
		x, err := binaryDecodeUint(d.buf)
		if err != nil {
			return 0, err
		}
		if shift := 64 - ubitlen; x != (x<<shift)>>shift {
			return 0, fmt.Errorf(errFmt, tt, ubitlen, x)
		}
		return x, nil
	case vdl.Int8, vdl.Int16, vdl.Int32, vdl.Int64:
		x, err := binaryDecodeInt(d.buf)
		if err != nil {
			return 0, err
		}
		ux := uint64(x)
		if shift := 64 - ubitlen; x < 0 || ux != (ux<<shift)>>shift {
			return 0, fmt.Errorf(errFmt, tt, ubitlen, x)
		}
		return ux, nil
	case vdl.Float32, vdl.Float64:
		x, err := binaryDecodeFloat(d.buf)
		if err != nil {
			return 0, err
		}
		ux := uint64(x)
		if shift := 64 - ubitlen; x != float64(ux) || ux != (ux<<shift)>>shift {
			return 0, fmt.Errorf(errFmt, tt, ubitlen, x)
		}
		return ux, nil
	}
	return 0, errIncompatibleDecode(tt, fmt.Sprintf("uint%d", ubitlen))
}

func (d *decoder81) DecodeInt(bitlen int) (int64, error) {
	tt := d.Type()
	if tt == nil {
		return 0, errEmptyDecoderStack
	}
	return d.decodeInt(tt, uint(bitlen))
}

func (d *decoder81) decodeInt(tt *vdl.Type, ubitlen uint) (int64, error) { //nolint:gocyclo
	const errFmt = "vom: conversion from %v into int%d loses precision: %v"
	switch tt.Kind() {
	case vdl.Byte:
		x, err := d.binaryDecodeByte()
		if err != nil {
			return 0, err
		}
		// The only case that fails is if we're converting byte(x) to int8, and x
		// uses more than 7 bits (i.e. is greater than 127).
		if ubitlen <= 8 && x > 0x7f {
			return 0, fmt.Errorf(errFmt, tt, ubitlen, x)
		}
		return int64(x), nil
	case vdl.Uint16, vdl.Uint32, vdl.Uint64:
		x, err := binaryDecodeUint(d.buf)
		if err != nil {
			return 0, err
		}
		ix := int64(x)
		// The shift uses 65 since the topmost bit is the sign bit.  I.e. 32 bit
		// numbers should be shifted by 33 rather than 32.
		if shift := 65 - ubitlen; ix < 0 || x != (x<<shift)>>shift {
			return 0, fmt.Errorf(errFmt, tt, ubitlen, x)
		}
		return ix, nil
	case vdl.Int8, vdl.Int16, vdl.Int32, vdl.Int64:
		x, err := binaryDecodeInt(d.buf)
		if err != nil {
			return 0, err
		}
		if shift := 64 - ubitlen; x != (x<<shift)>>shift {
			return 0, fmt.Errorf(errFmt, tt, ubitlen, x)
		}
		return x, nil
	case vdl.Float32, vdl.Float64:
		x, err := binaryDecodeFloat(d.buf)
		if err != nil {
			return 0, err
		}
		ix := int64(x)
		if shift := 64 - ubitlen; x != float64(ix) || ix != (ix<<shift)>>shift {
			return 0, fmt.Errorf(errFmt, tt, ubitlen, x)
		}
		return ix, nil
	}
	return 0, errIncompatibleDecode(tt, fmt.Sprintf("int%d", ubitlen))
}

func (d *decoder81) DecodeFloat(bitlen int) (float64, error) {
	tt := d.Type()
	if tt == nil {
		return 0, errEmptyDecoderStack
	}
	return d.decodeFloat(tt, uint(bitlen))
}

func (d *decoder81) decodeFloat(tt *vdl.Type, ubitlen uint) (float64, error) { //nolint:gocyclo
	const errFmt = "vom: conversion from %v into float%d loses precision: %v"
	switch tt.Kind() {
	case vdl.Byte:
		x, err := d.binaryDecodeByte()
		if err != nil {
			return 0, err
		}
		return float64(x), nil
	case vdl.Uint16, vdl.Uint32, vdl.Uint64:
		x, err := binaryDecodeUint(d.buf)
		if err != nil {
			return 0, err
		}
		var max uint64
		if ubitlen > 32 {
			max = float64MaxInt
		} else {
			max = float32MaxInt
		}
		if x > max {
			return 0, fmt.Errorf(errFmt, tt, ubitlen, x)
		}
		return float64(x), nil
	case vdl.Int8, vdl.Int16, vdl.Int32, vdl.Int64:
		x, err := binaryDecodeInt(d.buf)
		if err != nil {
			return 0, err
		}
		var min, max int64
		if ubitlen > 32 {
			min, max = float64MinInt, float64MaxInt
		} else {
			min, max = float32MinInt, float32MaxInt
		}
		if x < min || x > max {
			return 0, fmt.Errorf(errFmt, tt, ubitlen, x)
		}
		return float64(x), nil
	case vdl.Float32, vdl.Float64:
		x, err := binaryDecodeFloat(d.buf)
		if err != nil {
			return 0, err
		}
		if ubitlen <= 32 && (x < -math.MaxFloat32 || x > math.MaxFloat32) {
			return 0, fmt.Errorf(errFmt, tt, ubitlen, x)
		}
		return x, nil
	}
	return 0, errIncompatibleDecode(tt, fmt.Sprintf("float%d", ubitlen))
}

func (d *decoder81) DecodeBytes(fixedLen int, v *[]byte) error {
	top := d.top()
	if top == nil {
		return errEmptyDecoderStack
	}
	tt := top.Type
	if !tt.IsBytes() {
		return vdl.DecodeConvertedBytes(d, fixedLen, v)
	}
	return d.decodeBytes(tt, top.LenHint, fixedLen, v)
}

func (d *decoder81) decodeBytes(tt *vdl.Type, lenHint, fixedLen int, v *[]byte) error {
	switch {
	case lenHint == -1:
		return fmt.Errorf("vom: LenHint is currently required, %v", tt)
	case fixedLen >= 0 && fixedLen != lenHint:
		return fmt.Errorf("vom: got %d bytes, want fixed len %d, %v", lenHint, fixedLen, tt)
	case lenHint == 0:
		*v = nil
		return nil
	case fixedLen >= 0:
		// Only re-use the existing buffer if we're filling in an array.  This
		// sacrifices some performance, but also avoids bugs when repeatedly
		// decoding into the same value.
		*v = (*v)[:lenHint]
	default:
		*v = make([]byte, lenHint)
	}
	return d.buf.ReadIntoBuf(*v)
}

func (d *decoder81) DecodeTypeObject() (*vdl.Type, error) {
	tt := d.Type()
	if tt == nil {
		return nil, errEmptyDecoderStack
	}
	if tt != vdl.TypeObjectType {
		return nil, errIncompatibleDecode(tt, "typeobject")
	}
	return d.binaryDecodeType()
}

func (d *decoder81) binaryDecodeType() (*vdl.Type, error) {
	typeIndex, err := binaryDecodeUint(d.buf)
	if err != nil {
		return nil, err
	}
	tid, err := d.refTypes.ReferencedTypeID(typeIndex)
	if err != nil {
		return nil, err
	}
	return d.typeDec.lookupType(tid)
}

func (d *decoder81) SkipValue() error {
	tt, err := d.dfsNextType()
	if err != nil {
		return err
	}
	if len(d.stack) == 0 {
		// Handle top-level values.  It's easy to determine the byte length of the
		// value, so we can just skip the bytes.
		valueLen, err := d.peekValueByteLen(tt)
		if err != nil {
			return err
		}
		if err := d.buf.Skip(valueLen); err != nil {
			return err
		}
		return d.endMessage()
	}
	return d.skipValue(tt)
}
