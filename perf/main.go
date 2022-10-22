// perf - Application that tests performance of Binson parsing.
// Runs on Raspberry Pi Pico as the reference embedded target.
// Works on a PC as well.
//
// STATUS: Work in progress.
//

package main

import (
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/assaabloy-ppi/binson-go-tiny/binson"
)

func main() {
	time.Sleep(time.Second * 4)

	pln("==== perf ====")
	pln("binson-go-tiny/perf/main.go")

	var arg1 string

	if len(os.Args) > 1 {
		arg1 = os.Args[1]
	} else {
		arg1 = "perf"
	}

	switch arg1 {
	case "req":
		reqRepeated()
	case "perf":
		perf()
	default:
		mypanic("unexpected arg1 " + arg1)
	}
}

func perf() {
	pln("perf() called")

	var batchSize int = 500
	var count int = 0
	var duration time.Duration

	inbuf := getReq()
	outbuf := make([]byte, 100)
	var d binson.Decoder = binson.Decoder{}
	var e binson.Encoder = binson.Encoder{}

	for {
		t0 := time.Now()

		for i := 0; i < batchSize; i++ {
			d.Init(inbuf)
			e.Init(outbuf)
			handleReq(d, e)

			if i%20 == 19 {
				runtime.GC()
			}
		}
		count += batchSize

		duration = time.Since(t0)
		micros := duration.Microseconds()
		perReq := float64(micros) / float64(batchSize)

		pln("** Total count: " + strconv.Itoa(count))
		pln("batch size: " + strconv.Itoa(batchSize))
		pln("Batch duration (us): " + strconv.Itoa(int(duration.Microseconds())))
		pln("PerReq: " + strconv.FormatFloat(perReq, 'e', 4, 64) + " micros")

		if duration.Microseconds() < 1_000_000 {
			batchSize = batchSize * 2
		}
	}
}

var getTimeBytes = []byte("getTime")

// Mock for handling an RPC request. The test request (getTime) is parsed
// and a result is serialized to Binson. For performance testing.
func handleReq(d binson.Decoder, e binson.Encoder) {
	ok := d.NextField()

	if !ok {
		mypanic("NextField !ok")
	}

	if d.ValueType != binson.String {
		mypanic("unexpected field type")
	}

	if equals(d.BytesValue, getTimeBytes) {
		handleGetTime(e)
	} else {
		mypanic("expected getTime")
	}
}

func equals(a, b []byte) bool {
	length := len(a)
	if len(b) != length {
		return false
	}

	for i := 0; i < length; i++ {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func handleGetTime(e binson.Encoder) {
	e.Begin()
	e.Name("time")
	e.String("15:04:05.000")
	e.End()
	e.Flush()
}

func reqRepeated() {
	for {
		req()
		time.Sleep(time.Second * 4)
	}
}

func req() {
	pln("== req(): prints the test request")

	buf := make([]byte, 100)
	e := binson.NewEncoderFromBytes(buf)

	e.Begin()
	e.Name("a")
	e.String("getTime")
	e.End()
	e.Flush()

	len := e.Offset()
	pln("Size: " + strconv.Itoa(len))

	for i := range buf {
		print(strconv.Itoa(int(buf[i])) + ", ")
	}
	pln("")

	// RESULT
	// Size: 14
	// 64, 20, 1, 97, 20, 7, 103, 101, 116, 84, 105, 109, 101, 65,
}

// Returns the test request bytes
func getReq() []byte {
	// {a="getTime"}
	return []byte{64, 20, 1, 97, 20, 7, 103, 101, 116, 84, 105, 109, 101, 65}
}

func mypanic(s string) {
	pln(s)
	panic(s)
}

// Println for this application.
// Minicom etc prefers \r\n for line ending.
func pln(s string) {
	print(s + "\r\n")
}

func print(s string) {
	_, err := os.Stdout.WriteString(s)
	if err != nil {
		panic(err)
	}
}
