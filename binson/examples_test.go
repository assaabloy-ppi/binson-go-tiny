// Executable code examples, run with ordinary "go test" command.
// Note, unlike the other tests, these examples are not run with "tinygo test".
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
	d.Field("s")
	fmt.Println(string(d.BytesValue))
	// Output: Hello world!
}
