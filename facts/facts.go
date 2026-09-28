package facts

/*
#cgo CFLAGS: -I/home/kapitan/Downloads/gits/FuzzyCLIPS/source
#cgo LDFLAGS: -L/home/kapitan/Downloads/gits/FuzzyCLIPS/source -lfzclips -lm

#include <clips.h>
#include <tmpltdef.h>
#include <tmpltbin.h>
#include <factmngr.h>
#include <tmpltdef.h>
#include <symbol.h>
#include <multifld.h>
#include <router.h>
#include <stdlib.h>
#include <string.h>

long Go_GetMFLength(void *target)
{
    return GetMFLength(target);
}


int Go_DeftemplateImplied(void *theDeftemplate)
{
    return ((struct deftemplate *)theDeftemplate)->implied;
}

int Go_DeftemplateWatch(void *theDeftemplate)
{
    return ((struct deftemplate *)theDeftemplate)->watch;
}

int Go_DeftemplateInScope(void *theDeftemplate)
{
    return ((struct deftemplate *)theDeftemplate)->inScope;
}

int Go_DeftemplateNumberOfSlots(void *theDeftemplate)
{
    return ((struct deftemplate *)theDeftemplate)->numberOfSlots;
}

const char *Go_ValueToString(void *target)
{
    return ValueToString(target);
}

const char *Go_GetDeftemplateName(void *theDeftemplate)
{
    return GetDeftemplateName(theDeftemplate);
}

// Helper function to convert DATA_OBJECT to string
static char* data_object_to_string(DATA_OBJECT* theValue) {
    static char buffer[1024];
    buffer[0] = '\0';

    if (theValue == NULL) {
        return "<null>";
    }

    switch (theValue->type) {
        case INTEGER:
            snprintf(buffer, sizeof(buffer), "%ld", ValueToLong(theValue->value));
            break;
        case FLOAT:
            snprintf(buffer, sizeof(buffer), "%g", ValueToDouble(theValue->value));
            break;
        case SYMBOL:
            strncpy(buffer, ValueToString(theValue->value), sizeof(buffer)-1);
            buffer[sizeof(buffer)-1] = '\0';
            break;
        case STRING:
            strncpy(buffer, ValueToString(theValue->value), sizeof(buffer)-1);
            buffer[sizeof(buffer)-1] = '\0';
            break;
        case MULTIFIELD: {
            struct multifield* mf = (struct multifield*)theValue->value;
            long length = GetMFLength(mf);
            int pos = 0;
            pos += snprintf(buffer+pos, sizeof(buffer)-pos, "[");
            for (long i = 1; i <= length; i++) {
                struct field* fieldPtr = GetMFPtr(mf, i);
                char temp[64];
                temp[0] = '\0';
                switch (fieldPtr->type) {
                    case INTEGER:
                        snprintf(temp, sizeof(temp), "%ld", ValueToLong(fieldPtr->value));
                        break;
                    case FLOAT:
                        snprintf(temp, sizeof(temp), "%g", ValueToDouble(fieldPtr->value));
                        break;
                    case SYMBOL:
                        strncpy(temp, ValueToString(fieldPtr->value), sizeof(temp)-1);
                        temp[sizeof(temp)-1] = '\0';
                        break;
                    case STRING:
                        strncpy(temp, ValueToString(fieldPtr->value), sizeof(temp)-1);
                        temp[sizeof(temp)-1] = '\0';
                        break;
                    default:
                        strncpy(temp, "<unknown>", sizeof(temp)-1);
                        temp[sizeof(temp)-1] = '\0';
                        break;
                }
                if (i > 1) {
                    pos += snprintf(buffer+pos, sizeof(buffer)-pos, ", ");
                }
                pos += snprintf(buffer+pos, sizeof(buffer)-pos, "%s", temp);
            }
            pos += snprintf(buffer+pos, sizeof(buffer)-pos, "]");
            break;
        }
        case FACT_ADDRESS: {
            struct fact* f = (struct fact*)theValue->value;
            snprintf(buffer, sizeof(buffer), "<Fact-%ld>", f->factIndex);
            break;
        }
        case EXTERNAL_ADDRESS:
            strncpy(buffer, "<external>", sizeof(buffer)-1);
            buffer[sizeof(buffer)-1] = '\0';
            break;
        default:
            snprintf(buffer, sizeof(buffer), "<type:%d>", theValue->type);
            break;
    }
    return buffer;
}
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// PrintFactsToBuffer iterates through facts and prints their names and values to the provided buffer
// Returns the number of bytes written (including null terminator) or -1 on error
//
// Parameters:
//   - buffer: pointer to a character array where the output will be stored
//   - bufferLength: maximum number of bytes to write (including null terminator)
//
// The output format for each fact:
//   Fact-<index> (<deftemplate-name>):
//     <slot-name>: <slot-value>
//     ...
//   (blank line between facts)
//
// For ordered (implied) facts, slots are shown as indices:
//   Fact-<index> (<deftemplate-name>):
//     [0]: <value>
//     [1]: <value>
//     ...
//
func PrintFactsToBuffer(buffer *C.char, bufferLength C.int) C.int {
	if buffer == nil || bufferLength <= 0 {
		return -1
	}

	// Initialize buffer
	C.memset(unsafe.Pointer(buffer), 0, C.size_t(bufferLength))

	var pos C.int = 0
	var firstFact C.int = 1

	// Iterate through facts in current module
	var factPtr unsafe.Pointer = C.GetNextFactInScope(nil)
	for factPtr != nil {
		// Add separator between facts
		if firstFact != 0 {
			if pos+2 >= bufferLength {
				break // Buffer full
			}
			C.memset(unsafe.Pointer(uintptr(unsafe.Pointer(buffer))+uintptr(pos)), '\n', 1)
			pos++
			C.memset(unsafe.Pointer(uintptr(unsafe.Pointer(buffer))+uintptr(pos)), '\n', 1)
			pos++
		} else {
			firstFact = 0
		}

		// Get fact index (ID)
		factIndex := C.FactIndex(factPtr)

		// Get deftemplate
		deftemplate := ((*C.struct_fact)(factPtr)).whichDeftemplate

		// Get deftemplate name
		var deftemplateName *C.char
		if deftemplate != nil {
			deftemplateName = C.Go_GetDeftemplateName(unsafe.Pointer(deftemplate))
		} else {
			deftemplateName = C.CString("<unknown>")
			defer C.free(unsafe.Pointer(deftemplateName))
		}

		// Write fact header
		header := C.CString(fmt.Sprintf("Fact-%ld (%s):\n", factIndex, C.GoString(deftemplateName)))
		headerLen := C.int(len(C.GoString(header)))
		if pos+headerLen >= bufferLength {
			C.free(unsafe.Pointer(header))
			break
		}
		C.memcpy(unsafe.Pointer(uintptr(unsafe.Pointer(buffer))+uintptr(pos)), unsafe.Pointer(header), C.size_t(headerLen))
		pos += headerLen
		C.free(unsafe.Pointer(header))

		// Check if implied (ordered) fact
		implied := C.int(C.Go_DeftemplateImplied(unsafe.Pointer(deftemplate)))
		if implied != 0 {
			// Ordered fact - single multifield slot
			var value C.DATA_OBJECT

			value._type = C.MULTIFIELD;
			value.value = (*C.struct_fact)(factPtr).theProposition.theFields[0].value;

			// Get the multifield
			multifield := (*C.struct_multifield)(value.value)
			length := C.Go_GetMFLength(unsafe.Pointer(multifield))

			for i := C.long(0); i < length; i++ {
				
				field := (unsafe.Pointer(uintptr(unsafe.Pointer(&multifield.theFields)) + uintptr(i*C.sizeof_struct_field)))

				// Since struct field and DATA_OBJECT have the same layout for type and value,
				// we can cast the field pointer to DATA_OBJECT*
				valueStr := C.GoString(C.data_object_to_string((*C.DATA_OBJECT)(unsafe.Pointer(field))))

				fieldLine := C.CString(fmt.Sprintf("  [%d]: %s\n", i, valueStr))
				lineLen := C.int(len(C.GoString(fieldLine)))
				if pos+lineLen >= bufferLength {
					C.free(unsafe.Pointer(fieldLine))
					break
				}
				C.memcpy(unsafe.Pointer(uintptr(unsafe.Pointer(buffer))+uintptr(pos)), unsafe.Pointer(fieldLine), C.size_t(lineLen))
				pos += lineLen
				C.free(unsafe.Pointer(fieldLine))
			}
		} else {
			// Named fact - iterate through slots
			slotPtr := (*C.struct_deftemplate)(unsafe.Pointer(deftemplate)).slotList
			for slotPtr != nil {
				slotName := C.Go_ValueToString(unsafe.Pointer(slotPtr.slotName))
				slotNameGo := C.GoString(slotName)

				// Get slot value
				var value C.DATA_OBJECT
				success := C.GetFactSlot(factPtr, slotName, &value)
				if success != 0 && value._type != C.RVOID {
					valueStr := C.GoString(C.data_object_to_string(&value))
					line := C.CString(fmt.Sprintf("  %s: %s\n", slotNameGo, valueStr))
					lineLen := C.int(len(C.GoString(line)))
					if pos+lineLen >= bufferLength {
						C.free(unsafe.Pointer(line))
						break
					}
					C.memcpy(unsafe.Pointer(uintptr(unsafe.Pointer(buffer))+uintptr(pos)), unsafe.Pointer(line), C.size_t(lineLen))
					pos += lineLen
					C.free(unsafe.Pointer(line))
				}

				slotPtr = (*C.struct_templateSlot)(unsafe.Pointer(slotPtr)).next
			}
		}
		// Get next fact
		factPtr = C.GetNextFactInScope(factPtr)
	}

	// Ensure null termination
	if pos >= bufferLength {
		pos = bufferLength - 1
	}

	C.memset(unsafe.Pointer(uintptr(unsafe.Pointer(buffer))+uintptr(pos)), 0, 1)
	return pos + 1 // Include null terminator
}
