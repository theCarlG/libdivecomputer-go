package main

import (
	"fmt"
	"unsafe"

	core "github.com/calle-gunnarsson/libdivecomputer-go/core"
)

func main() {
	var iterator *core.Iterator
	var descriptor *core.Descriptor

	rc := core.DescriptorIterator(&iterator)
	if rc != core.DcStatusSuccess {
		fmt.Errorf("Error creating the device descriptor iterator.\n")
		return
	}

	for core.IteratorNext(iterator, (unsafe.Pointer(&descriptor))) == core.DcStatusSuccess {
		fmt.Printf("%s %s\n",
			core.DescriptorGetVendor(descriptor),
			core.DescriptorGetProduct(descriptor))
		core.DescriptorFree(descriptor)
	}

	core.IteratorFree(iterator)
}
