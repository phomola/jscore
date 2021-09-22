package jscore

import (
	"fmt"
	"reflect"
	"runtime/cgo"
	"unsafe"
)

/*
#cgo darwin LDFLAGS: -framework JavaScriptCore
#include <JavaScriptCore/JavaScriptCore.h>
*/
import "C"

// Value is a JS value.
type Value struct {
	val C.JSValueRef
}

// NewValue creates a new JS value.
func NewValue(ctx Context, v interface{}) Value {
	if v == nil {
		return NewNull(ctx)
	}
	switch v := v.(type) {
	case int:
		return NewNumber(ctx, float64(v))
	case float32:
		return NewNumber(ctx, float64(v))
	case float64:
		return NewNumber(ctx, v)
	case bool:
		return NewBoolean(ctx, v)
	case string:
		return NewString(ctx, v)
	}
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Slice {
		s := make([]Value, val.Len())
		for i := 0; i < val.Len(); i++ {
			s[i] = NewValue(ctx, val.Index(i).Interface())
		}
		return NewArray(ctx, s...).Value()
	}
	if val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Struct {
		return NewGoObject(ctx, v).Value()
	}
	panic(fmt.Sprint("can't create JS value from ", v))
}

// NewNull creates a new null JS value.
func NewNull(ctx Context) Value {
	return Value{C.JSValueMakeNull(ctx.jsContext())}
}

// NewNumber creates a new JS number.
func NewNumber(ctx Context, x float64) Value {
	return Value{C.JSValueMakeNumber(ctx.jsContext(), C.double(x))}
}

// NewBoolean creates a new JS boolean.
func NewBoolean(ctx Context, b bool) Value {
	return Value{C.JSValueMakeBoolean(ctx.jsContext(), C.bool(b))}
}

// NewString creates a new JS string.
func NewString(ctx Context, str string) Value {
	s := createJSString(str)
	defer releaseJSString(s)
	return Value{C.JSValueMakeString(ctx.jsContext(), s)}
}

// NewSymbol creates a new JS symbol.
func NewSymbol(ctx Context, str string) Value {
	s := createJSString(str)
	defer releaseJSString(s)
	return Value{C.JSValueMakeSymbol(ctx.jsContext(), s)}
}

// Protect protects the JS value from being garbage collected.
func (v Value) Protect(ctx Context) {
	C.JSValueProtect(ctx.jsContext(), v.val)
}

// Unprotect unprotects the JS value from being garbage collected.
func (v Value) Unprotect(ctx Context) {
	C.JSValueUnprotect(ctx.jsContext(), v.val)
}

// Type returns the type of the JS value.
func (v Value) Type(ctx Context) Type {
	t := C.JSValueGetType(ctx.jsContext(), v.val)
	switch t {
	case C.kJSTypeUndefined:
		return Undefined
	case C.kJSTypeNull:
		return Null
	case C.kJSTypeBoolean:
		return Boolean
	case C.kJSTypeNumber:
		return Number
	case C.kJSTypeString:
		return String
	case C.kJSTypeObject:
		return JSObject
	case C.kJSTypeSymbol:
		return Symbol
	default:
		panic(fmt.Sprint("unknown JS type: ", t))
	}
}

// Number returns the JS value as a number.
func (v Value) Number(ctx Context) float64 {
	return float64(C.JSValueToNumber(ctx.jsContext(), v.val, nil))
}

// Boolean returns the JS value as a boolean.
func (v Value) Boolean(ctx Context) bool {
	return bool(C.JSValueToBoolean(ctx.jsContext(), v.val))
}

// Object returns the JS value as an object.
func (v Value) Object(ctx Context) Object {
	return Object{C.JSValueToObject(ctx.jsContext(), v.val, nil)}
}

// String returns the textual representation of the JS value.
func (v Value) String(ctx Context) string {
	s := C.JSValueToStringCopy(ctx.jsContext(), v.val, nil)
	defer C.JSStringRelease(s)
	maxLen := C.JSStringGetMaximumUTF8CStringSize(s)
	cptr := (*C.char)(C.malloc(maxLen))
	defer C.free(unsafe.Pointer(cptr))
	C.JSStringGetUTF8CString(s, cptr, maxLen)
	return C.GoString(cptr)
}

// Interface returns the JS value as an interface{}.
// If the value is a wrapped Go object, then the wrapped object is returned.
// For other JS objects, a map of type map[string]interface{} is returned.
func (v Value) Interface(ctx Context) interface{} {
	t := v.Type(ctx)
	switch t {
	case Undefined:
		return nil
	case Null:
		return nil
	case Boolean:
		return v.Boolean(ctx)
	case Number:
		return v.Number(ctx)
	case String:
		return v.String(ctx)
	case JSObject:
		o := C.JSValueToObject(ctx.jsContext(), v.val, nil)
		p := C.JSObjectGetPrivate(o)
		if p != nil {
			return cgo.Handle(p).Value()
		} else {
			names := C.JSObjectCopyPropertyNames(ctx.jsContext(), o)
			defer C.JSPropertyNameArrayRelease(names)
			l := int(C.JSPropertyNameArrayGetCount(names))
			m := make(map[string]interface{}, l)
			for i := 0; i < l; i++ {
				s := C.JSPropertyNameArrayGetNameAtIndex(names, C.ulong(i))
				v := Value{C.JSObjectGetProperty(ctx.jsContext(), o, s, nil)}
				maxLen := C.JSStringGetMaximumUTF8CStringSize(s)
				cptr := (*C.char)(C.malloc(maxLen))
				defer C.free(unsafe.Pointer(cptr))
				C.JSStringGetUTF8CString(s, cptr, maxLen)
				name := C.GoString(cptr)
				m[name] = v.Interface(ctx)
			}
			return m
		}
	case Symbol:
		panic("JS value is a symbol")
	default:
		panic(fmt.Sprint("unknown JS type: ", int(t)))
	}
}
