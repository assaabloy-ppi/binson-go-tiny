// perf - Application that tests performance of Binson parsing.
// Runs on Raspberry Pi Pico as the reference embedded target.
// Works on a PC as well.
//
// STATUS: Work in progress.
//

package main

import "binson"

var buffer = [...]byte{0x40, 0x14, 0x03, 0x63, 0x69, 0x64, 0x10, 0x26, 0x14, 0x01, 0x7a, 0x40, 0x41, 0x41}

func main() {
	//println("Hello from perf/main.go")

	// {"cid":38, "z":{}}
	d := binson.NewDecoderFromBytes(buffer[:])

	d.NextField()
	println("field: " + string(d.Name))
}
