// Test internal (private) serialization / deserialization routines
package binson

import (
	"bytes"
	"encoding/hex"
	"math"
	"testing"
)

// Binson INTEGER internal representation test data table
var intTable = []struct {
	val int64
	raw []byte
}{
	// int8
	{0, []byte("\x10\x00")},
	{-1, []byte("\x10\xff")},
	{math.MaxInt8, []byte("\x10\x7f")},
	{math.MaxInt8 + 1, []byte("\x11\x80\x00")},
	{math.MinInt8, []byte("\x10\x80")},
	{math.MinInt8 - 1, []byte("\x11\x7f\xff")},

	// int16
	{math.MaxInt16, []byte("\x11\xff\x7f")},
	{math.MaxInt16 + 1, []byte("\x12\x00\x80\x00\x00")},
	{math.MinInt16, []byte("\x11\x00\x80")},
	{math.MinInt16 - 1, []byte("\x12\xff\x7f\xff\xff")},

	// int32
	{math.MaxInt32, []byte("\x12\xff\xff\xff\x7f")},
	{math.MaxInt32 + 1, []byte("\x13\x00\x00\x00\x80\x00\x00\x00\x00")},
	{math.MinInt32, []byte("\x12\x00\x00\x00\x80")},
	{math.MinInt32 - 1, []byte("\x13\xff\xff\xff\x7f\xff\xff\xff\xff")},

	// int64
	{math.MaxInt64, []byte("\x13\xff\xff\xff\xff\xff\xff\xff\x7f")},
	{math.MinInt64, []byte("\x13\x00\x00\x00\x00\x00\x00\x00\x80")},
}

// Binson BOOLEAN internal representation test data table
var boolTable = []struct {
	val bool
	raw []byte
}{
	{true, []byte("\x44")},
	{false, []byte("\x45")},
}

// Binson DOUBLE internal representation test data table
var doubleTable = []struct {
	val float64
	raw []byte
}{
	{0.0, []byte("\x46\x00\x00\x00\x00\x00\x00\x00\x00")},
	{math.Copysign(0, -1), []byte("\x46\x00\x00\x00\x00\x00\x00\x00\x80")},
	{+3.1415e+10, []byte("\x46\x00\x00\x00\x6f\xeb\x41\x1d\x42")},
	{-3.1415e-10, []byte("\x46\xfc\x17\xac\xd2\x95\x96\xf5\xbd")},
	{math.NaN(), []byte("\x46\x00\x00\x00\x00\x00\x00\xf8\x7f")},
	{math.Inf(+1), []byte("\x46\x00\x00\x00\x00\x00\x00\xf0\x7f")},
	{math.Inf(-1), []byte("\x46\x00\x00\x00\x00\x00\x00\xf0\xff")},
}

// Binson STRING internal representation test data table
var stringTable = []struct {
	val string
	raw []byte
}{
	{"", []byte("\x14\x00")},
	{"abc", []byte("\x14\x03\x61\x62\x63")},
	{"größer", []byte("\x14\x08\x67\x72\xc3\xb6\xc3\x9f\x65\x72")},
}

// Binson BYTES internal representation test data table
var bytesTable = []struct {
	val []byte
	raw []byte
}{
	{[]byte(""), []byte("\x18\x00")},
	{[]byte("\x00"), []byte("\x18\x01\x00")},
	{[]byte("\x00\x01\x02\xff\x00"), []byte("\x18\x05\x00\x01\x02\xff\x00")},
}

func TestTableInts(t *testing.T) {
	for _, record := range intTable {
		b := make([]byte, 100)
		enc := NewEncoderFromBytes(b)

		// test Encoder

		enc.Integer(record.val)
		enc.Flush()
		if !bytes.Equal(record.raw, b[:len(record.raw)]) {
			t.Errorf("Binson int encoder failed: val %v, expected 0x%v != recieved: 0x%v",
				record.val, hex.EncodeToString(record.raw), hex.EncodeToString(b[:len(record.raw)]))
		}

		// test Decoder
		typeBeforeValue := record.raw[0]
		dec := NewDecoderFromBytes(record.raw[1:])
		dec.parseValue(typeBeforeValue, 0)

		if record.val != dec.ValueInteger {
			t.Errorf("Binson int decoder failed: expected %v != recieved: %v", record.val, dec.ValueInteger)
		}
	}
}

