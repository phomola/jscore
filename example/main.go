package main

import (
	"fmt"

	"github.com/phomola/jscore"
)

type Person struct {
	Name string
	Age  int
}

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

	obj = jscore.NewObject(ctx)
	obj.Set(ctx, "p1", jscore.NewNumber(ctx, 1))
	obj.Set(ctx, "p2", jscore.NewNumber(ctx, 2))
	obj.Set(ctx, "p3", jscore.NewNumber(ctx, 3))
	r := obj.Value().Interface(ctx)
	fmt.Printf("%v %T\n", r, r)

	obj = jscore.NewGoObject(ctx, &Person{"Maeve", 18})
	r = obj.Value().Interface(ctx)
	fmt.Printf("%v %T\n", r, r)

	obj = jscore.NewArray(ctx, jscore.NewNumber(ctx, 1234), jscore.NewString(ctx, "Test"))
	r = obj.Value().Interface(ctx)
	fmt.Printf("%v %T\n", r, r)
}
