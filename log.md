log.md
======

Log of the project. Reverse chronological order.



2022-10-20
==========

## Status

Tests are working. Work in progress, API is expected to change.

Dynamic memory allocation removed
(need to be verified). Started with code from 
repo github.com/assaabloy-ppi/binson-go.
Removing dynamic memory allocation. In general, making the code
more suitable for an embedded target. The Raspberry Pi Pico
is currently used as the reference target (TinyGo).

To do: verify zero dyn allocs. Make examples executable. And more.
This is work in progress.


## Naming: binson-go-tiny

Going with binson-go-tiny as the repo name. binson-go-light is possibly 
confusing since binson-go is based on binson-java-light. The "tiny" 
associates with TinyGo and with tiny targets. And this lib is for tiny
targets!



2022-10-16
==========

## WORK NOTES 2022-10-16 14.00

Frans. Removed dynalloc, should not use any dynamic memory allocation now.
How to test this?


## WORK NOTES 2022-10-16 12.00

Frans. Starting with binson-go. Change to using []byte buffers for 
input and output. Removed dependencies on fmt, math, encoding/binary.
Next step: remove dynamic memory allocation (dynalloc). For bytes/string types
and field names. Use nameOffset, nameSize, valueOffset, valueSize. Get functions can
write to memory provided by caller if needed.
Note, slices can be returned to the caller (lib user). That should not result in
dynalloc.

