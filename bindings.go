package jsondiffcpp

// #include "c_wrapper.hxx"
// #include <stdlib.h>
import "C"
import (
	"errors"
	"fmt"
	"sync"
	"unsafe"
)

var ErrGeneratingPatch = errors.New("Error generating json patch")

const BufferSize = 1000 * 1000 * 10
const MaxBuffers = 10

type buffer struct {
	index int
	ptr unsafe.Pointer
}

var buffers [MaxBuffers]buffer
var freeBuffers map[int]struct{}
var bufferLock sync.Mutex
var oncer sync.Once

func getBuffer() *buffer {
	for {
		bufferLock.Lock()
		for k := range freeBuffers {
			delete(freeBuffers, k)
			bufferLock.Unlock()
			return &buffers[k]
		}
		bufferLock.Unlock()
	}
}

func putBuffer(buf *buffer) {
	bufferLock.Lock()
	freeBuffers[buf.index] = struct{}{}
	bufferLock.Unlock()
}

func initialize() {
		freeBuffers = make(map[int]struct{})
		for i:=0; i<MaxBuffers; i++ {
			buffers[i] = buffer{
				index: i,
				ptr: C.malloc(C.sizeof_char * BufferSize),
			}
			freeBuffers[i] = struct{}{}
		}
}

// JSONDiff Creates a json patch
func GeneratePatch(originalJSON, newJSON []byte) (patch string, err error) {
	oncer.Do(initialize)
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Error generating patch: %s", r)
			err = ErrGeneratingPatch
		}
	}()
	cOldJSON := C.CString(string(originalJSON))
	defer C.free(unsafe.Pointer(cOldJSON))

	cNewJSON := C.CString(string(newJSON))
	defer C.free(unsafe.Pointer(cNewJSON))

	buf := getBuffer()
	defer putBuffer(buf)

	length := C.json_patch((*C.char)(cOldJSON), (*C.char)(cNewJSON), (*C.char)(buf.ptr))
	patch = C.GoString((*C.char)(buf.ptr))[:length]
	return
}
