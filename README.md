binson-go-tiny/README.md
=========================

A light-weight one-file Golang implementation of a Binson parser (decoder) and writer (encoder).

Binson is like JSON, but faster, binary and even simpler.
See [binson.org](https://binson.org/).

This library is a high-performance, low-level library suitable for embedded targets 
and the TinyGo compiler. The dependencies are limited, code size is small, and no dynamic
memory allocation is required to use the library.

STATUS: Work on progress.



WORK NOTES
==========

NEXT STEP. Investigated heap allocs. Got:

    perf$ tinygo run -print-allocs=. main.go 
    binson.go:318:31: object allocated on the heap: escapes at line 318
    binson.go:322:31: object allocated on the heap: escapes at line 322
    binson.go:95:17: object allocated on the heap: escapes at line 95
    ../task_stack.go:67:15: object allocated on the heap: object size 65536 exceeds maximum stack allocation size 256
    ../task_stack.go:102:12: object allocated on the heap: escapes at line 104

TO DO: remove string() allocation. Lines 318, 322. Add Decoder function getStringValue() 
that returns a slice of the string. And: getBytesValue(), also a slice. 
Add fields: Decoder.valueOffset, valueSize. Set them instead of allocating memory.


## Tasks

* DONE! Must remove depency on test library for "tinygo test" to work.


## VS Code tip

Open Folder "binson" in one VS Code window and "perf" in another window.
