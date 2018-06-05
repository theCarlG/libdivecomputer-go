package helpers

import (
	"fmt"
	"strconv"
	"unsafe"

	"github.com/calle-gunnarsson/libdivecomputer-go/core"
)

func UsbhidOpen(out **core.Iostream, context *core.Context, descriptor *core.Descriptor) core.Status {
	var (
		iostream *core.Iostream     = nil
		iterator *core.Iterator     = nil
		device   *core.UsbhidDevice = nil
	)

	core.UsbhidIteratorNew(&iterator, context, descriptor)
	for core.IteratorNext(iterator, unsafe.Pointer(&device)) == core.DcStatusSuccess {
		break
	}
	core.IteratorFree(iterator)
	defer core.UsbhidDeviceFree(device)

	if device == nil {
		fmt.Println("No dive computer found")
		return core.DcStatusNodevice
	}

	rc := core.UsbhidOpen(&iostream, context, device)
	if CheckError("Failed to open usbhid device.", rc) {
		return rc
	}

	*out = iostream

	return core.DcStatusSuccess
}

func IrdaOpen(out **core.Iostream, context *core.Context, descriptor *core.Descriptor, devname string) core.Status {
	var (
		iostream *core.Iostream = nil
		address  uint32         = 0
	)

	if devname != "" {
		address, err := strconv.ParseUint(devname, 10, 32)
		if err != nil {
			fmt.Printf("Error when converting address %s: %s\n", address, err)
			return core.DcStatusInvalidargs
		}
	} else {
		// Discover the device address.
		var (
			iterator *core.Iterator   = nil
			device   *core.IrdaDevice = nil
		)
		core.IrdaIteratorNew(&iterator, context, descriptor)
		for core.IteratorNext(iterator, unsafe.Pointer(&device)) == core.DcStatusSuccess {
			address = core.IrdaDeviceGetAddress(device)
			core.IrdaDeviceFree(device)
			break
		}
		core.IteratorFree(iterator)
	}

	if address == 0 {
		if devname != "" {
			fmt.Println("No valid device address specified.")
		} else {
			fmt.Println("No dive computer found.")
		}
		return core.DcStatusNodevice
	}

	// Open the irda socket.
	rc := core.IrdaOpen(&iostream, context, address, 1)
	if CheckError("Failed to open irda socket", rc) {
		return rc
	}

	*out = iostream

	return core.DcStatusSuccess
}

func BluetoothOpen(out **core.Iostream, context *core.Context, descriptor *core.Descriptor, devname string) core.Status {
	var (
		iostream *core.Iostream        = nil
		address  core.BluetoothAddress = 0
	)

	if devname != "" {
		// Use the address.
		address = core.BluetoothStr2addr(devname)
	} else {
		var (
			iterator *core.Iterator        = nil
			device   *core.BluetoothDevice = nil
		)
		core.BluetoothIteratorNew(&iterator, context, descriptor)
		for core.IteratorNext(iterator, unsafe.Pointer(&device)) == core.DcStatusSuccess {
			address = core.BluetoothDeviceGetAddress(device)
			break
		}
		core.IteratorFree(iterator)
	}

	if address == 0 {
		if devname != "" {
			fmt.Println("No valid device address specified.")
		} else {
			fmt.Println("No dive computer found.")
		}
		return core.DcStatusNodevice
	}

	// Open the bluetooth socket.
	rc := core.BluetoothOpen(&iostream, context, address, 0)
	if CheckError("Failed to open bluetooth socket", rc) {
		return rc
	}

	*out = iostream

	return core.DcStatusSuccess
}

func IOStreamOpen(iostream **core.Iostream, context *core.Context, descriptor *core.Descriptor, transport core.Transport, devname string) core.Status {
	switch transport {
	case core.DcTransportSerial:
		return core.SerialOpen(iostream, context, devname)
	case core.DcTransportUsb:
		return core.DcStatusSuccess
	case core.DcTransportUsbhid:
		return UsbhidOpen(iostream, context, descriptor)
	case core.DcTransportIrda:
		return IrdaOpen(iostream, context, descriptor, devname)
	case core.DcTransportBluetooth:
		return BluetoothOpen(iostream, context, descriptor, devname)
	default:
		return core.DcStatusUnsupported
	}
}
