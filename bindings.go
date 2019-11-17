package jsondiffcpp

// #include "c_wrapper.hxx"
// #include <stdlib.h>
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

var ErrGeneratingPatch = errors.New("Error generating json patch")

const BUFFER_SIZE = 1000 * 1000 * 10

// JSONDiff Creates a json patch
func GeneratePatch(originalJSON, newJSON string) (patch string, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Error generating patch: %s", r)
			err = ErrGeneratingPatch
		}
	}()
	cOldJSON := C.CString(originalJSON)
	cNewJSON := C.CString(newJSON)
	cPatch := C.malloc(C.sizeof_char * BUFFER_SIZE)
	length := C.json_patch(cOldJSON, cNewJSON, (*C.char)(cPatch))
	patch = C.GoString((*C.char)(cPatch))[:length]
	C.free(unsafe.Pointer(cOldJSON))
	C.free(unsafe.Pointer(cNewJSON))
	C.free(unsafe.Pointer(cPatch))
	return
}
