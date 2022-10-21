// Package binson implements s small, high-performance implementation of Binson, see binson.org.
package binson

import (
	"unsafe"
)

// ValueType is type signature for each binson item
type ValueType uint

// Binson item types enumeration
const (
	Boolean ValueType = iota
	Integer
	Double
	String
	Bytes
	Array
	Object
)

// Binson item signatures
const (
	sigBegin      byte = 0x40
	sigEnd        byte = 0x41
	sigBeginArray byte = 0x42
	sigEndArray   byte = 0x43
	sigTrue       byte = 0x44
	sigFalse      byte = 0x45
	sigInteger1   byte = 0x10
	sigInteger2   byte = 0x11
	sigInteger4   byte = 0x12
	sigInteger8   byte = 0x13
	sigDouble     byte = 0x46
	sigString1    byte = 0x14
	sigString2    byte = 0x15
	sigString4    byte = 0x16
	sigBytes1     byte = 0x18
	sigBytes2     byte = 0x19
	sigBytes4     byte = 0x1a
)

const intLengthMask byte = 0x03
const oneByte byte = 0x00
const twoBytes byte = 0x01
const fourBytes byte = 0x02
const eightBytes byte = 0x03
const twoTo7 int64 = 128
const twoTo15 int64 = 32768
const twoTo31 int64 = 2147483648

// Binson Decoder private constants
const (
	stateZero = iota
	stateBeforeField
	stateBeforeArrayValue
	stateBeforeArray
	stateEndOfArray
	stateBeforeObject
	stateEndOfObject
)

// Error codes for err field of Encoder and Decoder
const ErrorNone = 0                // All ok, no error
const ErrorEOF = 1                 // Reached EOF (input bffuer ends before Binson object ends)
const ErrorEndOfObject = 2         // At end of Binson object, cannot parse further
const ErrorNotReadyToReadField = 3 // Not ready to read a field
const ErrorUnexpectedTypeByte = 4
const ErrorNotBeforeArrayValue = 5
const ErrorNotBeforeObject = 6
const ErrorNotBeforeArray = 7
const ErrorCannotGoUpToObject = 8
const ErrorCannotGoUpToArray = 9
const ErrorUnexpectedType = 10
const ErrorNegativeLength = 11
const ErrorLengthTooLarge = 12
const ErrorExpectedBegin = 13

// ======== Decoder ========

// A Decoder represents an Binson parser reading a particular input stream.
type Decoder struct {
	buf     []byte
	offset  int
	state   int
	sigByte byte
	err     int

	Name      []byte
	Value     interface{}
	ValueType ValueType
}

func NewDecoderFromBytes(buf []byte) *Decoder {
	return &Decoder{
		state:  stateZero,
		buf:    buf,
		offset: 0,
	}
}

// Field parses until an expected field with the given name is found
// in the current Binson object.
func (d *Decoder) Field(name string) bool {
	for d.NextField() {
		if d.err != ErrorNone {
			return false
		}

		// This does no heap alloc.
		if name == string(d.Name) {
			return true
		}
	}

	return false
}

// NextField reads next field, returns true if a field was found and false
// if end-of-object was reached.
// If  boolean/integer/double/bytes/string was found, the value is also read
// and is available in `Value` field
func (d *Decoder) NextField() bool {
	switch d.state {
	case stateZero:
		d.parseBegin()
	case stateEndOfObject:
		d.err = ErrorEndOfObject
		return false
	case stateBeforeObject:
		d.state = stateBeforeField
		for d.NextField() {
			if d.err != ErrorNone {
				return false
			}
		}
		d.state = stateBeforeField
	case stateBeforeArray:
		d.state = stateBeforeArrayValue
		for d.NextArrayValue() {
			if d.err != ErrorNone {
				return false
			}
		}
		d.state = stateBeforeField
	}

	if d.state != stateBeforeField {
		d.err = ErrorNotReadyToReadField
		return false
	}

	typeBeforeName := d.readOne()
	if d.err != ErrorNone {
		return false
	}
	if typeBeforeName == sigEnd {
		d.state = stateEndOfObject
		return false
	}
	d.parseFieldName(typeBeforeName)

	typeBeforeValue := d.readOne()
	if d.err != ErrorNone {
		return false
	}
	d.parseValue(typeBeforeValue, stateBeforeField)

	return true
}

