// Executable code examples, run with ordinary "go test" command.
// Note, unlike the other tests, these examples are not executed
// with "tinygo test".

package binson

import "fmt"

func Example1() {
	//
	// {"a":123, "s":"Hello world!"}
	//

	buf := make([]byte, 100)

	e := Encoder{}
	e.Init(buf)
	e.Begin()
	e.Name("a")
	e.Integer(123)
	e.Name("s")
	e.String("Hello world!")
	e.End()

	d := Decoder{}
	d.Init(buf)
	d.Field("a")

	fmt.Println(d.ValueInteger)
	// Output: 123
}

func Example2() {
	//
	// {"a":123, "s":"Hello world!"}
	//

	buf := make([]byte, 100)

	e := Encoder{}
	e.Init(buf)
	e.Begin()
	e.Name("a")
	e.Integer(123)
	e.Name("s")
	e.String("Hello world!")
	e.End()

	d := Decoder{}
	d.Init(buf)
	d.Field("s") // Skips field "a".

	fmt.Println(string(d.ValueBytes))
	// Output: Hello world!
}
