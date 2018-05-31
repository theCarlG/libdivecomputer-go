package helpers

/*
#cgo pkg-config: libdivecomputer
#include <stddef.h>
#include "libdivecomputer/parser.h"
#include "libdivecomputer/buffer.h"
#include "libdivecomputer/device.h"
#include "libdivecomputer/version.h"
#include "libdivecomputer/context.h"
#include "libdivecomputer/descriptor.h"
#include <stdlib.h>
#include "../../../core/cgo_helpers.h"
*/
import "C"

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"unsafe"

	core "github.com/calle-gunnarsson/libdivecomputer-go/core"
)

func DescriptorSearch(out **core.Descriptor, name string, family core.Family, model uint32) core.Status {
	var iterator *core.Iterator

	rc := core.DescriptorIterator(&iterator)
	if rc != core.DcStatusSuccess {
		fmt.Errorf("Error creating the device descriptor iterator.\n")
		return rc
	}

	var (
		descriptor *core.Descriptor
		current    *core.Descriptor
	)
	name = strings.ToLower(name)

	for core.IteratorNext(iterator, (unsafe.Pointer(&descriptor))) == core.DcStatusSuccess {
		if name != "" {
			vendor := strings.ToLower(core.DescriptorGetVendor(descriptor))
			product := strings.ToLower(core.DescriptorGetProduct(descriptor))

			if name == fmt.Sprintf("%s %s", vendor, product) || name == product {
				current = descriptor
				break
			}
		} else {
			if family == core.DescriptorGetType(descriptor) {
				if model == core.DescriptorGetModel(descriptor) {
					core.DescriptorFree(current)
					current = descriptor
					break
				} else {
					if current == nil {
						current = descriptor
						descriptor = nil
					}
				}
			}
		}
		core.DescriptorFree(descriptor)
	}

	if rc != core.DcStatusSuccess && rc != core.DcStatusDone {
		core.DescriptorFree(current)
		core.IteratorFree(iterator)
		fmt.Errorf("Error iterating the device descriptors.\n")
		return rc
	}

	core.IteratorFree(iterator)
	*out = current

	return core.DcStatusSuccess
}

func FamilyModel(t core.Family) uint {
	for _, b := range backends {
		if b.Type == t {
			return b.Model
		}
	}

	return 0
}

func FamilyType(name string) core.Family {
	for _, b := range backends {
		if b.Name == name {
			return b.Type
		}
	}

	return core.DcFamilyNull
}

func FamilyName(t core.Family) string {
	for _, b := range backends {
		if b.Type == t {
			return b.Name
		}
	}

	return ""
}

func EventCallback(device *core.Device, event core.EventType, data unsafe.Pointer, userdata unsafe.Pointer) {
	switch event {
	case core.DcEventWaiting:
		fmt.Printf("Event: waiting for user action\n")
		break
	case core.DcEventProgress:
		progress := core.NewEventProgressRef(unsafe.Pointer(data))
		progress.Deref()
		fmt.Printf("Event: progress %3.2f%% (%d/%d)\n", 100.0*float32(progress.Current)/float32(progress.Maximum), progress.Current, progress.Maximum)
		break
	case core.DcEventDevinfo:
		devinfo := core.NewEventDevinfoRef(unsafe.Pointer(data))
		devinfo.Deref()
		fmt.Printf("Event: model=%d (0x%08x), firmware=%d (0x%08x), serial=%d (0x%08x)\n",
			devinfo.Model, devinfo.Model,
			devinfo.Firmware, devinfo.Firmware,
			devinfo.Serial, devinfo.Serial)
		break
	case core.DcEventClock:
		clock := core.NewEventClockRef(unsafe.Pointer(data))
		clock.Deref()
		fmt.Printf("Event: systime=%d, devtime=%d\n",
			clock.Systime, clock.Devtime)
		break
	case core.DcEventVendor:
		vendor := core.NewEventVendorRef(data)
		vendor.Deref()

		size := vendor.Size
		data := []byte(vendor.Data)
		fmt.Printf("Event: vendor=")
		for i := uint32(0); i < size; i++ {
			fmt.Printf("%02X", data[i])
		}
		fmt.Println()
		break
	default:
		break
	}
}

func Logger(context *core.Context, loglevel core.Loglevel, file string, line uint32, function string, message string, userdata unsafe.Pointer) {
	loglevels := []string{"NONE", "ERROR", "WARNING", "INFO", "DEBUG", "ALL"}

	//fmt.Printf("%s: %s [in %s:%d (%s)]\n", loglevels[loglevel], message, file, line, function)
	fmt.Printf("%s: %s\n", loglevels[loglevel], message)
}

func FileWrite(filename string, buffer *core.Buffer) error {
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	data := core.BufferGetData(buffer)
	size := core.BufferGetSize(buffer)
	b := C.GoBytes(unsafe.Pointer(data), C.int(size))
	err = binary.Write(out, binary.LittleEndian, b)

	return nil
}

func FileRead(filename string) (*core.Buffer, error) {
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	buffer := core.BufferNew(0)
	core.BufferAppend(buffer, string(data), uint(len(data)))

	return buffer, nil
}

func ErrMsg(status core.Status) string {
	switch status {
	case core.DcStatusSuccess:
		return "Success"
	case core.DcStatusUnsupported:
		return "Unsupported operation"
	case core.DcStatusInvalidargs:
		return "Invalid arguments"
	case core.DcStatusNomemory:
		return "Out of memory"
	case core.DcStatusNodevice:
		return "No device found"
	case core.DcStatusNoaccess:
		return "Access denied"
	case core.DcStatusIo:
		return "Input/output error"
	case core.DcStatusTimeout:
		return "Timeout"
	case core.DcStatusProtocol:
		return "Protocol error"
	case core.DcStatusDataformat:
		return "Data format error"
	case core.DcStatusCancelled:
		return "Cancelled"
	default:
		return "Unknown error"
	}
}

func ExitOnError(msg string, rc core.Status) {
	if rc != core.DcStatusSuccess {
		fmt.Printf("%s: %s\n", msg, ErrMsg(rc))
		os.Exit(int(rc))
	}
}
func CheckError(msg string, rc core.Status) bool {
	if rc != core.DcStatusSuccess {
		fmt.Printf("%s: %s\n", msg, ErrMsg(rc))
		return true
	}
	return false
}