// NextArrayValue reads next binson ARRAY value,
// returns true if a field was found and false, if end-of-object was reached.
// If boolean/integer/double/bytes/string was found, the value is also read
// and is available in `Value` field
func (d *Decoder) NextArrayValue() bool {
	if d.state == stateBeforeArray {
		d.state = stateBeforeArrayValue
		for d.NextArrayValue() {
			if d.err != ErrorNone {
				return false
			}
		}
		d.state = stateBeforeArrayValue
	}

	if d.state == stateBeforeObject {
		d.state = stateBeforeField
		for d.NextField() {
			if d.err != ErrorNone {
				return false
			}
		}
		d.state = stateBeforeArrayValue
	}

	if d.state != stateBeforeArrayValue {
		//d.err = parseError("not before array value")
		d.err = ErrorNotBeforeArrayValue
		return false
	}

	sig := d.readOne()
	if d.err != ErrorNone {
		return false
	}
	if sig == sigEndArray {
		d.state = stateEndOfArray
		return false
	}
	d.parseValue(sig, stateBeforeArrayValue)

	return true
}

// GoIntoObject navigates decoder inside the expected OBJECT
func (d *Decoder) GoIntoObject() {
	if d.state != stateBeforeObject {
		//d.err = parseError("not an object field") XXX
		d.err = ErrorNotBeforeObject
		return
	}
	d.state = stateBeforeField
}

// GoIntoArray navigates decoder inside the expected ARRAY
func (d *Decoder) GoIntoArray() {
	if d.state != stateBeforeArray {
		//d.err = parseError("unexpected parser state, not an array field") // XXX
		d.err = ErrorNotBeforeArray
		return
	}
	d.state = stateBeforeArrayValue
}

// GoUpToObject navigates decoder to the parent OBJECT
func (d *Decoder) GoUpToObject() {
	if d.state == stateBeforeArrayValue {
		for d.NextArrayValue() {
			if d.err != ErrorNone {
				return
			}
		}
	}

	if d.state == stateBeforeField {
		for d.NextField() {
			if d.err != ErrorNone {
				return
			}
		}
	}

	if d.state != stateEndOfObject && d.state != stateEndOfArray {
		//d.err = parseError("unexpected parser state")
		d.err = ErrorCannotGoUpToObject
		return
	}

	d.state = stateBeforeField
}

// GoUpToArray navigates decoder to the parent ARRAY
func (d *Decoder) GoUpToArray() {
	if d.state == stateBeforeArrayValue {
		for d.NextArrayValue() {
			if d.err != ErrorNone {
				return
			}
		}
	}

	if d.state == stateBeforeField {
		for d.NextField() {
			if d.err != ErrorNone {
				return
			}
		}
	}

	if d.state != stateEndOfObject && d.state != stateEndOfArray {
		//d.err = parseError("unexpected parser state")
		d.err = ErrorCannotGoUpToArray
		return
	}

	d.state = stateBeforeArrayValue
}

// Private methods

func (d *Decoder) parseValue(sigByte byte, afterValueState int) {
	switch sigByte {
	case sigBegin:
		d.ValueType = Object
		d.state = stateBeforeObject
	case sigBeginArray:
		d.ValueType = Array
		d.state = stateBeforeArray
	case sigFalse, sigTrue:
		d.ValueType = Boolean
		d.Value = sigByte == sigTrue
		d.state = afterValueState
	case sigDouble:
		var d64 float64
		var i64 int64
		d.ValueType = Double
		d.readInt64(&i64)
		d64 = float64frombits(uint64(i64))
		d.Value = d64
		d.state = afterValueState
	case sigInteger1, sigInteger2, sigInteger4, sigInteger8:
		d.ValueType = Integer
		d.Value = d.parseInteger(sigByte)
		d.state = afterValueState
	case sigString1, sigString2, sigString4:
		d.ValueType = String
		d.Value = d.parseStringBytes(sigByte)
		d.state = afterValueState
	case sigBytes1, sigBytes2, sigBytes4:
		d.ValueType = Bytes
		d.Value = d.parseStringBytes(sigByte)
		d.state = afterValueState
	default:
		d.err = ErrorUnexpectedTypeByte
	}
}

func (d *Decoder) parseFieldName(sigBeforeName byte) {
	switch sigBeforeName {
	case sigString1, sigString2, sigString4:
		d.Name = d.parseStringBytes(sigBeforeName)
	default:
		d.err = ErrorUnexpectedType
	}
}

func (d *Decoder) parseBegin() {
	d.sigByte = d.readOne()

	if d.sigByte != sigBegin {
		d.err = ErrorExpectedBegin
		return
	}
	d.state = stateBeforeField
}

