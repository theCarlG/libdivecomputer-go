libdivecomputer-go [![GoDoc](https://godoc.org/github.com/calle-gunnarsson/libdivecomputer-go/core?status.svg)](https://godoc.org/github.com/calle-gunnarsson/libdivecomputer-go/core)
========

This project provides Go bindings for libdivecomputer v0.7.0-devel - a cross-platform and open source library for communication with dive computers from various manufacturers.<br />
All the binding code has automatically been generated with rules defined in [core.yml](/core.yml). There are future plans to write a high level wrapper for the bindings

Before start you must install [libdivecomputer](https://www.libdivecomputer.org/download.html).

### Usage

```
$ go get github.com/calle-gunnarsson/libdivecomputer-go/core
```

### Demo
These are simple ports of dctool <list,dump,download> and will probably contain some memory leaks or other bugs.

```bash
# List all supported devices
$ go get github.com/calle-gunnarsson/libdivecomputer-go/cmd/dc_list
$ dc_list

# Memory dump of device
$ go get github.com/calle-gunnarsson/libdivecomputer-go/cmd/dc_dump
$ dc_dump -name d9 -family suunto -filename dump.bin -transport serial -devname /dev/ttyS1

# Save dives as json to file
$ go get github.com/calle-gunnarsson/libdivecomputer-go/cmd/dc_download
$ dc_download -name d9 -family suunto -filename dump.json -transport serial -devname /tmp/ttyS1
```

### Rebuilding the package

You will need to get the [c-for-go](https://github.com/xlab/c-for-go) tool installed first.

```
$ git clone github.com/calle-gunnarsson/libdivecomputer-go && cd libdivecomputer-go
$ make clean && make
```
