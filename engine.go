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
extern void jsObjectInitialize(JSContextRef ctx, JSObjectRef obj);
extern void jsObjectFinalize(JSObjectRef obj);
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

// Context is a JS context.
type Context interface {
	jsContext() C.JSContextRef
	// GlobalObject returns the global object associated with the context.
	GlobalObject() *Object
}

// GlobalContext is a global JS context.
type GlobalContext struct {
	ctx C.JSGlobalContextRef
}

// NewGlobalContext creates a new global JS context.
func NewGlobalContext() *GlobalContext {
	return &GlobalContext{C.JSGlobalContextCreate(nil)}
}

func (c *GlobalContext) jsContext() C.JSContextRef {
	return C.JSContextRef(c.ctx)
}

// GlobalObject returns the global object associated with the context.
func (c *GlobalContext) GlobalObject() *Object {
	return &Object{C.JSContextGetGlobalObject(c.ctx)}
}

// Release releases the context.
func (c *GlobalContext) Release() {
	C.JSGlobalContextRelease(c.ctx)
}

// CheckScriptSyntax checks the syntax of the script.
func CheckScriptSyntax(ctx Context, script string) bool {
	s := createJSString(script)
	defer releaseJSString(s)
	return bool(C.JSCheckScriptSyntax(ctx.jsContext(), s, nil, 1, nil))
}

// EvaluateScript evaluates the provided script.
func EvaluateScript(ctx Context, script string) *Value {
	s := createJSString(script)
	defer releaseJSString(s)
	r := C.JSEvaluateScript(ctx.jsContext(), s, nil, nil, 1, nil)
	return &Value{r}
}

// Type is a JS type.
type Type int

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

const (
	Undefined Type = iota
	Null
	Boolean
	Number
	String
	JSObject
	Symbol
)

//export jsObjectInitialize
func jsObjectInitialize(ctx C.JSContextRef, obj C.JSObjectRef) {
	// h := cgo.Handle(C.JSObjectGetPrivate(obj))
	// fmt.Println("init:", h.Value())
}

//export jsObjectFinalize
func jsObjectFinalize(obj C.JSObjectRef) {
	h := cgo.Handle(C.JSObjectGetPrivate(obj))
	// fmt.Println("fin:", h.Value())
	h.Delete()
}

var (
	emptyJSClass C.JSClassRef
	goJSClass    C.JSClassRef
)

func init() {
	// empty class
	emptyClassDef := C.kJSClassDefinitionEmpty
	emptyJSClass = C.JSClassCreate(&emptyClassDef)
	// Go class
	goClassDef := C.kJSClassDefinitionEmpty
	goClassDef.initialize = (C.JSObjectInitializeCallback)(C.jsObjectInitialize)
	goClassDef.finalize = (C.JSObjectFinalizeCallback)(C.jsObjectFinalize)
	goJSClass = C.JSClassCreate(&goClassDef)
}

// Object is a JS object.
type Object struct {
	obj C.JSObjectRef
}

// NewObject creates a new JS object.
func NewObject(ctx Context) *Object {
	return &Object{C.JSObjectMake(ctx.jsContext(), emptyJSClass, nil)}
}

// NewGoObject creates a new JS object that wraps a Go object.
func NewGoObject(ctx Context, data interface{}) *Object {
	return &Object{C.JSObjectMake(ctx.jsContext(), goJSClass, unsafe.Pointer(cgo.NewHandle(data)))}
}

// NewArray creates a JS array.
func NewArray(ctx Context, values ...*Value) *Object {
	vals := make([]C.JSValueRef, len(values))
	for i, v := range values {
		vals[i] = v.val
	}
	ptr := (*C.JSValueRef)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&vals)).Data))
	return &Object{C.JSObjectMakeArray(ctx.jsContext(), C.ulong(len(vals)), ptr, nil)}
}

