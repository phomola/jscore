package jscore

/*
#cgo darwin LDFLAGS: -framework JavaScriptCore
#include <JavaScriptCore/JavaScriptCore.h>
*/
import "C"

// Context is a JS context.
type Context interface {
	jsContext() C.JSContextRef
	// GlobalObject returns the global object associated with the context.
	GlobalObject() Object
}

// GlobalContext is a global JS context.
type GlobalContext struct {
	ctx C.JSGlobalContextRef
}

// NewGlobalContext creates a new global JS context.
func NewGlobalContext() GlobalContext {
	return GlobalContext{C.JSGlobalContextCreate(nil)}
}

func (c GlobalContext) jsContext() C.JSContextRef {
	return C.JSContextRef(c.ctx)
}

// GlobalObject returns the global object associated with the context.
func (c GlobalContext) GlobalObject() Object {
	return Object{C.JSContextGetGlobalObject(c.ctx)}
}

// Release releases the context.
func (c GlobalContext) Release() {
	C.JSGlobalContextRelease(c.ctx)
}

type contextWrapper struct {
	ctx C.JSContextRef
}

func (c contextWrapper) jsContext() C.JSContextRef {
	return c.ctx
}

func (c contextWrapper) GlobalObject() Object {
	return Object{C.JSContextGetGlobalObject(c.ctx)}
}
