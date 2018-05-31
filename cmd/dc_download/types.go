package main

/*
#include <stdlib.h>
#include "libdivecomputer/device.h"
#include "libdivecomputer/parser.h"

typedef struct {
	unsigned int number;
	unsigned int size;
	char   *fingerprint;
	dc_datetime_t datetime;
	unsigned int divetime;
	double   maxdepth;
	char   *divemode;
	dc_gasmix_t **gasmixes;
	unsigned int ngasmixes;
	dc_tank_t  **tanks;
	unsigned int ntanks;
	dc_salinity_t *salinity;
	double atmosphere;
} DIVE_RECORD;

typedef struct {
	dc_device_t *device;
	unsigned int number;
	DIVE_RECORD *dives;
} DIVE_DATA;
*/
import "C"

import (
	"unsafe"

	"github.com/calle-gunnarsson/libdivecomputer-go/core"
)

type diveData struct {
	Device *core.Device `json:"-"`
	Number uint32
	Dives  *[]Dive

	ref *C.DIVE_DATA
}

func newDiveData() *diveData {
	return (*diveData)(allocDiveDataMemory(1))
}

func NewDiveDataRef(ref unsafe.Pointer) *diveData {
	if ref == nil {
		return nil
	}
	obj := new(diveData)
	obj.ref = (*C.DIVE_DATA)(unsafe.Pointer(ref))
	return obj
}

const sizeOfDiveDataValue = unsafe.Sizeof([1]C.DIVE_DATA{})

func allocDiveDataMemory(n int) unsafe.Pointer {
	mem, err := C.calloc(C.size_t(n), (C.size_t)(sizeOfDiveDataValue))
	if err != nil {
		panic("memory alloc error: " + err.Error())
	}
	return mem
}

func (x *diveData) PassRef() *C.DIVE_DATA {
	if x == nil {
		return nil
	} else if x.ref != nil {
		return x.ref
	}
	mem := allocDiveDataMemory(1)
	x.ref = (*C.DIVE_DATA)(mem)
	x.ref.number = (C.uint)(x.Number)
	x.ref.device = (*C.dc_device_t)(unsafe.Pointer(x.Device))
	x.ref.dives = (*C.DIVE_RECORD)(unsafe.Pointer(x.Dives))

	return x.ref
}

func (x *diveData) Deref() {
	if x.ref == nil {
		return
	}
	x.Device = (*core.Device)(unsafe.Pointer(x.ref.device))
	x.Number = (uint32)(x.ref.number)
	*x.Dives = (*[1 << 30]Dive)(unsafe.Pointer(x.ref.dives))[:x.ref.number:x.ref.number]
}

func (x *diveData) Ref() *C.DIVE_DATA {
	if x == nil {
		return nil
	}
	return (*C.DIVE_DATA)(unsafe.Pointer(x))
}

// Free cleanups the referenced memory using C free.
func (x *diveData) Free() {
	if x != nil {
		for _, d := range *x.Dives {
			d.Free()
		}
		C.free(unsafe.Pointer(x))
	}
}

type Dive struct {
	Number      uint32
	Size        uint32
	Fingerprint string
	Datetime    core.Datetime
	DiveTime    uint32
	MaxDepth    float64
	DiveMode    string
	GasMixes    []core.Gasmix
	Tanks       []core.Tank
	Salinity    core.Salinity
	Atmosphere  float64

	ref *C.DIVE_RECORD
}

func newDive() *Dive {
	return (*Dive)(allocDiveMemory(1))
}

func NewDiveRef(ref unsafe.Pointer) *Dive {
	if ref == nil {
		return nil
	}
	obj := new(Dive)
	obj.ref = (*C.DIVE_RECORD)(unsafe.Pointer(ref))
	return obj
}

const sizeOfDiveValue = unsafe.Sizeof([1]C.DIVE_RECORD{})

func allocDiveMemory(n int) unsafe.Pointer {
	mem, err := C.calloc(C.size_t(n), (C.size_t)(sizeOfDiveValue))
	if err != nil {
		panic("memory alloc error: " + err.Error())
	}
	return mem
}

func (x *Dive) PassRef() *C.DIVE_RECORD {
	if x == nil {
		return nil
	} else if x.ref != nil {
		return x.ref
	}
	mem := allocDiveMemory(1)
	x.ref = (*C.DIVE_RECORD)(mem)
	x.ref.number = (C.uint)(x.Number)

	return x.ref
}

func (x *Dive) Deref() {
	if x.ref == nil {
		return
	}
	x.Number = (uint32)(x.ref.number)
}

func (x *Dive) Ref() *C.DIVE_RECORD {
	if x == nil {
		return nil
	}
	return (*C.DIVE_RECORD)(unsafe.Pointer(x))
}

func (x *Dive) Free() {
	if x != nil {
		for _, g := range x.GasMixes {
			g.Free()
		}
		for _, t := range x.Tanks {
			t.Free()
		}
		x.Datetime.Free()
		x.Salinity.Free()

		C.free(unsafe.Pointer(x.ref))
	}
}
