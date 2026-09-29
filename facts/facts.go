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

long Go_ValueToInteger(void *target)
{
    return ValueToInteger(target);
}

double Go_ValueToDouble(void *target)
{
    return ValueToDouble(target);
}

void *Go_GetMFValue(void *target, long index)
{
    return GetMFValue(target, index);
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

*/
import "C"

import (
	"fmt"
	"strings"
	"unsafe"
)



type Fact struct {
	Name string
	Fields map[string]any
}

// This is the Go equivalent of the old PrintAtom function.
func AtomToString(typ C.int, value unsafe.Pointer) string {
	switch typ {
	case C.FLOAT:
		return fmt.Sprintf("%f", C.Go_ValueToDouble(value))

	case C.INTEGER:
		return fmt.Sprintf("%d", C.Go_ValueToInteger(value))

	case C.SYMBOL:
		return C.GoString(C.Go_ValueToString(value))

	case C.STRING:
		s := C.GoString(C.Go_ValueToString(value))
		return `"` + s + `"`

	case C.EXTERNAL_ADDRESS:
		return fmt.Sprintf("<Pointer-%p>", value)

	case C.INSTANCE_NAME:
		return "[" + C.GoString(C.Go_ValueToString(value)) + "]"

	case C.RVOID:
		return ""

	default:
		return fmt.Sprintf("<unknown atom type %d>", typ)
	}
}

// DataObjectToString converts a CLIPS DATA_OBJECT to a Go string.



func DataObjectToString(d *C.DATA_OBJECT) string {
	if d == nil {
		return "<nil DATA_OBJECT>"
	}

	switch d._type {
	case C.INTEGER:
		return fmt.Sprintf("%d", C.Go_ValueToInteger(unsafe.Pointer(d)))

	case C.FLOAT:
		return fmt.Sprintf("%f", C.Go_ValueToDouble(unsafe.Pointer(d)))

	case C.SYMBOL:
		return C.GoString(C.Go_ValueToString(d.value))

	case C.STRING:
		return `"` + C.GoString(C.Go_ValueToString(d.value)) + `"`

	case C.MULTIFIELD:
		length := C.Go_GetMFLength(d.value)
		var b strings.Builder
		b.WriteString("(")

		for i := C.long(1); i <= length; i++ {
			if i > 1 {
				b.WriteByte(' ')
			}
			field := C.Go_GetMFValue(d.value, i)
			b.WriteString(fieldToString(field))
		}

		b.WriteString(")")
		return b.String()

	case C.FACT_ADDRESS:
		return fmt.Sprintf("<Fact-Address %p>", d.value)

	case C.INSTANCE_ADDRESS:
		return fmt.Sprintf("<Instance-Address %p>", d.value)

	case C.EXTERNAL_ADDRESS:
		return fmt.Sprintf("<External-Address %p>", d.value)

	case C.INSTANCE_NAME:
		return C.GoString(C.Go_ValueToString(d.value))

	default:
		return fmt.Sprintf("<Unknown type=%d value=%p>", d._type, d.value)
	}
}

func fieldToString(field unsafe.Pointer) string {
	f := (*C.struct_field)(field)

	switch f._type {
	case C.INTEGER:
		return fmt.Sprintf("%d", C.Go_ValueToInteger(unsafe.Pointer(f)))
	case C.FLOAT:
		return fmt.Sprintf("%f", C.Go_ValueToDouble(unsafe.Pointer(f)))
	case C.SYMBOL:
		return C.GoString(C.Go_ValueToString(f.value))
	case C.STRING:
		return `"` + C.GoString(C.Go_ValueToString(f.value)) + `"`
	default:
		return fmt.Sprintf("<type=%d>", f._type)
	}
}


// PrintFactsToBuffer iterates through facts and returns their names and values as a string
// Returns the string containing facts formatted as:
//   Fact-<index> (<deftemplate-name>):
//     <slot-name>: <slot-value>
//     ...
//   (blank line between facts)
// Returns an empty string on error
func PrintFactsToBuffer() string {
	var output strings.Builder
	var firstFact bool = true

	// Iterate through facts in current module
	var factPtr unsafe.Pointer = C.GetNextFactInScope(nil)
	for factPtr != nil {
		// Add separator between facts
		if !firstFact {
			output.WriteString("\n\n")
		} else {
			firstFact = false
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
		fmt.Fprintf(&output, "F-%d (%s) CF %.2f:\n", int(factIndex), C.GoString(deftemplateName), float32((*C.struct_fact)(factPtr).factCF))

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

				field := (*C.struct_field)(unsafe.Pointer(uintptr(unsafe.Pointer(&multifield.theFields)) + uintptr(i*C.sizeof_struct_field)))

				// Since struct field and DATA_OBJECT have the same layout for type and value,
				// we can cast the field pointer to DATA_OBJECT*
				//valueStr := C.GoString(C.data_object_to_string((*C.DATA_OBJECT)(unsafe.Pointer(field))))
				valueStr := AtomToString(C.int(field._type), field.value) 
				fmt.Fprintf(&output, "  [%d]: %s\n", int(i), valueStr)
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
					valueStr := DataObjectToString(&value)
					fmt.Fprintf(&output, "  %s: %s\n", slotNameGo, valueStr)
				}

				slotPtr = (*C.struct_templateSlot)(unsafe.Pointer(slotPtr)).next
			}
		}
		// Get next fact
		factPtr = C.GetNextFactInScope(factPtr)
	}

	return output.String()
}
