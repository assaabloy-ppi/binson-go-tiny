package binson

import "fmt"

func Example1() {
	//
	// {"a":123, "s":"Hello world!"}
	//

	b := make([]byte, 100)
	e := NewEncoderFromBytes(b)

	e.Begin()
	e.Name("a")
	e.Integer(123)
	e.Name("s")
	e.String("Hello world!")
	e.End()
	e.Flush()

	var d = NewDecoderFromBytes(b)

	//d.Field("a")
	//fmt.Println(d.Value)
	// xxx Output: 123

	d.Field("s")
	fmt.Println(string(d.Value.([]byte)))
	// Output: Hello world!
}
