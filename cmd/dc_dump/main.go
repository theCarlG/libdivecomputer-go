package main

import (
	"flag"
	"fmt"
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
	devname := flag.String("devname", "", "Tty device name")

	flag.Parse()

	if *filename == "" {
		fmt.Println("You need to specify filename")
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

	fmt.Printf("Opening the device (%s %s, %s).\n",
		core.DescriptorGetVendor(descriptor),
		core.DescriptorGetProduct(descriptor),
		*devname)

	rc = core.DeviceOpen(&device, context, descriptor, *devname)
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

	buffer := core.BufferNew(0)
	defer core.BufferFree(buffer)
	// Download the memory dump.
	fmt.Println("Downloading the memory dump.")
	rc = core.DeviceDump(device, buffer)
	helpers.ExitOnError("Error when downloading the memory dump", rc)

	err := helpers.FileWrite(*filename, buffer)
	if err != nil {
		fmt.Printf("Error when writing buffer: %s\n", err)
	}
}
