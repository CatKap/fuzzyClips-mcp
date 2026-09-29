package facts

/*
#cgo LDFLAGS: -lfzclips
#include "factmngr.h"
*/
import "C"

// TestStructAccess tests if we can access C struct fields
func TestStructAccess() {
	var f C.struct_fact
	// Try accessing a field
	_ = f.factIndex
}