package jscore

import (
	"reflect"
	"runtime/cgo"
	"strings"
	"sync"
	"unsafe"
)

/*
#cgo darwin LDFLAGS: -framework JavaScriptCore
#include <JavaScriptCore/JavaScriptCore.h>
extern void jsObjectInitialize(JSContextRef ctx, JSObjectRef obj);
extern void jsObjectFinalize(JSObjectRef obj);
extern JSValueRef jsObjectGetProperty(JSContextRef ctx, JSObjectRef obj, JSStringRef name, JSValueRef* exc);
*/
import "C"

var (
	emptyJSClass C.JSClassRef
	jsClasses    sync.Map
)

type classInfo struct {
	jsClass      C.JSClassRef
	fieldIndices map[string][]int
}

func init() {
	classDef := C.kJSClassDefinitionEmpty
	emptyJSClass = C.JSClassCreate(&classDef)
}

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

//export jsObjectGetProperty
func jsObjectGetProperty(ctx C.JSContextRef, obj C.JSObjectRef, name C.JSStringRef, exc *C.JSValueRef) C.JSValueRef {
	o := cgo.Handle(C.JSObjectGetPrivate(obj)).Value()
	v, _ := jsClasses.Load(reflect.TypeOf(o).Elem())
	ci := v.(*classInfo)
	maxLen := C.JSStringGetMaximumUTF8CStringSize(name)
	cptr := (*C.char)(C.malloc(maxLen))
	defer C.free(unsafe.Pointer(cptr))
	C.JSStringGetUTF8CString(name, cptr, maxLen)
	n := C.GoString(cptr)
	if index, ok := ci.fieldIndices[n]; ok {
		return NewValue(contextWrapper{ctx}, reflect.ValueOf(o).Elem().FieldByIndex(index).Interface()).val
	} else {
		return nil
	}
}

func jsClassForType(t reflect.Type) C.JSClassRef {
	if c, ok := jsClasses.Load(t); ok {
		return c.(*classInfo).jsClass
	}
	classDef := C.kJSClassDefinitionEmpty
	classDef.initialize = (C.JSObjectInitializeCallback)(C.jsObjectInitialize)
	classDef.finalize = (C.JSObjectFinalizeCallback)(C.jsObjectFinalize)
	classDef.getProperty = (C.JSObjectGetPropertyCallback)(C.jsObjectGetProperty)
	c := C.JSClassCreate(&classDef)
	fi := make(map[string][]int, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("jscore")
		if tag == "-" {
			continue
		}
		n := f.Name
		if tag != "" {
			comps := strings.Split(tag, ",")
			if comps[0] != "" {
				n = comps[0]
			}
		}
		fi[n] = f.Index
	}
	jsClasses.Store(t, &classInfo{c, fi})
	return c
}
