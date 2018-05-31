package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"unsafe"

	core "github.com/calle-gunnarsson/libdivecomputer-go/core"
	"github.com/calle-gunnarsson/libdivecomputer-go/internal/pkg/helpers"
)

func diveCb(data string, size uint32, fingerprint string, fsize uint32, userdata unsafe.Pointer) int32 {
	var parser *core.Parser

	dData := (*diveData)(unsafe.Pointer(userdata))
	dData.Number++

	fmt.Printf("Dive: number=%d, size=%d, fingerprint=", dData.Number, size)
	if fsize > 0 {
		fmt.Printf("%X", []byte(fingerprint))
	}
	fmt.Println()

	// Create the parser.
	fmt.Println("Creating the parser.")
	rc := core.ParserNew(&parser, core.NewDeviceRef(unsafe.Pointer(dData.Device.Ref())))
	if helpers.CheckError("Error when creating the parser", rc) {
		return 1
	}
	defer core.ParserDestroy(parser)

	// Register the data.
	fmt.Println("Registering the data")
	rc = core.ParserSetData(parser, data, size)
	if helpers.CheckError("Error when registering the data", rc) {
		return 1
	}

	// Parse the dive data.

	fmt.Println("Parsing the dive data.")

	dive := Dive{}
	dive.Number = dData.Number
	dive.Size = size
	dive.Fingerprint = fmt.Sprintf("%X", fingerprint)

	fmt.Println("Parsing the datetime.")
	rc = core.ParserGetDatetime(parser, &dive.Datetime)
	if rc != core.DcStatusUnsupported && helpers.CheckError("Error when parsing the datetime.", rc) {
		return 1
	}
	dive.Datetime.Deref()

	fmt.Println("Parsing the maxdepth.")
	maxdepth := C.double(0)
	rc = core.ParserGetField(parser, core.DcFieldMaxdepth, 0, unsafe.Pointer(&maxdepth))
	if rc != core.DcStatusUnsupported && helpers.CheckError("Error when parsing the max depth.", rc) {
		return 1
	}
	dive.MaxDepth = float64(maxdepth)

	fmt.Println("Parsing the dive time.")
	divetime := C.uint(0)
	rc = core.ParserGetField(parser, core.DcFieldDivetime, 0, unsafe.Pointer(&divetime))
	if rc != core.DcStatusUnsupported && helpers.CheckError("Error when parsing dive time.", rc) {
		return 1
	}
	dive.DiveTime = uint32(maxdepth)

	fmt.Println("Parsing the dive mode.")
	divemode := C.int(0)
	rc = core.ParserGetField(parser, core.DcFieldDivemode, 0, unsafe.Pointer(&divemode))
	if false && rc != core.DcStatusUnsupported && helpers.CheckError("Error when parsing divet mode.", rc) {
		return 1
	} else if rc != core.DcStatusUnsupported {
		names := []string{"freedive", "gauge", "oc", "cc"}
		dive.DiveMode = names[int32(divemode)]
	}

	fmt.Println("Parsing the gas mixes.")
	ngases := C.uint(0)
	rc = core.ParserGetField(parser, core.DcFieldGasmixCount, 0, unsafe.Pointer(&ngases))
	if rc != core.DcStatusUnsupported && helpers.CheckError("Error when parsing the gas mix count.", rc) {
		return 1
	}

	for i := uint32(0); i < uint32(ngases); i++ {
		gasmix := core.Gasmix{}
		rc = core.ParserGetField(parser, core.DcFieldGasmix, i, unsafe.Pointer(&gasmix))
		if rc != core.DcStatusUnsupported && helpers.CheckError("Error when parsing the gas mix.", rc) {
			return 1
		}

		dive.GasMixes = append(dive.GasMixes, gasmix)
	}

	fmt.Println("Parsing the tanks.")
	ntanks := uint32(0)
	rc = core.ParserGetField(parser, core.DcFieldTankCount, 0, unsafe.Pointer(&ntanks))
	if rc != core.DcStatusUnsupported && helpers.CheckError("Error when parsing the tanks count", rc) {
		return 1
	}

	for i := uint32(0); i < ntanks; i++ {
		tank := core.Tank{}

		rc = core.ParserGetField(parser, core.DcFieldTank, i, unsafe.Pointer(tank.Ref()))
		if rc != core.DcStatusUnsupported && helpers.CheckError("Error when parsing tank", rc) {
			return 1
		}
		tank.Deref()

		dive.Tanks = append(dive.Tanks, tank)
	}

	fmt.Println("Parsing the salinity")
	rc = core.ParserGetField(parser, core.DcFieldSalinity, 0, unsafe.Pointer(dive.Salinity.Ref()))
	if rc != core.DcStatusUnsupported && helpers.CheckError("Error when parsing the salinity", rc) {
		return 1
	}
	dive.Salinity.Deref()

	fmt.Println("Parsing the atmospheric pressure.")
	atmosphere := C.double(0)
	rc = core.ParserGetField(parser, core.DcFieldAtmospheric, 0, unsafe.Pointer(&atmosphere))
	if rc != core.DcStatusUnsupported && helpers.CheckError("Error when parsing the atmospheric pressure", rc) {
		return 1
	}
	dive.Atmosphere = float64(atmosphere)

	*dData.Dives = append(*dData.Dives, dive)
	dData.PassRef()

	return 1
}
