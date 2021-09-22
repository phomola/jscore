package jscore

import (
	"reflect"
	"runtime/cgo"
	"unsafe"
)

/*
#cgo darwin LDFLAGS: -framework JavaScriptCore
#include <JavaScriptCore/JavaScriptCore.h>
*/
import "C"

// Object is a JS object.
type Object struct {
	obj C.JSObjectRef
}

// NewObject creates a new JS object.
func NewObject(ctx Context) Object {
	return Object{C.JSObjectMake(ctx.jsContext(), emptyJSClass, nil)}
}

// NewGoObject creates a new JS object that wraps a Go object.
func NewGoObject(ctx Context, data interface{}) Object {
	t := reflect.TypeOf(data)
	if t.Kind() != reflect.Ptr {
		goto panic
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		goto panic
	}
	return Object{C.JSObjectMake(ctx.jsContext(), jsClassForType(t), unsafe.Pointer(cgo.NewHandle(data)))}
panic:
	panic("Go object wrapped by a JS object must be a pointer to a struct")
}

// NewArray creates a JS array.
func NewArray(ctx Context, values ...Value) Object {
	vals := make([]C.JSValueRef, len(values))
	for i, v := range values {
		vals[i] = v.val
	}
	ptr := (*C.JSValueRef)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&vals)).Data))
	return Object{C.JSObjectMakeArray(ctx.jsContext(), C.ulong(len(vals)), ptr, nil)}
}

// Value returns the object as a JS value.
func (o Object) Value() Value {
	return Value{(C.JSValueRef)(unsafe.Pointer(o.obj))}
}

// Has checks whether the object has a property with the given name.
func (o Object) Has(ctx Context, propertyName string) bool {
	n := createJSString(propertyName)
	defer releaseJSString(n)
	return bool(C.JSObjectHasProperty(ctx.jsContext(), o.obj, n))
}

// Get returns the property with the given name.
func (o Object) Get(ctx Context, propertyName string) Value {
	n := createJSString(propertyName)
	defer releaseJSString(n)
	return Value{C.JSObjectGetProperty(ctx.jsContext(), o.obj, n, nil)}
}

// At returns the property at the given index.
func (o Object) At(ctx Context, index int) Value {
	return Value{C.JSObjectGetPropertyAtIndex(ctx.jsContext(), o.obj, C.uint(index), nil)}
}

// Set sets the property with the given name.
func (o Object) Set(ctx Context, propertyName string, value Value) {
	n := createJSString(propertyName)
	defer releaseJSString(n)
	C.JSObjectSetProperty(ctx.jsContext(), o.obj, n, value.val, 0, nil)
}

// Array returns the JS object as a slice.
func (o Object) Slice(ctx Context) []interface{} {
	var (
		s []interface{}
		i int
	)
	for {
		v := o.At(ctx, i)
		if v.Type(ctx) == Undefined {
			break
		}
		s = append(s, v.Interface(ctx))
		i++
	}
	return s
}
