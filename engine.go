package jscore

import (
	"unsafe"
)

/*
#cgo darwin LDFLAGS: -framework JavaScriptCore
#include <JavaScriptCore/JavaScriptCore.h>
*/
import "C"

func createJSString(str string) C.JSStringRef {
	s := C.CString(str)
	defer C.free(unsafe.Pointer(s))
	return C.JSStringCreateWithUTF8CString(s)
}

func releaseJSString(s C.JSStringRef) {
	C.JSStringRelease(s)
}

// CheckScriptSyntax checks the syntax of the script.
func CheckScriptSyntax(ctx Context, script string) bool {
	s := createJSString(script)
	defer releaseJSString(s)
	return bool(C.JSCheckScriptSyntax(ctx.jsContext(), s, nil, 1, nil))
}

// EvaluateScript evaluates the provided script.
func EvaluateScript(ctx Context, script string) Value {
	s := createJSString(script)
	defer releaseJSString(s)
	r := C.JSEvaluateScript(ctx.jsContext(), s, nil, nil, 1, nil)
	return Value{r}
}