// Parses name or string or bytes.
func (d *Decoder) parseStringBytes(sigByte byte) []byte {
	ln := d.parseInteger(sigByte)

	if ln < 0 {
		//d.err = parseError("negative name/string/bytes length")
		d.err = ErrorNegativeLength
		return nil
	}

	if ln >= 2^31 {
		//d.err = parseError("name/string/bytes length too large")
		d.err = ErrorLengthTooLarge
		return nil
	}

	// TODO: check if length exceeds total length of buffer
	// Note, should there be an expected max length? That may be shorter
	// then the total size of the buffer.

	intLn := int(ln)

	if d.offset+intLn > len(d.buf) {
		//d.err = parseError("EOF")
		d.err = ErrorEOF
		return nil
	}

	result := d.buf[d.offset : d.offset+intLn]
	d.offset += intLn

	//var buf = make([]byte, ln) // TODO remove dynalloc
	//err := d.readToBuffer(buf)
	//if err != nil {
	//		return nil
	//}

	return result

	//if sigByte >= sigBytes1 {
	//	return buf
	//}
	//return string(buf)
}

func (d *Decoder) parseInteger(sigByte byte) int64 {
	switch sigByte & intLengthMask {
	case oneByte:
		var i1 int8
		d.readInt8(&i1)
		return int64(i1)
	case twoBytes:
		var i2 int16
		d.readInt16(&i2)
		return int64(i2)
	case fourBytes:
		var i4 int32
		d.readInt32(&i4)
		return int64(i4)
	case eightBytes:
		var i8 int64
		d.readInt64(&i8)
		return int64(i8)
	default:
		panic("never happens")
	}
}

// Reads one byte from the buffer.
func (d *Decoder) readOne() byte {
	if d.offset >= len(d.buf) {
		d.err = ErrorEOF
		return 0
	}
	b := d.buf[d.offset]
	d.offset++
	return b
}

func (d *Decoder) readInt8(a *int8) {
	if d.offset+1 > len(d.buf) {
		*a = 0
		d.err = ErrorEOF
		return
	}
	*a = int8(d.buf[d.offset])
	d.offset++
}

func (d *Decoder) readInt16(a *int16) {
	if d.offset+2 > len(d.buf) {
		*a = 0
		d.err = ErrorEOF
		return
	}

	myUint16 := getUint16(d.buf[d.offset:])
	*a = int16(myUint16)
	d.offset += 2
}

func (d *Decoder) readInt32(a *int32) {
	if d.offset+4 > len(d.buf) {
		*a = 0
		d.err = ErrorEOF
	}

	myUint32 := getUint32(d.buf[d.offset:])
	*a = int32(myUint32)
	d.offset += 4
}

func (d *Decoder) readInt64(a *int64) {
	if d.offset+8 > len(d.buf) {
		d.err = ErrorEOF
		*a = 0
		return
	}

	myUint64 := getUint64(d.buf[d.offset:])
	*a = int64(myUint64)
	d.offset += 8
}

func (d *Decoder) readToBuffer(toBuffer []byte) {
	ln := len(toBuffer)
	if d.offset+ln > len(d.buf) {
		//return parseError("EOF")
		d.err = ErrorEOF
		return
	}

	for i := 0; i < ln; i++ {
		toBuffer[i] = d.buf[d.offset+i]
	}

	d.offset += ln
}

// ========= Encoder ========

// An Encoder writes binson data to an output stream.
type Encoder struct {
	err    error
	buf    []byte // buffer to write output to
	offset int    // next position in buf to write to
}

func NewEncoderFromBytes(buf []byte) *Encoder {
	return &Encoder{buf: buf}
}

// Flush encoder buffers
// Does nothing in this implementation.
// Keep it if we want same API, binson-go and binson-go-tiny.
func (e *Encoder) Flush() {
}

// Begin writes OBJECT begin signature to output stream
func (e *Encoder) Begin() {
	e.writeOne(sigBegin)
}

// End writes OBJECT end signature to output stream
func (e *Encoder) End() {
	e.writeOne(sigEnd)
}

// BeginArray writes ARRAY begin signature to output stream
func (e *Encoder) BeginArray() {
	e.writeOne(sigBeginArray)
}

// EndArray writes ARRAY end signature to output stream
func (e *Encoder) EndArray() {
	e.writeOne(sigEndArray)
}

// Bool writes specified boolean value to output stream
func (e *Encoder) Bool(val bool) {
	var sig = sigTrue
	if !val {
		sig = sigFalse
	}
	e.writeOne(sig)
}

// Integer writes specified integer value to output stream
func (e *Encoder) Integer(val int64) {
	e.writeIntegerOrLength(sigInteger1, val)
}

// Double writes float64 value to output stream
func (e *Encoder) Double(val float64) {
	e.writeOne(sigDouble)
	var myUint uint64 = float64bits(val)
	e.writeInt64(int64(myUint))
}

