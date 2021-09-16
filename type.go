package jscore

import (
	"fmt"
)

/*
#cgo darwin LDFLAGS: -framework JavaScriptCore
#include <JavaScriptCore/JavaScriptCore.h>
*/
import "C"

// Type is a JS type.
type Type int

const (
	Undefined Type = iota
	Null
	Boolean
	Number
	String
	JSObject
	Symbol
)

// String returns the textual representation of the type.
func (t Type) String() string {
	switch t {
	case Undefined:
		return "undefined"
	case Null:
		return "null"
	case Boolean:
		return "boolean"
	case Number:
		return "number"
	case String:
		return "string"
	case JSObject:
		return "object"
	case Symbol:
		return "symbol"
	default:
		return fmt.Sprint("unknown JS type: ", int(t))
	}
}
