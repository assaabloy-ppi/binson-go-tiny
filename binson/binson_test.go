package binson

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncoderEmptyBinsonObject(t *testing.T) {
	exp := []byte("\x40\x41") // {}
	b := make([]byte, 100)
	e := NewEncoderFromBytes(b)

	e.Begin()
	e.End()
	e.Flush()

	if e.offset != 2 {
		t.Errorf("unexpected e.offset: %d", e.offset)
		return
	}

	slice := b[:len(exp)]
	if !bytes.Equal(exp, slice) {
		t.Errorf("Binson encoder failure: expected 0x%v, got %v", hex.EncodeToString(exp),
			hex.EncodeToString(slice))
	}
}

func TestEncoderEmptyBinsonArray(t *testing.T) {
	exp := []byte("\x42\x43") // []
	b := make([]byte, 10)
	e := NewEncoderFromBytes(b)

	e.BeginArray()
	e.EndArray()
	e.Flush()

	if !bytes.Equal(exp, b[:len(exp)]) {
		t.Errorf("Binson encoder failure: expected 0x%v", hex.EncodeToString(exp))
	}
}

func TestEncoderEmptyBinsonArray2(t *testing.T) {
	exp := []byte("\x40\x14\x00\x42\x43\x41") // {""=[]}
	b := make([]byte, 100)
	e := NewEncoderFromBytes(b)

	e.Begin()
	e.Name("")
	e.BeginArray()
	e.EndArray()
	e.End()
	e.Flush()

	if !bytes.Equal(exp, b[:len(exp)]) {
		t.Errorf("Binson encoder failure: expected 0x%v", hex.EncodeToString(exp))
	}
}

func TestEncoderObjectWithUTF8Name(t *testing.T) {
	exp := []byte("\x40\x14\x06\xe7\x88\x85\xec\x9b\xa1\x10\x7b\x41") // {"爅웡":123}
	b := make([]byte, 100)
	e := NewEncoderFromBytes(b)

	e.Begin()
	e.Name("爅웡")
	e.Integer(123)
	e.End()
	e.Flush()

	if !bytes.Equal(exp, b[:len(exp)]) {
		t.Errorf("Binson encoder failure: expected 0x%v", hex.EncodeToString(exp))
	}
}

func TestEncoderNestedObjectsWithEmptyKeyNames(t *testing.T) {
	// {"":{"":{"":{}}}}
	exp := []byte("\x40\x14\x00\x40\x14\x00\x40\x14\x00\x40\x41\x41\x41\x41")
	b := make([]byte, 100)
	e := NewEncoderFromBytes(b)

	e.Begin()
	e.Name("")
	e.Begin()
	e.Name("")
	e.Begin()
	e.Name("")
	e.Begin()
	e.End()
	e.End()
	e.End()
	e.End()
	e.Flush()

	if !bytes.Equal(exp, b[:len(exp)]) {
		t.Errorf("Binson encoder failure: expected 0x%v", hex.EncodeToString(exp))
	}
}

func TestEncoderNestedArraysAsObjectValue(t *testing.T) {
	// {"b":[[[]]]}
	exp := []byte("\x40\x14\x01\x62\x42\x42\x42\x43\x43\x43\x41")
	b := make([]byte, 100)
	e := NewEncoderFromBytes(b)

	e.Begin()
	e.Name("b")
	e.BeginArray()
	e.BeginArray()
	e.BeginArray()
	e.EndArray()
	e.EndArray()
	e.EndArray()
	e.End()
	e.Flush()

	if !bytes.Equal(exp, b[:len(exp)]) {
		t.Errorf("Binson encoder failure: expected 0x%v", hex.EncodeToString(exp))
	}
}

func TestEncoderNestedStructures1AsObjectValue(t *testing.T) {
	// {"b":[[],{},[]]}
	exp := []byte("\x40\x14\x01\x62\x42\x42\x43\x40\x41\x42\x43\x43\x41")
	b := make([]byte, 100)
	e := NewEncoderFromBytes(b)

	e.Begin()
	e.Name("b")
	e.BeginArray()
	e.BeginArray()
	e.EndArray()
	e.Begin()
	e.End()
	e.BeginArray()
	e.EndArray()
	e.EndArray()
	e.End()
	e.Flush()

	if !bytes.Equal(exp, b[:len(exp)]) {
		t.Errorf("Binson encoder failure: expected 0x%v", hex.EncodeToString(exp))
	}
}

