package facts

/*
#cgo LDFLAGS: -lfzclips
#include <stdio.h>

void my_puts(const char* s) {
    printf("%s\n", s);
}
*/
import "C"

// TestCGO tests if CGO is working
func TestCGO() {
	C.my_puts(C.CString("Hello from C"))
}