// String writes string value to output stream
func (e *Encoder) String(val string) {
	e.writeIntegerOrLength(sigString1, int64(len(val)))
	e.write([]byte(val)) // TODO check this!
}

// Bytes writes []byte value to output stream
func (e *Encoder) Bytes(val []byte) {
	e.writeIntegerOrLength(sigBytes1, int64(len(val)))
	e.write(val)
}

// Name writes string value as OBJECT item's name to output stream
func (e *Encoder) Name(val string) {
	e.String(val)
}

/* === private methods === */

func (e *Encoder) writeIntegerOrLength(baseType byte, val int64) {
	switch {
	case val >= -twoTo7 && val < twoTo7:
		e.writeOne(baseType | oneByte)
		e.writeInt8(int8(val))
	case val >= -twoTo15 && val < twoTo15:
		e.writeOne(baseType | twoBytes)
		e.writeInt16(int16(val))
	case val >= -twoTo31 && val < twoTo31:
		e.writeOne(baseType | fourBytes)
		e.writeInt32(int32(val))
	default:
		e.writeOne(baseType | eightBytes)
		e.writeInt64(int64(val))
	}
}

// Returns true if s bytes can be written to output.
// If not, e.err is set to EOF and false is returned.
func (e *Encoder) available(s int) bool {
	// TODO check, should we use len() or capacity?
	if e.offset+s > len(e.buf) {
		e.err = writeError("EOF")
		return false
	}
	return true
}

func (e *Encoder) writeOne(b byte) {
	if !e.available(1) {
		return
	}

	e.buf[e.offset] = b
	e.offset += 1
}

func (e *Encoder) write(b []byte) {
	lenb := len(b)
	if !e.available(lenb) {
		return
	}
	for i := 0; i < lenb; i++ {
		e.buf[e.offset+i] = b[i]
	}
	e.offset += lenb
}

func (e *Encoder) writeInt8(i int8) {
	e.writeOne(byte(i))
}

func (e *Encoder) writeInt16(i int16) {
	if !e.available(2) {
		return
	}
	putUint16(e.buf[e.offset:], uint16(i))
	e.offset += 2
}

func (e *Encoder) writeInt32(i int32) {
	if !e.available(4) {
		return
	}
	putUint32(e.buf[e.offset:], uint32(i))
	e.offset += 4
}

func (e *Encoder) writeInt64(i int64) {
	if !e.available(8) {
		return
	}
	putUint64(e.buf[e.offset:], uint64(i))
	e.offset += 8
}

type WriteError struct {
	text string
}

func (e WriteError) Error() string {
	return e.text
}

func writeError(text string) WriteError {
	return WriteError{text}
}

// ======== Instead of math ========
// Code in this section removes dependency on math package.

// Float64bits returns the IEEE 754 binary representation of f.
// Equivalent to math.Float64bits().
func float64bits(f float64) uint64 {
	return *(*uint64)(unsafe.Pointer(&f))
}

// Float64frombits returns the floating-point number corresponding
// to the IEEE 754 binary representation b.
// Equivalent to math.Float64frombits().
func float64frombits(b uint64) float64 {
	return *(*float64)(unsafe.Pointer(&b))
}

// ======== Instead of binary ========
// Code in this section removes dependency on binary package.
// Little-endian encoding is assumed. As used by Binson.
// Early bounds checks (decreasing indexes) see code in binary package and
// golang.org/issue/14808. Can improve performance. Check TinyGo performance.
//
// TODO: check source, copyright etc.
// Given 32-bit machine (TinyGo on Cortex M0, for example), can we optimized
// this? On little-ending machines, should be possible to do fast conversions
// using Unsafe.

func getUint64(b []byte) uint64 {
	_ = b[7] // bounds check hint to compiler; see golang.org/issue/14808
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
}

func putUint64(b []byte, v uint64) {
	_ = b[7] // early bounds check to guarantee safety of writes below
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	b[4] = byte(v >> 32)
	b[5] = byte(v >> 40)
	b[6] = byte(v >> 48)
	b[7] = byte(v >> 56)
}

func getUint32(b []byte) uint32 {
	_ = b[3] // bounds check hint to compiler; see golang.org/issue/14808
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

func putUint32(b []byte, v uint32) {
	_ = b[3] // early bounds check to guarantee safety of writes below
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
}

func getUint16(b []byte) uint16 {
	_ = b[1] // early bounds check; see golang.org/issue/14808
	return uint16(b[0]) | uint16(b[1])<<8
}

func putUint16(b []byte, v uint16) {
	_ = b[1] // early bounds check
	b[0] = byte(v)
	b[1] = byte(v >> 8)
}