func TestEncoderNestedStructures2AsObjectValue(t *testing.T) {
	// {"b":[[{}],[{}]]}
	exp := []byte("\x40\x14\x01\x62\x42\x42\x40\x41\x43\x42\x40\x41\x43\x43\x41")
	b := make([]byte, 100)
	e := NewEncoderFromBytes(b)

	e.Begin()
	e.Name("b")
	e.BeginArray()
	e.BeginArray()
	e.Begin()
	e.End()
	e.EndArray()
	e.BeginArray()
	e.Begin()
	e.End()
	e.EndArray()
	e.EndArray()
	e.End()
	e.Flush()

	if !bytes.Equal(exp, b[:len(exp)]) {
		t.Errorf("Binson encoder failure: expected 0x%v", hex.EncodeToString(exp))
	}
}

func TestEncoderComplexObjectStructure1(t *testing.T) {
	// {"abc":{"cba":{}}, "b":{"abc":{}}}
	exp := []byte("\x40\x14\x03\x61\x62\x63\x40\x14\x03\x63\x62\x61\x40\x41\x41\x14\x01\x62\x40\x14\x03\x61\x62\x63\x40\x41\x41\x41")
	b := make([]byte, 100)
	e := NewEncoderFromBytes(b)

	e.Begin()
	e.Name("abc")
	e.Begin()
	e.Name("cba")
	e.Begin()
	e.End()
	e.End()
	e.Name("b")
	e.Begin()
	e.Name("abc")
	e.Begin()
	e.End()
	e.End()
	e.End()
	e.Flush()

	if !bytes.Equal(exp, b[:len(exp)]) {
		t.Errorf("Binson encoder failure: expected 0x%v", hex.EncodeToString(exp))
	}
}

func TestEncoderComplexObjectStructure2(t *testing.T) {
	// {"b":[true,13,"cba",{"abc":false, "b":"0x008100ff00", "cba":"abc"},9223372036854775807]}
	exp := []byte(
		"\x40\x14\x01\x62\x42\x44\x10\x0d\x14\x03\x63\x62\x61\x40\x14\x03" +
			"\x61\x62\x63\x45\x14\x01\x62\x18\x05\x00\x81\x00\xff\x00\x14\x03" +
			"\x63\x62\x61\x14\x03\x61\x62\x63\x41\x13\xff\xff\xff\xff\xff\xff" +
			"\xff\x7f\x43\x41",
	)
	b := make([]byte, 100)
	e := NewEncoderFromBytes(b)

	e.Begin()
	e.Name("b")
	e.BeginArray()
	e.Bool(true)
	e.Integer(13)
	e.String("cba")
	e.Begin()
	e.Name("abc")
	e.Bool(false)
	e.Name("b")
	e.Bytes([]byte("\x00\x81\x00\xff\x00"))
	e.Name("cba")
	e.String("abc")
	e.End()
	e.Integer(9223372036854775807)
	e.EndArray()
	e.End()
	e.Flush()

	if !bytes.Equal(exp, b[:len(exp)]) {
		t.Errorf("Binson encoder failure: expected 0x%v", hex.EncodeToString(exp))
	}
}

func TestDecoderObjectEmpty(t *testing.T) {
	d := NewDecoderFromBytes([]byte("\x40\x41"))
	gotField := d.NextField()
	assert.Equal(t, false, gotField)

	if d.err != nil {
		t.Errorf("Binson decoder error: %v", d.err)
	}
}

