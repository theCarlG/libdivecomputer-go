package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"unsafe"

	core "github.com/calle-gunnarsson/libdivecomputer-go/core"
	"github.com/calle-gunnarsson/libdivecomputer-go/internal/pkg/helpers"
)

var cancel int32 = 0

func cancelCb(userdata unsafe.Pointer) int32 {
	return cancel
}

func main() {
	var (
		descriptor *core.Descriptor
		iostream   *core.Iostream
		context    *core.Context
		device     *core.Device
	)

	cancelChan := make(chan os.Signal)

	signal.Notify(cancelChan)
	go func() {
		sig := <-cancelChan
		fmt.Printf("caught sig: %+v\n", sig)
		cancel = 1
	}()

	name := flag.String("name", "", "Device name")
	compFamily := flag.String("family", "", "Device family type")
	model := flag.Uint("model", 0, "Device model number")
	filename := flag.String("filename", "", "Dump file")
	transportType := flag.String("transport", "", "Transport type")
	devname := flag.String("devname", "", "Tty device name")

	flag.Parse()

	if *filename == "" {
		fmt.Println("You need to specify filename")
		os.Exit(1)
	}

	transport := helpers.TransportType(*transportType)
	if transport == core.DcTransportNone {
		fmt.Println("No valid transport type specified")
		os.Exit(1)
	}

	family := helpers.FamilyType(*compFamily)
	if *name == "" && family == core.DcFamilyNull {
		fmt.Println("No device name or family type specified.")
		os.Exit(1)
	}

	// Set the default model number.
	if family != 0 && *model == uint(0) {
		*model = helpers.FamilyModel(family)
	}

	rc := helpers.DescriptorSearch(&descriptor, *name, family, uint32(*model))
	helpers.ExitOnError("Descriptor error", rc)

	if descriptor == nil {
		if name != nil {
			fmt.Printf("No supported device found: %s\n", *name)
		} else {
			fmt.Printf("No supported device found: %s, 0x%X\n", helpers.FamilyName(family), *model)
		}
		os.Exit(1)
	}

	rc = core.ContextNew(&context)
	helpers.ExitOnError("Error when creating context", rc)

	core.ContextSetLoglevel(context, core.DcLoglevelDebug)
	core.ContextSetLogfunc(context, helpers.Logger, nil)

	// Open the I/O stream.
	fmt.Printf("Opening the I/O stream (%s, %s).\n", helpers.TransportName(transport), *devname)
	rc = helpers.IOStreamOpen(&iostream, context, descriptor, transport, *devname)
	helpers.ExitOnError("Error when opening the I/O stream.", rc)
	defer core.IostreamClose(iostream)

	fmt.Printf("Opening the device (%s %s, %s).\n",
		core.DescriptorGetVendor(descriptor),
		core.DescriptorGetProduct(descriptor),
		*devname)

	rc = core.DeviceOpen(&device, context, descriptor, iostream)
	defer core.DeviceClose(device)

	helpers.ExitOnError("Error opening the device", rc)

	fmt.Println("Registering the event handler.")
	events := core.DcEventWaiting | core.DcEventProgress | core.DcEventDevinfo | core.DcEventClock | core.DcEventVendor
	rc = core.DeviceSetEvents(device, uint32(events), helpers.EventCallback, nil)
	helpers.ExitOnError("Error when registering the event handler", rc)

	// Register the cancellation handler.
	fmt.Println("Registering the cancellation handler.")
	rc = core.DeviceSetCancel(device, cancelCb, nil)
	helpers.ExitOnError("Error when registering the cancellation handler", rc)

	// Initialize the dive data.
	divedata := newDiveData()
	divedata.Device = device
	divedata.Number = 0
	divedata.Dives = new([]Dive)

	rc = core.DeviceForeach(device, diveCb, unsafe.Pointer(divedata.Ref()))
	helpers.ExitOnError("Error when downloading the dives.", rc)

	data, err := json.Marshal(divedata)
	if err != nil {
		log.Fatal(err)
	}
	divedata.Free()

	err = ioutil.WriteFile(*filename, data, 0644)
	if err != nil {
		fmt.Printf("Error when writing output: %s\n", err)
	}
}
