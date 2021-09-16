package main

import (
	"fmt"

	"github.com/phomola/jscore"
)

type Person struct {
	Name string `jscore:"name"`
	Age  int    `jscore:"age"`
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

	obj.Set(ctx, "person", jscore.NewGoObject(ctx, &Person{"Aoife", 17}).Value())
	fmt.Println("person:", obj.Has(ctx, "person"), obj.Get(ctx, "person").String(ctx), obj.Get(ctx, "person").Type(ctx))
	script := `person.age`
	if jscore.CheckScriptSyntax(ctx, script) {
		r := jscore.EvaluateScript(ctx, script)
		fmt.Println("result:", r.Interface(ctx))
	} else {
		panic("syntax error")
	}

	script = `[1, 2, 3]`
	if jscore.CheckScriptSyntax(ctx, script) {
		r := jscore.EvaluateScript(ctx, script)
		fmt.Println("result:", r.Interface(ctx), r.Object(ctx).Slice(ctx), r.Type(ctx))
	} else {
		panic("syntax error")
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

	obj = jscore.NewGoObject(ctx, &Person{"Moira", 19})
	r = obj.Value().Interface(ctx)
	fmt.Printf("%v %T\n", r, r)

	obj = jscore.NewArray(ctx, jscore.NewNumber(ctx, 1234), jscore.NewString(ctx, "Test"))
	r = obj.Value().Interface(ctx)
	fmt.Printf("%v %T\n", r, r)
}