func TestTableBooleans(t *testing.T) {
	for _, record := range boolTable {
		b := make([]byte, 100)
		enc := NewEncoderFromBytes(b)

		// test Encoder
		enc.Bool(record.val)
		enc.Flush()

		slice := b[:len(record.raw)]
		if !bytes.Equal(record.raw, slice) {
			t.Errorf("Binson boolean encoder failed: val %v, expected 0x%v != recieved: 0x%v",
				record.val, hex.EncodeToString(record.raw), hex.EncodeToString(slice))
		}

		// test Decoder
		typeBeforeValue := record.raw[0]
		dec := NewDecoderFromBytes(record.raw[1:])

		dec.parseValue(typeBeforeValue, 0)

		if record.val != dec.ValueBoolean {
			t.Errorf("Binson boolean decoder failed: expected %v != recieved: %v", record.val, dec.ValueBoolean)
		}
	}
}

func TestTableDoubles(t *testing.T) {
	for _, record := range doubleTable {
		b := make([]byte, 100)

		// Encoder
		enc := NewEncoderFromBytes(b)
		enc.Double(record.val)
		enc.Flush()
		slice := b[:len(record.raw)]

		if !bytes.Equal(record.raw, slice) && !math.IsNaN(record.val) {
			t.Errorf("Binson double encoder failed: val %v, expected 0x%v != got: 0x%v",
				record.val, hex.EncodeToString(record.raw), hex.EncodeToString(slice))
		}

		// Decoder
		typeBeforeValue := record.raw[0]
		dec := NewDecoderFromBytes(record.raw[1:])

		dec.parseValue(typeBeforeValue, 0)

		if record.val != dec.ValueDouble && !math.IsNaN(record.val) {
			t.Errorf("Binson double decoder failed: expected %v != recieved: %v", record.val, dec.ValueDouble)
		}
	}
}

func TestTableStrings(t *testing.T) {
	for _, record := range stringTable {
		b := make([]byte, 100)
		enc := NewEncoderFromBytes(b)

		enc.String(record.val)
		enc.Flush()
		slice := b[:len(record.raw)]
		if !bytes.Equal(record.raw, slice) {
			t.Errorf("Encoder failed: val %v, expected 0x%v != recieved: 0x%v",
				record.val, hex.EncodeToString(record.raw), hex.EncodeToString(slice))
		}

		// test Decoder
		typeBeforeValue := record.raw[0]
		dec := NewDecoderFromBytes(record.raw[1:])

		dec.parseValue(typeBeforeValue, 0)

		if record.val != string(dec.ValueBytes) {
			t.Errorf("Decoder failed: expected %v != recieved: %v", record.val, dec.ValueBytes)
		}
	}
}

func TestTableBytes(t *testing.T) {
	for _, record := range bytesTable {
		b := make([]byte, 100)
		enc := NewEncoderFromBytes(b)
		enc.Bytes(record.val)
		enc.Flush()
		slice := b[:len(record.raw)]

		if !bytes.Equal(record.raw, slice) {
			t.Errorf("Binson bytes encoder failed: val %v, expected 0x%v != recieved: 0x%v",
				record.val, hex.EncodeToString(record.raw), hex.EncodeToString(slice))
		}

		// test Decoder
		typeBeforeValue := record.raw[0]
		dec := NewDecoderFromBytes(record.raw[1:])

		dec.parseValue(typeBeforeValue, 0)

		if !bytes.Equal(record.val, dec.ValueBytes) {
			t.Errorf("Binson bytes decoder failed: expected %v != recieved: %v",
				hex.EncodeToString(record.val), hex.EncodeToString(dec.ValueBytes))
		}
	}
}