func TestDecoder0(t *testing.T) {
	// {"cid":38, "z":{}}
	d := NewDecoderFromBytes([]byte("\x40\x14\x03\x63\x69\x64\x10\x26\x14\x01\x7a\x40\x41\x41"))

	gotField := d.NextField()
	assert.Equal(t, true, gotField)
	assert.Equal(t, Integer, d.ValueType)
	assert.Equal(t, "cid", string(d.Name))
	assert.Equal(t, int64(38), d.Value)

	gotField = d.NextField()
	assert.Equal(t, true, gotField)
	assert.Equal(t, Object, d.ValueType)
	assert.Equal(t, "z", string(d.Name))

	gotField = d.NextField()
	assert.Equal(t, false, gotField)

	if d.err != nil {
		t.Errorf("Binson decoder error: %v", d.err)
	}
}

func TestDecoderNested1(t *testing.T) {
	// {"a":{"b":2}}
	d := NewDecoderFromBytes([]byte("\x40\x14\x01\x61\x40\x14\x01\x62\x10\x02\x41\x41"))

	gotField := d.NextField()
	assert.Equal(t, true, gotField)
	assert.Equal(t, Object, d.ValueType)
	assert.Equal(t, "a", string(d.Name))

	d.GoIntoObject()
	gotField = d.NextField()
	assert.Equal(t, true, gotField)
	assert.Equal(t, Integer, d.ValueType)
	assert.Equal(t, "b", string(d.Name))
	assert.Equal(t, int64(2), d.Value)
	d.GoUpToObject()

	gotField = d.NextField()
	assert.Equal(t, false, gotField)

	if d.err != nil {
		t.Errorf("Binson decoder error: %v", d.err)
	}
}

func TestDecoderExample4a(t *testing.T) {
	// {"a":1,"b":{"c":3},"d":4}
	d := NewDecoderFromBytes([]byte("\x40\x14\x01\x61\x10\x01\x14\x01\x62\x40\x14\x01\x63\x10\x03\x41\x14\x01\x64\x10\x04\x41"))

	gotField := d.NextField()
	assert.Equal(t, true, gotField)
	assert.Equal(t, "a", string(d.Name))
	assert.Equal(t, Integer, d.ValueType)
	assert.Equal(t, int64(1), d.Value)

	gotField = d.NextField()
	assert.Equal(t, true, gotField)
	assert.Equal(t, Object, d.ValueType)
	assert.Equal(t, "b", string(d.Name))

	gotField = d.NextField()
	assert.Equal(t, true, gotField)
	assert.Equal(t, "d", string(d.Name))
	assert.Equal(t, Integer, d.ValueType)
	assert.Equal(t, int64(4), d.Value)

	if d.err != nil {
		t.Errorf("Binson decoder error: %v", d.err)
	}
}

func TestDecoderExample4b(t *testing.T) {
	// {"a":1,"b":{"c":3},"d":4}
	d := NewDecoderFromBytes([]byte("\x40\x14\x01\x61\x10\x01\x14\x01\x62\x40\x14\x01\x63\x10\x03\x41\x14\x01\x64\x10\x04\x41"))

	gotField := d.NextField()
	gotField = d.NextField()

	d.GoIntoObject()
	gotField = d.NextField()
	assert.Equal(t, true, gotField)
	assert.Equal(t, "c", string(d.Name))
	assert.Equal(t, int64(3), d.Value)
	d.GoUpToObject()

	gotField = d.NextField()
	assert.Equal(t, true, gotField)
	assert.Equal(t, "d", string(d.Name))
	assert.Equal(t, int64(4), d.Value)

	assert.Equal(t, false, d.NextField())

	if d.err != nil {
		t.Errorf("Binson decoder error: %v", d.err)
	}
}

func TestDecoderExample4c(t *testing.T) {
	// {"a":1,"b":{"c":3},"d":4}
	d := NewDecoderFromBytes([]byte("\x40\x14\x01\x61\x10\x01\x14\x01\x62\x40\x14\x01\x63\x10\x03\x41\x14\x01\x64\x10\x04\x41"))

	d.Field("b")
	d.GoIntoObject()
	d.Field("c")
	assert.Equal(t, int64(3), d.Value)
	d.GoUpToObject()
	d.Field("d")
	assert.Equal(t, int64(4), d.Value)

	if d.err != nil {
		t.Errorf("Binson decoder error: %v", d.err)
	}
}

