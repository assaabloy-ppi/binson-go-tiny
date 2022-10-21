// perf - Application that tests performance of Binson parsing.
// Runs on Raspberry Pi Pico as the reference embedded target.
// Works on a PC as well.
//
// STATUS: Work in progress.
//

package main

import "binson"

func main() {
	println("Hello from perf/main.go")

	// {"cid":38, "z":{}}
	d := binson.NewDecoderFromBytes([]byte("\x40\x14\x03\x63\x69\x64\x10\x26\x14\x01\x7a\x40\x41\x41"))

	d.NextField()
	println("field: " + string(d.Name))
}