// Value returns the object as a JS value.
func (o *Object) Value() *Value {
	return &Value{(C.JSValueRef)(unsafe.Pointer(o.obj))}
}

// Has checks whether the object has a property with the given name.
func (o *Object) Has(ctx Context, propertyName string) bool {
	n := createJSString(propertyName)
	defer releaseJSString(n)
	return bool(C.JSObjectHasProperty(ctx.jsContext(), o.obj, n))
}

// Get returns the property with the given name.
func (o *Object) Get(ctx Context, propertyName string) *Value {
	n := createJSString(propertyName)
	defer releaseJSString(n)
	return &Value{C.JSObjectGetProperty(ctx.jsContext(), o.obj, n, nil)}
}

// Set sets the property with the given name.
func (o *Object) Set(ctx Context, propertyName string, value *Value) {
	n := createJSString(propertyName)
	defer releaseJSString(n)
	C.JSObjectSetProperty(ctx.jsContext(), o.obj, n, value.val, 0, nil)
}

// Value is a JS value.
type Value struct {
	val C.JSValueRef
}

// NewNull creates a new null JS value.
func NewNull(ctx Context) *Value {
	return &Value{C.JSValueMakeNull(ctx.jsContext())}
}

// NewNumber creates a new JS number.
func NewNumber(ctx Context, x float64) *Value {
	return &Value{C.JSValueMakeNumber(ctx.jsContext(), C.double(x))}
}

// NewBoolean creates a new JS boolean.
func NewBoolean(ctx Context, b bool) *Value {
	return &Value{C.JSValueMakeBoolean(ctx.jsContext(), C.bool(b))}
}

// NewString creates a new JS string.
func NewString(ctx Context, str string) *Value {
	s := createJSString(str)
	defer releaseJSString(s)
	return &Value{C.JSValueMakeString(ctx.jsContext(), s)}
}

// NewSymbol creates a new JS symbol.
func NewSymbol(ctx Context, str string) *Value {
	s := createJSString(str)
	defer releaseJSString(s)
	return &Value{C.JSValueMakeSymbol(ctx.jsContext(), s)}
}

// Protect protects the JS value from being garbage collected.
func (v *Value) Protect(ctx Context) {
	C.JSValueProtect(ctx.jsContext(), v.val)
}

// Unprotect unprotects the JS value from being garbage collected.
func (v *Value) Unprotect(ctx Context) {
	C.JSValueUnprotect(ctx.jsContext(), v.val)
}

// Type returns the type of the JS value.
func (v *Value) Type(ctx Context) Type {
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
func (v *Value) Number(ctx Context) float64 {
	return float64(C.JSValueToNumber(ctx.jsContext(), v.val, nil))
}

// Boolean returns the JS value as a boolean.
func (v *Value) Boolean(ctx Context) bool {
	return bool(C.JSValueToBoolean(ctx.jsContext(), v.val))
}

// Object returns the JS value as an object.
func (v *Value) Object(ctx Context) *Object {
	return &Object{C.JSValueToObject(ctx.jsContext(), v.val, nil)}
}

// String returns the textual representation of the JS value.
func (v *Value) String(ctx Context) string {
	s := C.JSValueToStringCopy(ctx.jsContext(), v.val, nil)
	defer C.JSStringRelease(s)
	maxLen := C.JSStringGetMaximumUTF8CStringSize(s)
	cptr := (*C.char)(C.malloc(maxLen))
	defer C.free(unsafe.Pointer(cptr))
	C.JSStringGetUTF8CString(s, cptr, maxLen)
	return C.GoString(cptr)
}

// Interface returns the JS value as an interface{}.
// If the value is a wrapped Go object, then the wrapper object is returned.
// For other JS objects, a map of type map[string]interface{} is returned.
func (v *Value) Interface(ctx Context) interface{} {
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
				v := &Value{C.JSObjectGetProperty(ctx.jsContext(), o, s, nil)}
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
		return fmt.Sprint("unknown JS type: ", int(t))
	}
}