func TestDecoderNonExistantField(t *testing.T) {
	// {"cid":38, "z":{}}
	d := NewDecoderFromBytes([]byte("\x40\x14\x03\x63\x69\x64\x10\x26\x14\x01\x7a\x40\x41\x41"))
	assert.Equal(t, false, d.Field("height"))
}

func TestDecoderExampleArray1(t *testing.T) {
	// {"a":[1, "hello"]}
	d := NewDecoderFromBytes([]byte("\x40\x14\x01\x61\x42\x10\x01\x14\x05\x68\x65\x6c\x6c\x6f\x43\x41"))

	d.Field("a")
	d.GoIntoArray()

	gotField := d.NextArrayValue()
	assert.Equal(t, true, gotField)
	assert.Equal(t, Integer, d.ValueType)
	assert.Equal(t, int64(1), d.Value)

	gotField = d.NextArrayValue()
	assert.Equal(t, true, gotField)
	assert.Equal(t, String, d.ValueType)
	assert.Equal(t, []byte("hello"), d.Value)

	d.GoUpToArray()

	if d.err != nil {
		t.Errorf("Binson decoder error: %v", d.err)
	}
}

func TestDecoderSkipArrayFields(t *testing.T) {
	// {"a":1,"b":[10,20],"c":3}
	d := NewDecoderFromBytes([]byte("\x40\x14\x01\x61\x10\x01\x14\x01\x62\x42\x10\x0a\x10\x14\x43\x14\x01\x63\x10\x03\x41"))

	d.Field("a")
	assert.Equal(t, Integer, d.ValueType)
	assert.Equal(t, int64(1), d.Value)

	d.Field("c")
	assert.Equal(t, Integer, d.ValueType)
	assert.Equal(t, int64(3), d.Value)

	if d.err != nil {
		t.Errorf("Binson decoder error: %v", d.err)
	}
}

func TestDecoderFieldInTheMiddle1(t *testing.T) {
	// {"a":1,"b":[10,20],"c":3}
	d := NewDecoderFromBytes([]byte("\x40\x14\x01\x61\x10\x01\x14\x01\x62\x42\x10\x0a\x10\x14\x43\x14\x01\x63\x10\x03\x41"))

	d.Field("b")
	d.GoIntoArray()

	d.NextArrayValue()
	assert.Equal(t, int64(10), d.Value)

	d.NextArrayValue()
	assert.Equal(t, int64(20), d.Value)

	d.GoUpToObject()
	d.Field("c")
	assert.Equal(t, int64(3), d.Value)

	if d.err != nil {
		t.Errorf("Binson decoder error: %v", d.err)
	}
}

func TestDecoderArrayInArray1(t *testing.T) {
	// {"a":1,"b":[10,[100,101],20],"c":3}
	b := []byte(
		"\x40\x14\x01\x61\x10\x01\x14\x01\x62\x42\x10\x0a\x42" +
			"\x10\x64\x10\x65\x43\x10\x14\x43\x14\x01\x63\x10\x03\x41")
	d := NewDecoderFromBytes(b)

	d.Field("b")
	d.GoIntoArray()

	gotValue := d.NextArrayValue()
	assert.Equal(t, true, gotValue)
	assert.Equal(t, int64(10), d.Value)

	gotValue = d.NextArrayValue()
	assert.Equal(t, true, gotValue)
	assert.Equal(t, Array, d.ValueType)

	d.GoIntoArray()
	gotValue = d.NextArrayValue()
	assert.Equal(t, true, gotValue)
	assert.Equal(t, Integer, d.ValueType)
	assert.Equal(t, int64(100), d.Value)

	gotValue = d.NextArrayValue()
	assert.Equal(t, true, gotValue)
	assert.Equal(t, Integer, d.ValueType)
	assert.Equal(t, int64(101), d.Value)

	d.GoUpToArray()

	gotValue = d.NextArrayValue()
	assert.Equal(t, true, gotValue)
	assert.Equal(t, Integer, d.ValueType)
	assert.Equal(t, int64(20), d.Value)

	if d.err != nil {
		t.Errorf("Binson decoder error: %v", d.err)
	}
}
