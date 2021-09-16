package main

import (
	"fmt"

	"github.com/phomola/jscore"
)

func main() {
	ctx := jscore.NewGlobalContext()
	defer ctx.Release()

	obj := ctx.GlobalObject()
	fmt.Println(obj.Has(ctx, "test"), obj.Get(ctx, "test").String(ctx), obj.Get(ctx, "test").Type(ctx))
	obj.Set(ctx, "test", jscore.NewString(ctx, "TEST"))
	fmt.Println(obj.Has(ctx, "test"), obj.Get(ctx, "test").String(ctx), obj.Get(ctx, "test").Type(ctx))
	obj.Set(ctx, "test", jscore.NewArray(ctx, jscore.NewNumber(ctx, 1234), jscore.NewString(ctx, "Test")).Value())
	fmt.Println(obj.Has(ctx, "test"), obj.Get(ctx, "test").String(ctx), obj.Get(ctx, "test").Type(ctx))
	obj.Set(ctx, "test", jscore.NewObject(ctx).Value())
	fmt.Println(obj.Has(ctx, "test"), obj.Get(ctx, "test").String(ctx), obj.Get(ctx, "test").Type(ctx))

	script := `12.34`
	corr := jscore.CheckScriptSyntax(ctx, script)
	fmt.Println("check:", corr)
	if corr {
		r := jscore.EvaluateScript(ctx, script)
		fmt.Println("result:", r.String(ctx), r.Type(ctx))
	}
}
