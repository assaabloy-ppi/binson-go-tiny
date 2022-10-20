package main

import (
	"bytes"
	"fmt"

	"github.com/assaabloy-ppi/binson-go/binson"
)

func example1() {
	//
	// {"a":123, "s":"Hello world!"}
	//
	var b bytes.Buffer
	var e = binson.NewEncoder(&b)

	e.Begin()
	e.Name("a")
	e.Integer(123)
	e.Name("s")
	e.String("Hello world!")
	e.End()
	e.Flush()

	var d = binson.NewDecoder(&b)

	d.Field("a")
	fmt.Println(d.Value) // -> 123
	d.Field("s")
	fmt.Println(d.Value) // -> Hello world!
}

func example2() {
	//
	// {"a":{"b":2},"c":3}
	//
	var b bytes.Buffer
	var e = binson.NewEncoder(&b)

	e.Begin()
	e.Name("a")
	e.Begin()
	e.Name("b")
	e.Integer(2)
	e.End()
	e.Name("c")
	e.Integer(3)
	e.End()
	e.Flush()

	var d = binson.NewDecoder(&b)

	d.Field("a")
	d.GoIntoObject()
	d.Field("b")
	fmt.Println(d.Value) // -> 2
	d.GoUpToObject()
	d.Field("c")
	fmt.Println(d.Value) // -> 3
}

func example3() {
	//
	// {"arr":[123, "hello"]}
	//
	var b bytes.Buffer
	var e = binson.NewEncoder(&b)

	e.Begin()
	e.Name("arr")
	e.BeginArray()
	e.Integer(123)
	e.String("hello")
	e.EndArray()
	e.End()
	e.Flush()

	var d = binson.NewDecoder(&b)

	d.Field("arr")
	d.GoIntoArray()
	gotValue := d.NextArrayValue()
	fmt.Println(gotValue)                      // -> true
	fmt.Println(binson.Integer == d.ValueType) // -> true
	fmt.Println(d.Value)                       // -> 123

	gotValue = d.NextArrayValue()
	fmt.Println(gotValue)                     // -> true
	fmt.Println(binson.String == d.ValueType) // -> true
	fmt.Println(d.Value)                      // -> hello
}

func main() {
	example1()
	fmt.Println()

	example2()
	fmt.Println()

	example3()
	fmt.Println()
}
